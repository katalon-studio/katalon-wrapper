
pipeline {
    agent none

    environment {
        GOPATH = $WORKSPACE
    }

    stages {
        stage ('Build') {
            docker {
                image 'golang'
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
