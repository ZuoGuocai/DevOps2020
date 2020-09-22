# DevOps2020


# 安装记录

# DevOps2020


#### 项目背景

依托于 AWS ECS  基于 Dockerfile + Jenkinsfile + EKS   搭建交付流水线



#### EKS 


```
1master + 2 worker 


aws --version
aws configure
AWS Access Key ID [None]: AKIAX2GTKKBQXT322
AWS Secret Access Key [None]: TATT9l4xJvg92Rk0ZZqyk1np
Default region name [None]: 
Default output format [None]: 

aws eks  --region   cn-northwest-1 update-kubeconfig   --name zuoguocai-eks-master

 
ls /root/.kube/config 


- nginx  ingress  controller 

参考 https://github.com/kubernetes/ingress-nginx/blob/master/deploy/static/provider/baremetal/deploy.yaml

更改副本数、网络、默认端口、时区


timedatectl set-timezone  Asia/Shanghai

kubectl create secret tls zuogucoai-secret   --cert=zuoguocai.xyz.cert  --key=zuoguocai.xyz.key   -n devops

kubectl  exec  -it  ingress-nginx-controller-798c579896-ctrg8    /bin/bash    -n ingress-nginx

```

#### harbor


```
wget https://github.com/goharbor/harbor/releases/download/v2.1.0/harbor-offline-installer-v2.1.0.tgz

https://harbor.zuoguocai.xyz:4443/




docker login  https://harbor.zuoguocai.xyz:4443


docker tag  nginx:latest  harbor.zuoguocai.xyz:4443/devops/nginx:v1

docker push harbor.zuoguocai.xyz:4443/devops/nginx:v1

```

#### dnspod



#### git

```
https://github.com/ZuoGuocai/DevOps2020.git

```

#### jenkins 搭建


```

wget http://ftp-chi.osuosl.org/pub/jenkins/war/2.256/jenkins.war

yum  -y  install java  java-devel

java -jar  jenkins.war   --httpPort=9000  --daemon




http://161.189.123.98:9000/




```
#### docker

```
yum  -y install docker

cat /etc/docker/daemon.json 
{
  "registry-mirrors": [
        "https://docker.mirrors.ustc.edu.cn",
        "https://registry.docker-cn.com"
    ],
  "insecure-registries": [
    "https://registry.xxx.com"
      ],
  "graph": "/data/docker"
}
 

systemctl start docker

systemctl enable docker


```


#### docker-compose 安装

```

wget https://github.com/docker/compose/releases/download/1.27.2/docker-compose-Linux-x86_64


mv docker-compose-Linux-x86_64   docker-compose

chmod +x docker-compose 

mv docker-compose  /usr/bin

```

<!--ignore-preflight-errors=NumCPU  失败  https://juejin.im/entry/6844903781314887694 -->


#### JFog


#### Elasticsearch

https://cloud.elastic.co/home
