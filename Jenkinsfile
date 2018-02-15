#!/usr/bin/groovy
@Library('k8s-jenkins-tools') _
import com.ultimaker.Slug

def credentialsFile = 'gcloud-jenkins-service-account'
def projectName = 'um-website-193311'
def clusterName = 'ultimaker-dev'
def zone = 'europe-west4-b'

def slugify = new Slug()
def branchSlug = slugify.slug(env.BRANCH_NAME)

podTemplate(label: 'jenkins-ns-cleaner-pipeline', inheritFrom: 'default', containers: [
  containerTemplate(name: 'golang', image: 'golang:1.9.2', ttyEnabled: true, command: 'cat')
], volumes: [
  hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
  hostPathVolume(mountPath: '/usr/bin/docker', hostPath: '/usr/bin/docker')
]) {
  node('jenkins-ns-cleaner-pipeline') {
    checkout scm

    stage('install dependencies') {
      container('golang') {
        sh '''
          mkdir -p $GOPATH/src/github.com/Ultimaker
          ln -s `pwd` $GOPATH/src/github.com/Ultimaker/k8s-ns-cleaner
          cd $GOPATH/src/github.com/Ultimaker/k8s-ns-cleaner &&
          go get -v ./
        '''
      }
    }

    stage('build binary') {
      container('golang') {
        sh '''
          cd $GOPATH/src/github.com/Ultimaker/k8s-ns-cleaner &&
          CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o bin/ns-cleaner main.go
        '''
      }
    }

    stage('build image') {
      container('jnlp') {
        sh "docker build -t eu.gcr.io/um-website-193311/kube-system/ns-cleaner:${branchSlug} ."
      }
    }

    stage('authenticate gcloud') {
      container('jnlp') {
        withCredentials([file(credentialsId: credentialsFile, variable: 'GCLOUD_KEY_FILE')]) {
          sh "gcloud auth activate-service-account --key-file=${GCLOUD_KEY_FILE}"
          sh "gcloud config set project ${projectName}"
          sh "gcloud container clusters get-credentials ${clusterName} --zone ${zone} --project ${projectName}"
        }
      }
    }

    stage('push images') {
      container('jnlp') {
        sh "gcloud docker -- push eu.gcr.io/um-website-193311/kube-system/ns-cleaner:${branchSlug}"
      }
    }

    stage('deploy') {
      container('jnlp') {
        sh "sed -i s#_BRANCH_SLUG_#${branchSlug}# k8s/cronjob.yaml"
        sh "kubectl apply -f k8s/cronjob.yaml"
      }
    }

  }
}
