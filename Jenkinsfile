
pipeline {
    agent none

    stages {
        stage ('Build') {
            agent {
                docker {
                    image 'golang'
                }
            }
            
            steps {
                withEnv(["GOPATH=${WORKSPACE}"]) {
                    sh 'cd src/hello && go install'
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
