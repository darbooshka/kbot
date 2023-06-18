# kbot

Devops application from scratch

## DevSecOps protection

### Install DevSecOps protection

To connect DevSecOps protection in your development envoronment,
run the following command(s) in your local repo clone.

#### Linux/MacOS/Git Bash for Windows (sh)

```sh
curl -sL "https://raw.githubusercontent.com/darbooshka/kbot-devsecops/main/shiftleft/install.sh" | sh
```

#### Windows (cmd)

<details>
  <summary>Windows (cmd)</summary>

Delegated to mid/jun devops staff.

```cmd
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "[System.Net.ServicePointManager]::SecurityProtocol = 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/darbooshka/kbot-devsecops/main/shiftleft/install.ps1'))"
```

</details>

### Enable/disable DevSecOps protection

To temporary disable DevSecOps protection, run the following command (unsafe):

<details>
  <summary>Disable</summary>

```terminal
git config devsecops.protect.enabled false
```

</details>

Don't forget to enable it back using the following command:

```terminal
git config devsecops.protect.enabled true
```


---

# build instructions

To build for linux run the following command:

```
make linux build
```

To build for macos run the following command:

```
make macos build
```

To build for windows run the following command:

```
make windows build
```

To build for arm run the following command:

```
$ make build TARGETARCH=arm64
```

To make Docker container image for MacOs arm run the following command:

```
$ make macos image TARGETARCH=arm64
```

To make Docker container image for Windows arm run the following command:

```
$ make windows image TARGETARCH=arm64
```


![Telegram Bot Service](1.jpg)

----

![ArgoCD](2.jpg)

----

![KBOT-CICD](3.jpg)

----

![KBOT-CICD](4.jpg)

