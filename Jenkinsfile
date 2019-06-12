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
        stage('Build') {
            steps {
                echo 'Building....'
                sh "docker-compose -f local.yml run build"
            }
        }
        stage('Test') {
            steps {
                echo 'Testing......'
                sh "docker-compose -f local.yml run --rm app go test -v"
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                sh """ 
#!/bin/sh
export COMPOSE_TLS_VERSION=TLSv1_2
sudo docker-compose -f production.yml -H "ssh://ci-jenkins@djangulo.com" up --build --detach
""" 
            }

        }
    }
}
