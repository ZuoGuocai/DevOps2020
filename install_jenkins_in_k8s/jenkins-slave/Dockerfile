##slave.jar通过url下载  wget  http://myci.zuoguocai.xyz:11180//jnlpJars/slave.jar
FROM centos:7
LABEL maintainer="zuoguocai@126.com"  version="2.0" description="jenkins-slave"


RUN yum install -y java-1.8.0-openjdk maven curl git libtool-ltdl-devel && \ 
    yum clean all && \
    rm -rf /var/cache/yum/* && \
    mkdir -p /usr/share/jenkins

COPY slave.jar /usr/share/jenkins/slave.jar 
COPY start-jenkins-slave.sh /usr/bin/start-jenkins-slave.sh
COPY settings.xml /etc/maven/settings.xml
RUN chmod +x /usr/bin/start-jenkins-slave.sh

ENTRYPOINT ["start-jenkins-slave.sh"]

