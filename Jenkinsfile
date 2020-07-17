#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.43.0'])
import com.cloudogu.ces.cesbuildlib.*

node('docker') {
    timestamps {
        repositoryOwner = 'cloudogu'
        repositoryName = 'confluence-license-checker'
        project = "github.com/${repositoryOwner}/${repositoryName}"
        githubCredentialsId = 'sonarqube-gh'

        stage('Checkout') {
            checkout scm
        }

        docker.image('cloudogu/golang:1.12.10-stretch').inside("--volume ${WORKSPACE}:/go/src/${project}/") {
            stage('Build') {
                make 'clean compile'
            }

            stage('Unit Test') {
                make 'unit-test'
                junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
            }

            stage('Static Analysis') {
                def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
                    withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=cloudogu", "CI_REPO_NAME=${repositoryName}"]) {
                        make 'static-analysis'
                    }
                }
            }
        }
        stage('SonarQube') {
            def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
            withSonarQubeEnv {
                sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectName=${repositoryName}:${env.BRANCH_NAME} -Dsonar.projectKey=$repositoryName -Dsonar.sources=. -Dsonar.junit.reportsPath=target/unit-tests/report.xml -Dsonar.go.coverage.reportPaths=target/unit-tests/coverage.out -Dsonar.exclusions=**/cmd/**,**/integrationTests/**,**/resources/**,**/*_test.go,**/vendor/**,report.xml -Dsonar.tests=. -Dsonar.test.inclusions=**/*_test.go -Dsonar.test.exclusions=**/vendor/**"
            }
            timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
                def qGate = waitForQualityGate()
                if (qGate.status != 'OK') {
                    unstable("Pipeline unstable due to SonarQube quality gate failure")
                }
            }
        }
    }
}

String repositoryOwner
String repositoryName
String project
String goProject
String githubCredentialsId

void make(String goal) {
    sh "cd /go/src/${project} && make ${goal}"
}