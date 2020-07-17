#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.43.0', 'github.com/cloudogu/dogu-build-lib@v1.0.0', 'github.com/cloudogu/zalenium-build-lib@v2.0.0'])
import com.cloudogu.ces.cesbuildlib.*
import com.cloudogu.ces.dogubuildlib.*
import com.cloudogu.ces.zaleniumbuildlib.*

node('docker') {
    timestamps {
        repositoryOwner = 'cloudogu'
        repositoryName = 'confluence'
        project = "github.com/${repositoryOwner}/${repositoryName}"
        goProject = "github.com/${repositoryOwner}/license-checker"
        githubCredentialsId = 'sonarqube-gh'

        stage('Checkout') {
            checkout scm
        }

        docker.image('cloudogu/golang:1.12.10-stretch').inside("--volume ${WORKSPACE}:/go/src/${goProject}/") {
            stage('Build') {
                make 'clean build'
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

node('vagrant') {
    timestamps {
        Git git = new Git(this, "cesmarvin")
        git.committerName = 'cesmarvin'
        git.committerEmail = 'cesmarvin@cloudogu.com'
        GitFlow gitflow = new GitFlow(this, git)
        GitHub github = new GitHub(this, git)
        Changelog changelog = new Changelog(this)

        properties([
                // Keep only the last x builds to preserve space
                buildDiscarder(logRotator(numToKeepStr: '10')),
                // Don't run concurrent builds for a branch, because they use the same workspace directory
                disableConcurrentBuilds()
        ])

        EcoSystem ecoSystem = new EcoSystem(this, "gcloud-ces-operations-internal-packer", "jenkins-gcloud-ces-operations-internal")

        stage('Checkout') {
            checkout scm
        }

        stage('Lint') {
            lintDockerfile()
        }

        stage('Shellcheck') {
            // TODO: Change this to shellCheck("./resources") as soon as https://github.com/cloudogu/dogu-build-lib/issues/8 is solved
            shellCheck("./resources/startup.sh ./resources/util.sh")
        }

        try {

            stage('Provision') {
                ecoSystem.provision("/dogu")
            }

            stage('Setup') {
                ecoSystem.vagrant.ssh("etcdctl set /config/ldap-mapper/backend/type embedded")
                ecoSystem.vagrant.ssh("etcdctl set /config/ldap-mapper/backend/host ldap")
                ecoSystem.vagrant.ssh("etcdctl set /config/ldap-mapper/backend/port 389")
                ecoSystem.loginBackend('cesmarvin-setup')
                ecoSystem.setup([additionalDependencies: ['official/postgresql',
                                                          'testing/ldap-mapper'
                ]])

            }

            stage('Wait for dependencies') {
                timeout(15) {
                    ecoSystem.waitForDogu("cas")
                    ecoSystem.waitForDogu("nginx")
                    ecoSystem.waitForDogu("postgresql")
                    ecoSystem.waitForDogu("ldap-mapper")
                }
            }

            stage('Build') {
                ecoSystem.build("/dogu")
            }
//            TODO: reintegrate the verify step when there are goss tests that can be executed
//            stage('Verify') {
//                ecoSystem.verify("/dogu")
//            }

            if (gitflow.isReleaseBranch()) {
                String releaseVersion = git.getSimpleBranchName()

                stage('Finish Release') {
                    gitflow.finishRelease(releaseVersion)
                }

                stage('Push Dogu to registry') {
                    ecoSystem.push("/dogu")
                }

                stage('Add Github-Release') {
                    github.createReleaseWithChangelog(releaseVersion, changelog)
                }
            }
        } finally {
            stage('Clean') {
                ecoSystem.destroy()
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