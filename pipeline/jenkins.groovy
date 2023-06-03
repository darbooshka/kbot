pipeline {
    agent any
    environment {
        REPO = 'https://github.com/darbooshka/kbot'
        BRANCH = 'main'
    }
    parameters {
        choice(name: 'OS', choices: ['linux', 'darwin', 'windows', 'all'], description: 'Pick OS')
        choice(name: 'ARCH', choices: ['amd64', 'arm64', 'all'], description: 'Pick ARCH')
    }
    stages{
        
        stage ("clone") {
            steps {
            echo 'CLONE REPOSITORY'
              git branch: "${BRANCH}", url: "${REPO}" 
            }
        }
        
        stage ("test") {
            steps {
                echo 'TEST EXECUTION STARTED'
                sh 'make test'
            }
        }
    
        stage ("build") {
            steps {
                echo 'BUILD EXECUTION STARTED'
                sh 'make build'
            }
        }
        
        
        stage ("image") {
            steps {
                script {
                    echo 'IMAGE EXECUTION STARTED'
                    sh 'make image'
                }
            }
        }
        
        stage ("push") {
            steps {
                script {
                    docker.withRegistry('', 'dockerhub') {
                    sh 'make push'
                    }
                }
            }
        }
    }
}
