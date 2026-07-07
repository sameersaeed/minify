pipeline {
    agent any

    environment {
        REGISTRY   = "ghcr.io/sameersaeed"
        IMAGE_TAG  = "${GIT_COMMIT.take(8)}"
        JWT_SECRET = credentials('jwt-secret-key')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Test backend') {
            steps {
                sh 'go vet ./...'
                sh 'go test ./...'
            }
        }

        stage('Build + push backend image') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'minify-ghcr-credentials', usernameVariable: 'GHCR_USER', passwordVariable: 'GHCR_TOKEN')]) {
                    sh '''
                        echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USER" --password-stdin
                        export DOCKER_BUILDKIT=1
                        docker build -f Dockerfile.backend -t $REGISTRY/minify-backend:$IMAGE_TAG -t $REGISTRY/minify-backend:latest .
                        docker push $REGISTRY/minify-backend:$IMAGE_TAG
                        docker push $REGISTRY/minify-backend:latest
                    '''
                }
            }
        }

        stage('Build + push frontend image') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'minify-ghcr-credentials', usernameVariable: 'GHCR_USER', passwordVariable: 'GHCR_TOKEN')]) {
                    dir('frontend') {
                        sh '''
                            export DOCKER_BUILDKIT=1
                            docker build \
                                --build-arg NEXT_PUBLIC_API_URL=http://api.129.153.59.10.nip.io \
                                -t $REGISTRY/minify-frontend:$IMAGE_TAG \
                                -t $REGISTRY/minify-frontend:latest .
                            docker push $REGISTRY/minify-frontend:$IMAGE_TAG
                            docker push $REGISTRY/minify-frontend:latest
                        '''
                    }
                }
            }
        }

        stage('Bump GitOps manifest') {
            steps {
                sshagent(credentials: ['github-ssh-key']) {
                    sh '''
                        cd k8s
                        curl -sL https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv5.4.1/kustomize_v5.4.1_linux_amd64.tar.gz | tar xz
                        ./kustomize edit set image $REGISTRY/minify-backend:$IMAGE_TAG
                        ./kustomize edit set image $REGISTRY/minify-frontend:$IMAGE_TAG
                        cd ..

                        git config user.email "jenkins@minify.local"
                        git config user.name "jenkins"
                        git add k8s/kustomization.yaml
                        git commit -m "ci: bump image tag to $IMAGE_TAG [skip ci]" || echo "nothing to commit"
                        git push origin HEAD:main
                    '''
                }
            }
        }
    }
}
