FROM golang:1.13.4-stretch 

RUN apt-get update 
RUN apt-get install -y remake

ENV PATH $PATH:/usr/local/go/bin
ENV PATH $PATH:/usr/local/
ENV GOROOT /usr/local/go
ENV GOPATH=/home/jenkins/agent
ENV PATH $PATH:$GOPATH/bin

# USER jenkins
WORKDIR /home/jenkins/agent

# jx
ENV JX_VERSION 2.0.1142
RUN curl -f -L https://github.com/jenkins-x/jx/releases/download/v${JX_VERSION}/jx-linux-amd64.tar.gz | tar xzv && \
mv jx /usr/bin/

# Docker
ENV DOCKER_VERSION 17.12.0
RUN curl -f https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION-ce.tgz | tar xvz && \
  mv docker/docker /usr/bin/ && \
  rm -rf docker

# helm
ENV HELM_VERSION 2.11.0
RUN curl -f https://storage.googleapis.com/kubernetes-helm/helm-v${HELM_VERSION}-linux-amd64.tar.gz  | tar xzv && \
  mv linux-amd64/helm /usr/bin/ && \
  mv linux-amd64/tiller /usr/bin/ && \
  rm -rf linux-amd64

# jx-release-version
ENV JX_RELEASE_VERSION 1.0.10
RUN curl -f -o ./jx-release-version -L https://github.com/jenkins-x/jx-release-version/releases/download/v${JX_RELEASE_VERSION}/jx-release-version-linux && \
  mv jx-release-version /usr/bin/ && \
  chmod +x /usr/bin/jx-release-version

# exposecontroller
ENV EXPOSECONTROLLER_VERSION 2.3.34
RUN curl -f -L https://github.com/fabric8io/exposecontroller/releases/download/v$EXPOSECONTROLLER_VERSION/exposecontroller-linux-amd64 > exposecontroller && \
  chmod +x exposecontroller && \
  mv exposecontroller /usr/bin/

# skaffold
ENV SKAFFOLD_VERSION 1.0.0
RUN curl -f -Lo skaffold https://github.com/GoogleCloudPlatform/skaffold/releases/download/v${SKAFFOLD_VERSION}/skaffold-linux-amd64 && \
  chmod +x skaffold && \
  mv skaffold /usr/bin

# kubectl
RUN curl -f -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
  chmod +x kubectl && \
  mv kubectl /usr/bin/  

CMD ["go","version"]