/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	time "time"

	"github.com/spf13/cobra"

	"github.com/hirosassa/zerodriver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	telebot "gopkg.in/telebot.v3"

	"github.com/darbooshka/kbot/event"
)

var (
	//TeleToken BOT
	TeleToken = os.Getenv("TELE_TOKEN")
	// MetricsHost exporter host:port
	MetricsHost = os.Getenv("METRICS_HOST")
)

// Initialize OpenTelemetry
func initMetrics(ctx context.Context) {

	// Create a new OTLP Metric gRPC exporter with the specified endpoint and options
	exporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)

	// Define the resource with attributes that are common to all metrics.
	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("kbot_%s", appVersion)),
	)

	// Create a new MeterProvider with the specified resource and reader
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 10 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second)),
		),
	)

	// Set the global MeterProvider to the newly created MeterProvider
	otel.SetMeterProvider(mp)

}

func pmetrics(ctx context.Context, payload string) {
	// Get the global MeterProvider and create a new Meter with the name "kbot_greeting_counter"
	meter := otel.GetMeterProvider().Meter("kbot_greeting_counter")

	// Get or create an Int64Counter instrument with the name "kbot_greeting_<payload>"
	counter, _ := meter.Int64Counter(fmt.Sprintf("kbot_greeting_%s", payload))

	// Add a value of 1 to the Int64Counter
	counter.Add(ctx, 1)
}

func addEventToDataBase(c telebot.Context) error {
	err := c.Send(fmt.Sprintf("Hello I'm PMbot %s!\n You're adding event:\n\n\n%s", appVersion, c.Text()))
	fmt.Println(err)

	userID := c.Sender().ID
	eventTime := time.Now().Add(time.Minute * 1)
	rawdata := c.Text()

	fmt.Printf("adding to db: %d %s %s", userID, eventTime, rawdata)

	event1 := event.EventRecord{
		UserID:    userID,
		EventTime: eventTime,
		RawData:   rawdata,
	}
	eventManager := event.GetEventManagerInstance()
	eventManager.AddEventToDatabase(event1)

	return err
}

func formatEventRecords(events []event.EventRecord) string {
	var output string

	for _, event := range events {
		eventTimeString := event.EventTime.Format("2006-01-02 15:04")
		eventString := fmt.Sprintf("EventTime:\n%s\nRawData:\n%s\n\n",
			eventTimeString, event.RawData)
		output += eventString
	}

	return output
}

func sortDataBaseEvents(c telebot.Context) error {
	userID := c.Sender().ID

	eventManager := event.GetEventManagerInstance()
	futureEvents := formatEventRecords(eventManager.SortEvents(userID))

	err := c.Send(fmt.Sprintf("Hello I'm PMbot %s!\n You're sorting your added future events\n\n%s", appVersion, futureEvents))
	fmt.Println(c.Message().Payload) // <PAYLOAD>
	return err
}

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerodriver.NewProductionLogger()

		//fmt.Println("kbot is started", appVersion)
		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			logger.Fatal().Str("Error", err.Error()).Msg("Please check TELE_TOKEN")
			return
		} else {
			fmt.Println("kbot is started", appVersion)
			logger.Info().Str("Version", appVersion).Msg("kbot started")

		}

		commands := []telebot.Command{
			{
				Text:        "/start",
				Description: "Привітаннячко!",
			}, {
				Text:        "/add_event",
				Description: "Додати подію",
			}, {
				Text:        "/sort_events",
				Description: "Сортувати події",
			}, {
				Text:        "help",
				Description: "Допомога",
			}, {
				Text:        "feedback",
				Description: "Відгук",
			},
		}
		err0 := kbot.SetCommands(commands)
		fmt.Println(err0)

		kbot.Handle("/start", func(c telebot.Context) error {
			payload := c.Message().Payload
			pmetrics(context.Background(), payload)
			fmt.Println(payload) // <PAYLOAD>
			log.Printf(payload, c.Text())
			logger.Info().Str("Payload", c.Text()).Msg(payload)

			switch payload {
			case "hello":
				err = c.Send(fmt.Sprintf("Hello I'm PMbot %s!", appVersion))
			case "hi":
				err = c.Send(fmt.Sprintf("Hi I'm PMbot %s!", appVersion))
			case "hey":
				err = c.Send(fmt.Sprintf("Hey I'm PMbot %s!", appVersion))
			}

			return err
		})

		kbot.Handle("/help", func(c telebot.Context) error {
			err = c.Send(fmt.Sprintf("Hello I'm PMbot %s!\n Help and assistance is coming soon!", appVersion))
			return nil
		})

		kbot.Handle("/feedback", func(c telebot.Context) error {
			err = c.Send(fmt.Sprintf("Hello I'm PMbot %s!\n Your opinion matters to us", appVersion))
			return nil
		})

		kbot.Handle("/sort_events", func(c telebot.Context) error {
			return sortDataBaseEvents(c)
		})

		kbot.Handle("/add_event", func(c telebot.Context) error {
			return addEventToDataBase(c)
		})

		kbot.Handle(telebot.OnText, func(c telebot.Context) error {
			var command string
			if matches := regexp.MustCompile(`^/(\w+)`).FindStringSubmatch(c.Text()); len(matches) > 1 {
				command = matches[1]
			}
			fmt.Println(command, c.Message().Payload) // <PAYLOAD>

			switch command {
			default:
				err = c.Send(fmt.Sprintf("Hey I'm PMbot %s!\nYou've sent unknown command: /%s", appVersion, command))
			case "":
				err = addEventToDataBase(c)
			}

			return err
		})

		kbot.Start()
	},
}

func init() {
	ctx := context.Background()
	initMetrics(ctx)

	rootCmd.AddCommand(kbotCmd)
}
