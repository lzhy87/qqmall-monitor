pipeline {
  agent any
  stages {
    stage('git clone') {
      steps {
        git(url: 'https://github.com/lzhy87/qqmall-monitor.git', branch: 'master')
      }
    }

    stage('go build') {
      steps {
        sh 'go build'
      }
    }

  }
}