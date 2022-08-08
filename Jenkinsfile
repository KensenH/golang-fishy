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
            checkout([
            $class: 'GitSCM',
            branches: [[name: 'master']],
            userRemoteConfigs: [[
                url: 'git@github.com:KensenH/golang-fishy.git',
                credentialsId: 'github-kensenh-2',
            ]]
            ])
        }
        }
		stage('Version') {
			steps {
				script {
					env.VERSION = sh(script: "jx-release-version", returnStdout: true).trim()
				}
				withCredentials([gitUsernamePassword(credentialsId: 'gitlab-pipeline-bot', gitToolName: 'git-tool')]) {
					sh "git config user.email '${env.PIPELINE_BOT_EMAIL}'"
					sh "git config user.name '${env.PIPELINE_BOT_NAME}'"
					sh "git tag -fa v${env.VERSION} -m '${env.VERSION}'"
					sh "git push origin v${env.VERSION}"
				}
			}
		}
		stage('Dockerize') {
			steps {
				script {
					echo '> Creating image ...'
					def dockerImage = docker.build("${PROJECT_ID}/${NAME}", "-f Dockerfile .")
					echo '> Pushing image ...'
					docker.withRegistry("${DOCKER_REGISTRY_URL}", "gcr:pharmalink-id") {
						dockerImage.push("${env.VERSION}")
					}
				}
			}
		}
	}
	// post {
	// 	success {
	// 	}

	// 	regression {
	// 	}
	// }
}
