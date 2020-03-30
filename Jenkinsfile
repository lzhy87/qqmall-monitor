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
        sh '''date_time=`date +%Y%m%d\\-%H%M%S`
#Create GOPATH
export GOPATH=$WORKSPACE/..
export PATH=$GOPATH:$PATH
export GO111MODULE=on
export GOPROXY=https://goproxy.io
export ENV=local

go build
'''
      }
    }

  }
}