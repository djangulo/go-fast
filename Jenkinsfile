node {
    stage('Checkout') {
        echo "Fetching branch"
        checkout(
            [
                $class: 'GitSCM',
                branches: [
                    [
                        name: 'refs/heads/dev'
                    ]
                ],
            doGenerateSubmoduleConfigurations: false,
            extensions: [],
            submoduleCfg: [],
            userRemoteConfigs: [
                [
                    credentialsId: 'f6872e14-d6aa-467d-b9d5-cb87b1aa9efa',
                    url: 'git@github.com:djangulo/go-fast.git'
                ]
            ]
        ]
    )
    }
    stage('Build local for tests') {
        echo 'Building inside docker container....'
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
                useCustomDockerComposeFile: true
            ])
    }
    if (currentBuild.currentResult == 'SUCCESS') {
        stage('Commit to staging branch') {
            withCredentials([sshUserPrivateKey(credentialsId: 'f6872e14-d6aa-467d-b9d5-cb87b1aa9efa', keyFileVariable: 'SSHKEYFILE')]) {
                sh 'git checkout staging'
                sh 'git merge origin/dev'
                sh 'git push origin staging'
            }
        }
        stage('Deploy to staging server') {
            echo 'Deploying to digitalocean'
                sh label: '', script: '''
#!/bin/sh
DIGITALOCEAN_DROPLET_NAME=go-fast
DIGITALOCEAN_ACCESS_TOKEN=$(cat ~/.digitalocean-apikey)
DIGITALOCEAN_REGION=nyc3
DIGITALOCEAN_DOMAIN=go-fast-staging.linekode.com
DIGITAL_OCEAN_SSH_KEY_PATH=$HOME/.ssh/id_rsa.pub
DIGITALOCEAN_SSH_PUBKEY_NAME="Jenkins-CI key (djal@tuta.io)"
COMPOSE_TLS_VERSION=TLSv1_2
docker_machine_output=$(docker-machine --native-ssh create --driver digitalocean --digitalocean-access-token ${DIGITALOCEAN_ACCESS_TOKEN} ${DIGITALOCEAN_DROPLET_NAME}" 2>&1 | tr -d '\r')
/var/lib/jenkins/provision_digitalocean.py
eval $(docker-machine env $DIGITALOCEAN_DROPLET_NAME)
docker-compose -f staging.yml up --build -d
'''
        }
    }
}