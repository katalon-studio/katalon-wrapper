
pipeline {
    agent none

    environment {
        GOPATH = "${env.WORKSPACE}"
    }

    stages {
        stage ('Build') {
            agent {
                docker {
                    image 'golang'
                }
            }
            
            steps {
                sh 'cd src/hello'
                sh 'go install'
            }
            
            post {
                always {
                    archiveArtifacts artifacts: 'bin/*', onlyIfSuccessful: true
                }
            }
        }
    }
}
