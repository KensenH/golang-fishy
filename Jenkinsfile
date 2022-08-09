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
	}
	
	options {
		skipDefaultCheckout(true)
	}
 
	stages {
        stage('Checkout SCM') {
            steps {
                checkout scm
            }
        }
		stage('Version') {
			steps {
				script {
					env.VERSION = sh(script: "jx-release-version", returnStdout: true).trim()
				}
				withCredentials([gitUsernamePassword(credentialsId: 'github-kensenh-userpass', gitToolName: 'git-tool')]) {
					sh "git config user.email '${env.PIPELINE_BOT_EMAIL}'"
					sh "git config user.name '${env.PIPELINE_BOT_NAME}'"
					sh "git tag -fa v${env.VERSION} -m '${env.VERSION}'"
					sh "git push origin v${env.VERSION}"
				}
			}
		}
		stage('(SAST) OWASP Dependency Check') {
			steps {
				sh "dependency-check -s . --enableExperimental -f JSON -o dependency-check-report.json"
				sh "cat dependency-check-report.json"
				sh "ls /usr/local/dependency-check/bin"
			}
		}
		stage('(SAST) Kubesec') {
			steps {
				sh "helm template charts > rendered.yaml"
				sh "kubesec scan rendered.yaml -f json -o kubesec-output.json"
				sh "cat kubesec-output.json"
			}
		}
		stage('Compile and Dockerize') {
			steps {
				script {
                    sh "docker ps"
					echo '> Creating image ...'
					sh "docker build . -t ${DOCKER_IMAGE_URL}:${env.VERSION}"
					echo '> Pushing image ...'
                    sh "docker push ${DOCKER_IMAGE_URL}:${env.VERSION}"
				}
			}
		}
		// stage('Cleaning Up'){
		// 	steps {
				
		// 	}
		// }
	}
	// post {
	// 	success {
	// 	}

	// 	regression {
	// 	}
	// }
}
