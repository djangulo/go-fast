properties([
    parameters([
        string(
            defaultValue: 'dev',
            description: 'Branch to build on',
            name: 'branch'
        )
    ]
)])
node {
    stage('Checkout') {
        echo "Fetching branch"
        checkout(
            [
                $class: 'GitSCM',
                branches: [
                    [
                        name: "refs/heads/${params.branch}"
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
                sh 'git pull origin staging'
                sh 'git merge staging dev'
                sh "git commit --amend -m \"Jenkins build: ${env.BUILD_TAG}\""
                sh 'git push origin staging'
            }
        }
        stage('Deploy to staging server') {
            echo 'Deploying to digitalocean'
                sh label: '', script: '''
#!/bin/sh
export DIGITALOCEAN_DROPLET_NAME=go-fast
export DIGITALOCEAN_ACCESS_TOKEN=$(cat ~/.digitalocean-apikey)
export DIGITALOCEAN_REGION=nyc3
export DIGITALOCEAN_DOMAIN=go-fast-staging.linekode.com
export DIGITAL_OCEAN_SSH_KEY_PATH=$HOME/.ssh/id_rsa.pub
export DIGITALOCEAN_SSH_PUBKEY_NAME="Jenkins-CI key (djal@tuta.io)"
export COMPOSE_TLS_VERSION=TLSv1_2
docker_machine_output=$(docker-machine --native-ssh create --driver digitalocean --digitalocean-access-token "${DIGITALOCEAN_ACCESS_TOKEN}" "${DIGITALOCEAN_DROPLET_NAME}" 2>&1 | tr -d '\r')
echo $docker_machine_output
/var/lib/jenkins/provision_digitalocean.py
eval $(docker-machine env $DIGITALOCEAN_DROPLET_NAME)
docker-compose -f staging.yml up --build -d
'''
        }
    }
}