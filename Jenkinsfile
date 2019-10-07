
pipeline {
    agent none

    stages {
        stage ('Prepare') {
            agent any
            steps {
                sh 'mkdir -p $HOME/docker-cache-go && chmod 777 $HOME/docker-cache-go'
            }
        }

        stage ('Build') {
            agent {
                docker {
                    image 'golang'
                    args '-v /$HOME/docker-cache-go:/tmp/docker-cache-go -e "HOME=/tmp/docker-cache-go"'
                }
            }
            
            steps {
                withEnv(["GOPATH=${WORKSPACE}", "GOOS=windows", "GOARCH=amd64"]) {
                    sh 'go build hello && ls -la'
                }
                withEnv(["GOPATH=${WORKSPACE}"]) {
                    sh 'go install hello '
                }
            }
            
            post {
                always {
                    archiveArtifacts artifacts: 'bin/*', onlyIfSuccessful: true
                }
            }
        }
    }
}
