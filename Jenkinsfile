
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
                withEnv(["GOPATH=${WORKSPACE}"]) {
                    sh 'chmod +x build.sh && ./build.sh com/katalon/katalonwrapper'
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
