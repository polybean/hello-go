import java.text.SimpleDateFormat

def props
def label = "jenkins-slave-${UUID.randomUUID().toString()}"
currentBuild.displayName = new SimpleDateFormat("yy.MM.dd").format(new Date()) + "-" + env.BUILD_NUMBER

podTemplate(
  label: label,
  namespace: "build", // Not allowed with declarative
  serviceAccount: "build",
  yaml: """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: helm
    image: vfarcic/helm:2.9.1
    command: ["cat"]
    tty: true
    volumeMounts:
    - name: build-config
      mountPath: /etc/config
  - name: kubectl
    image: vfarcic/kubectl
    command: ["cat"]
    tty: true
  - name: golang
    image: golang:1.12
    command: ["cat"]
    tty: true
  volumes:
  - name: build-config
    configMap:
      name: build-config
"""
) {
  node(label) {
    stage("Build") {
      container("helm") {
        sh "cp /etc/config/build-config.properties ."
        props = readProperties interpolate: true, file: "build-config.properties"
      }
      node("docker") { // Not allowed with declarative
        checkout scm
        k8sBuildBetaImage(props.image)
      }
    }
    stage("Functional Test") {
      try {
        container("helm") {
          checkout scm
          k8sHelmUpgradeBeta(props.project, props.domain, "--set replica.app=2 --set replica.db=1")
        }
        container("kubectl") {
          k8sRolloutBeta(props.project)
        }
        container("golang") {
          k8sFuncTestGolang(props.project, props.domain, "src", true)
        }
      } catch(e) {
          error "Failed functional tests"
      } finally {
        container("helm") {
          k8sDeleteBeta(props.project)
        }
      }
    }
    if ("${BRANCH_NAME}" == "master") {
      stage("Release") {
        node("docker") {
          k8sPushImage(props.image)
        }
        container("helm") {
          k8sPushChart(props.project, props.chartVer, "${props.cmSub}.${props.domain}")
        }
      }
      stage("Deploy") {
        try {
          container("helm") {
            k8sHelmUpgrade(props.project, "${props.project}.${props.domain}")
          }
          container("kubectl") {
            k8sRollout(props.project)
          }
          container("golang") {
            k8sProdTestGolang("${props.project}.${props.domain}", "src", true)
          }
        } catch(e) {
          container("helm") {
            k8sRollback(props.project)
          }
        }
      }
    }
  }
}
