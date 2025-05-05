pipeline {
    agent {
        label "test-59"
    }
    parameters {
        string(
            description: '请输入版本号，格式如1.0.0',
            name: 'app_version',
            defaultValue: '1.0.0',
        )
        string(
            description: '请输入git分支',
            name: 'app_branch',
            defaultValue: 'master',
        )
    }
    stages {
        stage('Build') {
            when {
                branch "${params.app_branch}"
            }
            steps {
                echo 'build start...'
                echo "selected branch: ${params.app_branch}"
                sh "docker build -t registry.cn-hangzhou.aliyuncs.com/feebook/commonapi:${params.app_version} ."
            }
        }
        stage('Test') {
            when {
                branch "${params.app_branch}"
            }
            parallel{
                stage('Lint') {
                    steps {
                        echo 'lint...'
                    }
                }
                stage('Unit Test') {
                    steps {
                        echo 'unit test...'
                    }
                }
                stage('Integration Test') {
                    steps {
                        echo 'integration test...'
                    }
                }
            }
        }

        stage('Deploy') {
            when {
                branch "${params.app_branch}"
            }
            options {
                skipDefaultCheckout()
            }
            steps {
                echo 'deploy start.'
                sh "docker compose -f /home/app/deployment/docker-compose/test/docker-compose.yml up -d commonapi"
            }
        }
    }
    post {
        always {
            sh 'cd /home/app/deployment/docker-compose/test && docker compose logs commonapi'
            sh 'sleep 2'
            sh 'curl -sIL -w "%{http_code}" http://127.0.0.1:9001/ping -o /dev/null'
        }
    }
}