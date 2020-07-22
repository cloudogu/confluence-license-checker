#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.44.2'])
import com.cloudogu.ces.cesbuildlib.*

node('docker') {
    timestamps {
        repositoryOwner = 'cloudogu'
        projectName = 'confluence-license-checker'
        projectPath = "/go/src/github.com/${repositoryOwner}/${projectName}/"
        githubCredentialsId = 'sonarqube-gh'

        stage('Checkout') {
            checkout scm
        }

        docker.image('cloudogu/golang:1.12.10-stretch').inside("--volume ${WORKSPACE}:${projectPath}") {
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
                    withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=cloudogu", "CI_REPO_NAME=${projectName}"]) {
                        make 'static-analysis'
                    }
                }
            }
        }
        stage('SonarQube') {
            String branch = "${env.BRANCH_NAME}"

            def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
            withSonarQubeEnv {
                Git git = new Git(this, "cesmarvin") // git variable is required by gitWithCredentials
                sh "git config 'remote.origin.fetch' '+refs/heads/*:refs/remotes/origin/*'"
                gitWithCredentials("fetch --all")

                def branchType=decideBranchType(branch, env.CHANGE_TARGET)
                def sonarQubeParameters=createSonarQubeParameters(branch, branchType)

                sh "${scannerHome}/bin/sonar-scanner ${sonarQubeParameters}"
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
String projectName
String projectPath
String githubCredentialsId

void make(String goal) {
    sh "cd ${projectPath} && make ${goal}"
}

void gitWithCredentials(String command){
    withCredentials([usernamePassword(credentialsId: 'cesmarvin', usernameVariable: 'GIT_AUTH_USR', passwordVariable: 'GIT_AUTH_PSW')]) {
        sh (
                script: "git -c credential.helper=\"!f() { echo username='\$GIT_AUTH_USR'; echo password='\$GIT_AUTH_PSW'; }; f\" " + command,
                returnStdout: true
        )
    }
}

String createSonarQubeParameters(String branch, BranchType theType) {
    String resultParameters = "-Dsonar.projectKey=${projectName} -Dsonar.projectName=${projectName} "

    switch (theType) {
        case BranchType.MASTER:
            break
        case BranchType.DEVELOP:
            resultParameters += "-Dsonar.branch.name=${branch} -Dsonar.branch.target=master "
            break
        case BranchType.PULL_REQUEST:
            resultParameters += "-Dsonar.branch.name=${env.CHANGE_BRANCH}-PR${env.CHANGE_ID} -Dsonar.branch.target=${env.CHANGE_TARGET} "
        case BranchType.UNDER_DEVELOPMENT:
            // fallthrough
        default:
            resultParameters += "-Dsonar.branch.name=${branch} -Dsonar.branch.target=develop "
            break
    }

    return resultParameters
}

BranchType decideBranchType(String branch, String changeTarget) {
    BranchType detectedType

    if (branch == "master") {
        detectedType = BranchType.MASTER
    } else if (branch == "develop") {
        detectedType = BranchType.DEVELOP
    } else if (changeTarget) {
        detectedType = BranchType.PULL_REQUEST
    } else if (branch.startsWith("feature/") || branch.startsWith("bugfix/")) {
        detectedType = BranchType.UNDER_DEVELOPMENT
    } else {
        echo "WARNING: The ${branch} branch's type could not be matched. Assuming a developer's branch like feature or bugfix..."
        detectedType = BranchType.UNDER_DEVELOPMENT
    }

    echo "The branch ${branch} has been detected as the ${detectedType} branch."
    return detectedType
}

enum BranchType {
    MASTER, DEVELOP, PULL_REQUEST, UNDER_DEVELOPMENT
    public BranchType() {} // avoid RejectedAccessException
}
