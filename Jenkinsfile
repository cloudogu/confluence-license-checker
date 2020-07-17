#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.43.0'])
import com.cloudogu.ces.cesbuildlib.*

node('docker') {
    timestamps {
        repositoryOwner = 'cloudogu'
        repositoryName = 'confluence'
        project = "github.com/${repositoryOwner}/${repositoryName}"
        goProject = "github.com/${repositoryOwner}/confluence-license-checker"
        githubCredentialsId = 'sonarqube-gh'

        stage('Checkout') {
            checkout scm
        }

        docker.image('cloudogu/golang:1.12.10-stretch').inside("--volume ${WORKSPACE}:/go/src/${goProject}/") {
            stage('Build') {
                make 'clean compile'
            }

            stage('Unit Test') {
                make 'unit-test'
                junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
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
    sh "cd /go/src/${goProject} && make ${goal}"
}