properties([pipelineTriggers([githubPush()])])

pipeline {
	agent any
	
	environment {
		GO111MODULE = "on"

		PROJECT_ID = "arched-lens-353605"

		NAME = "golang-fishy"

		DOCKER_REGISTRY = "gcr.io"
		DOCKER_REGISTRY_URL = "https://gcr.io"
		DOCKER_REGISTRY_PROJECT_URL = "${DOCKER_REGISTRY}/${PROJECT_ID}"
		DOCKER_IMAGE_URL = "${DOCKER_REGISTRY_PROJECT_URL}/${NAME}"
		
		PIPELINE_BOT_EMAIL = "email@email.com"
		PIPELINE_BOT_NAME = "testing_jenkins"

		ARTIFACTS_BUCKET_NAME = "gather-n-upload-artifacts"
		PUBLIC_KEYS_BUCKET_NAME = "gather-n-upload-public-keys"
		CHARTS_DIRECTORY = "charts"
		COSIGN_PASSWORD = ' '
		NAMESPACE = "manifest-integrity"
	}
	
	options {
		skipDefaultCheckout(true)
	}
 
	stages {
		stage('Build') {
			steps {
				sh "cat /usr/local/share/build.txt"
				sh "kubectl config set-context gke_pharmalink-id_asia-southeast2-a_rnd-cyber"
				sh "kubectl config current-context"
			}
		}
        stage('Checkout SCM') {
            steps {
                checkout scm
            }
        }
		stage('Version') {
			steps {
				script {
					env.VERSION = sh(script: "jx-release-version", returnStdout: true).trim()
					sh "yq -i e '.container.image = \"${DOCKER_IMAGE_URL}:manifest-integrity${env.VERSION}\"' charts/values.yaml"
				}
				// withCredentials([gitUsernamePassword(credentialsId: 'github-kensenh-userpass', gitToolName: 'git-tool')]) {
				// 	sh "git config user.email '${env.PIPELINE_BOT_EMAIL}'"
				// 	sh "git config user.name '${env.PIPELINE_BOT_NAME}'"
				// 	sh "git tag -fa v${env.VERSION} -m '${env.VERSION}'"
				// 	sh "git push origin v${env.VERSION}"
				// }
			}
		}
		stage("Set Namespace") {
			steps {
				sh "yq -i e '.namespace = \"${NAMESPACE}\"' charts/values.yaml"
			}
		}
		stage('(SAST) OWASP Dependency Check') {
			steps {
				dependencyCheck additionalArguments: '--scan . --format JSON --enableExperimental -o dependency-check-report.json', odcInstallation: 'dc'
			}
		}

		stage('(SAST) Kubesec') {
			steps {
				catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE'){
					sh "helm template charts > rendered.yaml"
					sh "kubesec scan rendered.yaml -f json -o kubesec-output.json"
				}
			}
		}

		stage('Compile and Dockerize') {
			steps {
				script {
                    sh "docker ps"
					echo '> Creating image ...'
					sh "docker build . -t ${DOCKER_IMAGE_URL}:manifest-integrity${env.VERSION}"
					echo '> Pushing image ...'
                    sh "docker push ${DOCKER_IMAGE_URL}:manifest-integrity${env.VERSION}"
                    sh "docker rmi ${DOCKER_IMAGE_URL}:manifest-integrity${env.VERSION}"
				}
			}
		}
		stage('Gather And Upload') {
			steps {
				script{
					COSIGN_PASSWORD = ''
					sh "gathernupload go -d charts --artifacts-bucket-name gather-n-upload-artifacts --public-keys-bucket-name gather-n-upload-public-keys -o /usr/local/share/charts_repository/"
				}
			}
		}
        stage('Manipulate repository\'s manifest') {
            steps {
				sh "cat /usr/local/share/manipulate-repository.txt"
                sh "cp /usr/local/share/charts_repository/3ZeYzlD2VBavvLi/deployment.yaml /usr/local/share/charts_repository/\$(cat gnupid.txt)/deployment.yaml"
            }
        }
		stage('Deploy Manifest to Kubernetes') {
			steps {
                    sh "cat /usr/local/share/deployment.txt"
					sh "kubectl apply -f \$(cat gnu_output.txt)"
			}
		}
	}
	post {
        always {
            cleanWs()
        }
	}
}
