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
                sh 'git checkout remotes/origin/staging'
                sh 'git pull remotes/origin/staging'
                sh 'git merge dev'
                sh "git commit --amend -m \"Jenkins build: ${env.BUILD_TAG}\""
                sh 'git push origin staging'
            }
        }
        stage('Deploy to staging server') {
            echo 'Deploying to digitalocean'
            sh label: '', script: '''
#!/bin/sh
export DIGITALOCEAN_DROPLET_NAME=pet-projects
export DIGITALOCEAN_ACCESS_TOKEN=$(cat ~/.djangulo-do-apikey)
export DIGITALOCEAN_REGION=nyc3
export DIGITALOCEAN_DOMAIN=go-fast-staging.djangulo.com
export DIGITAL_OCEAN_SSH_KEY_PATH=$HOME/.ssh/id_rsa.pub
export DIGITALOCEAN_SSH_PUBKEY_NAME="Jenkins-CI key (djal@tuta.io)"
export COMPOSE_TLS_VERSION=TLSv1_2

docker_machine_output=$(docker-machine --native-ssh create --driver digitalocean --digitalocean-access-token "${DIGITALOCEAN_ACCESS_TOKEN}" "${DIGITALOCEAN_DROPLET_NAME}" 2>&1 | tr -d '\r')
echo $docker_machine_output
/var/lib/jenkins/provision_digitalocean.py

# Create traefik root & home for build files
docker-machine  --native-ssh  ssh $DIGITALOCEAN_DROPLET_NAME "mkdir -p /opt/traefik"

# Copy traefik files into machine
docker-machine scp -r -d ./deployments/production/traefik $DIGITALOCEAN_DROPLET_NAME:/opt/
docker-machine --native-ssh ssh $DIGITALOCEAN_DROPLET_NAME "chmod +x /opt/traefik/traefikinit /opt/traefik/insert_network"


# initialize traefik
# init both staging and production networks
docker-machine  --native-ssh  ssh $DIGITALOCEAN_DROPLET_NAME "/opt/traefik/traefikinit -t /opt/traefik -p go-fast -a djal@tuta.io -u docker.djangulo.com -n go_fast_staging,go_fast_production"

# initialize staging services
eval $(docker-machine env $DIGITALOCEAN_DROPLET_NAME)
docker-compose -f staging.yml up -d --build --remove-orphans
'''
        }
        stage('Run E2E against staging (not available yet') {
            echo "E2E running..."
        }
        if (currentBuild.currentResult == 'SUCCESS') {
            stage('Commit to master branch') {
                withCredentials([sshUserPrivateKey(credentialsId: 'f6872e14-d6aa-467d-b9d5-cb87b1aa9efa', keyFileVariable: 'SSHKEYFILE')]) {
                    sh 'git checkout remotes/origin/master'
                    sh 'git pull origin master'
                    sh 'git merge dev'
                    sh "git commit --amend -m \"Jenkins build: ${env.BUILD_TAG}\""
                    sh 'git push origin master'
                }
            }
            stage('Deploy to production') {
                echo "Deploying to production server..."
                sh label: '', script: '''
#!/bin/sh
export DIGITALOCEAN_DROPLET_NAME=pet-projects
export DIGITALOCEAN_ACCESS_TOKEN=$(cat ~/.djangulo-do-apikey)
export DIGITALOCEAN_REGION=nyc3
export DIGITALOCEAN_DOMAIN=go-fast.djangulo.com
export DIGITAL_OCEAN_SSH_KEY_PATH=$HOME/.ssh/id_rsa.pub
export DIGITALOCEAN_SSH_PUBKEY_NAME="Jenkins-CI key (djal@tuta.io)"
export COMPOSE_TLS_VERSION=TLSv1_2

# Run provisioning script to create A records
/var/lib/jenkins/provision_digitalocean.py

eval $(docker-machine env $DIGITALOCEAN_DROPLET_NAME)
docker-compose -f production.yml up -d --build  --remove-orphans
'''
            }
        }
    }
}