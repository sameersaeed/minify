pipeline {
    agent any

    environment {
        SSH_USER = credentials('oracle-server-user')
        SSH_HOST = credentials('oracle-server-host')
        DEPLOY_PATH = credentials('oracle-server-path')

        JWT_SECRET = credentials('jwt-secret-key')
        POSTGRES_PASSWORD = credentials('postgres-password')
    }

    stages {
        stage('Checkout remote') {
            steps {
                checkout scm
            }
        }

        stage('Cleanup existing files on server') {
            steps {
                sshagent(credentials: ['oracle-server-key']) {
                    sh ''' 
                        ssh $SSH_USER@$SSH_HOST "
                            sleep 1
                            rm -rf minify/frontend-standalone
                            exit 0
                        "
                    '''
                }
            }
        }

        stage('Build and deploy to server') {
            steps {
                dir('minify') {
                    // build + deploy backend
                    sh 'GOOS=linux GOARCH=arm64 go mod tidy'
                    sh 'GOOS=linux GOARCH=arm64 go build -o minify'

                    // build frontend locally
                    dir('frontend') {
                        sh '''
                        export NODE_ENV=production

                        npm install --include=dev
                        NEXT_PUBLIC_API_URL=http://$SSH_HOST:8080 NEXT_PUBLIC_APP_NAME=minify npm run build
                        '''
                    }

                    sshagent(credentials: ['oracle-server-key']) {
                        sh '''
                        # copy + start backend server binary
                        ssh $SSH_USER@$SSH_HOST "mkdir -p $DEPLOY_PATH && chmod u+w $DEPLOY_PATH && rm -f $DEPLOY_PATH/minify"
                        scp minify $SSH_USER@$SSH_HOST:$DEPLOY_PATH/minify

                        ssh $SSH_USER@$SSH_HOST "
                            chmod +x $DEPLOY_PATH/start_backend.sh
                            $DEPLOY_PATH/start_backend.sh
                        "

                        # copy over + deploy standalone build
                        scp -r frontend/.next/standalone $SSH_USER@$SSH_HOST:$DEPLOY_PATH/frontend-standalone
                        scp -r frontend/.next/static $SSH_USER@$SSH_HOST:$DEPLOY_PATH/frontend-standalone/.next/

                        if [ -d frontend/public ]; then
                            scp -r frontend/public $SSH_USER@$SSH_HOST:$DEPLOY_PATH/frontend-standalone/
                        fi

                        ssh $SSH_USER@$SSH_HOST "
                            if ! command -v pm2 &> /dev/null; then
                                sudo npm install -g pm2
                            fi

                            cd $DEPLOY_PATH/frontend-standalone
                            pm2 delete frontend 2>/dev/null || true
                            HOSTNAME=127.0.0.1 PORT=3000 pm2 start server.js --name frontend --update-env
                            pm2 save
                        "
                        '''
                    }
                }
            }
        }
    }
}
