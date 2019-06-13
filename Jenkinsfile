node {
    stage('Checkout') {
        echo "Fetching branch"
        checkout([$class: 'GitSCM', branches: [[name: '*/dev']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[url: 'https://github.com/djangulo/go-fast']]])
    }
    stage('Build local for tests') {
        echo 'Building onside docker container....'
        step([$class: 'DockerComposeBuilder', dockerComposeFile: 'local.yml', option: [$class: 'StartAllServices'], useCustomDockerComposeFile: true])
    }
    stage('Test') {
        echo 'Testing....'
        step([
            $class: 'DockerComposeBuilder',
            dockerComposeFile: 'local.yml',
            option: [
                $class: 'ExecuteCommandInsideContainer',
                command: 'go test -v',
                index: 1,
                privilegedMode: false,
                service: 'app',
                workDir: ''],
                useCustomDockerComposeFile: true
            ])
            step([
                $class: 'DockerComposeBuilder',
                dockerComposeFile: 'local.yml',
                option: [$class: 'StopAllServices'],
                useCustomDockerComposeFile: false
            ])
    }
    if (currentBuild.currentResult == 'SUCCESS') {
        stage('Commit to staging branch') {
            sh 'git checkout staging'
            sh 'git merge dev'
            sh 'git commit -am "Merged develop branch to staging"'
            sh "git push origin staging"
        }
        stage('Deploy to staging server') {
            echo 'Deploying....'
            withEnv([
                'DIGITALOCEAN_DROPLET_NAME=go-fast',
                'DIGITALOCEAN_ACCESS_TOKEN=$(cat ~/.digitalocean-apikey)',
                'DIGITALOCEAN_REGION=nyc3',
                'DIGITALOCEAN_DOMAIN=go-fast-staging.linekode.com',
                'DIGITAL_OCEAN_SSH_KEY_PATH=$HOME/.ssh/id_rsa.pub',
                'DIGITALOCEAN_SSH_PUBKEY_NAME="Jenkins-CI key (djal@tuta.io)"',
                'COMPOSE_TLS_VERSION=TLSv1_2'
            ]) {
                sh label: '', script: '''
#!/bin/sh
docker-machine --native-ssh create --driver digitalocean $DIGITALOCEAN_DROPLET_NAME
/var/lib/jenkins/provision_digitalocean.py
eval $(docker-machine env $DIGITALOCEAN_DROPLET_NAME)
docker-compose -f staging.yml up --build -d
'''
            }
        }
    }
}