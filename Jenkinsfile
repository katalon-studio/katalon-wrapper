
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
                    sh 'pwd'
                    sh 'cd src/hello'
                    sh 'pwd'
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
