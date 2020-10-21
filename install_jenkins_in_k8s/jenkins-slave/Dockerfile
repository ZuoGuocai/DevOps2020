##agent.jar通过url下载  wget  http://myci.zuoguocai.xyz:11180//jnlpJars/agent.jar
FROM centos:7
LABEL maintainer="zuoguocai@126.com"  version="2.0" description="jenkins-slave"



RUN yum install -y java-1.8.0-openjdk maven curl git libtool-ltdl-devel && \ 
    yum clean all && \
    rm -rf /var/cache/yum/* && \
    mkdir -p /usr/share/jenkins

COPY agent.jar /usr/share/jenkins/agent.jar 
COPY settings.xml /etc/maven/settings.xml
COPY jenkins-agent /usr/local/bin/jenkins-agent
RUN chmod +x /usr/local/bin/jenkins-agent &&\
    ln -s /usr/local/bin/jenkins-agent /usr/local/bin/jenkins-slave

ENTRYPOINT ["jenkins-agent"]
