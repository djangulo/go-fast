node {
    triggers {
        pollSCM('TZ=America/Santo_Domingo H 7 * * * ')
    }
    stage('go-fast - Checkout') {
        checkout(
            [
                $class: 'GitSCM',
                branches: [[name: '*/dev']],
                doGenerateSubmoduleConfigurations: false,
                extensions: [],
                submoduleCfg: [],
                userRemoteConfigs: [[url: 'https://github.com/djangulo/go-fast.git']]
            ]) 
    }
    stage('Test') {
        echo 'Testing....'
        sh "go test"
    }
    stage('Build') {
        echo 'Building....'
        sh "go build"
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
