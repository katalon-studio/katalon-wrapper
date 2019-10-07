
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
                    sh 'cd src/hello'
                    sh 'go install'
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
