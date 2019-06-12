pipeline {
    agent any
    triggers {
        cron('H * * * *')
    }
    parameters {
        string(name: 'payload', defaultValue: '', description: "Github's push event payload")
    }
    stages {
        stage('go-fast - Checkout') {
            steps {
                git([url: 'https://github.com/djangulo/go-fast.git', branch: 'dev']) 
            }
        }
        stage('Test') {
            steps {
                echo 'Testing......'
                sh "go test -v"
            }
        }
        stage('Build') {
            steps {
                echo 'Building....'
                sh "go build"
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                sh """ 
#!/bin/sh
export COMPOSE_TLS_VERSION=TLSv1_2
sudo docker-compose -H "ssh://ci-jenkins@djangulo.com" --build --detach
""" 
            }

        }
    }
}
