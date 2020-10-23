# DevOps2020

# 项目背景:


依托于 AWS EC2 基于 Dockerfile + Jenkinsfile + EKS   搭建golang web 应用 交付流水线


![image](http://processon.com/chart_image/5f3e59245653bb06f2dd122c.png)


# 整体说明：

- 需求、迭代、任务、缺陷管理，看板、测试计划（功能、性能、接口、自动化测试...）、提测

  我这里使用的阿里云效（https://devops.aliyun.com/） ，也可以使用 jira和禅道



- 源代码管理github 

  项目仓库 ( https://github.com/ZuoGuocai/DevOps2020.git)

  git tag 管理

  合并分支请求  Github Pull  Requests 
  
  Code Review (https://github.com/features/code-review/)



<!--

Gitlab Merge Request
webhook
-->



- 单元测试、静态扫描

  go test ，
  golangci-lint
  
- 编译构建

  go build, Dockerfile
  
- 制品库harbor

  https://harbor.zuoguocai.xyz:4443/
  
  修改每个EKS worker 的 docker配置/etc/docker/daemon.json，加入制品库地址

- CICD jenkins pipeline 

  脚本式Jenkinsfile 、参数化构建、动态构建、Blue Ocean 、人工确认、钉钉通知
  
  http://52.82.121.46:9000/ （ec2单机版）
  
  http://myci.zuoguocai.xyz:11180/ （k8s版）

- 部署发布EKS

  基于nginx ingress controller 动静分离
  
  基于nginx ingress controller header 的灰度 +  客户端chrome插件Modheader

- 度量    

1. 监控 

   监控宝

2. 日志

   Elastic cloud Kibana、Elasticsearch、Filebeat

3. 请求、性能分析

   Elastic cloud APM-Golang




<!-- 动态构建 、自动化测试、 度量-apm、日志、监控、链路-->

---

# 功能点：

1. AI学习导航

   https://devops2020.zuoguocai.xyz:11443/


2. ipcat web 工具， 在pod里获取客户端真实源地址工具(负载均衡反向代理需要开启 支持proxy protocol 协议)

   https://devops2020.zuoguocai.xyz:11443/ipcat


3. 大赛彩蛋 

   https://devops2020.zuoguocai.xyz:11443/caidan/




> 通过ingress 做url路由分发

> caidan、daohang、live2d 等静态资源放在nginx里托管，

> ipcat 处理动态请求












# 运维记录:





## EKS 


 我创建了一个集群，一个计算节点组 ，包含2 个 worker 节点

1. 升级k8s 版本

把k8s 从1.17 升级到1.18 

先升级master，再升级worker

通过观察 kubectl get  nodes 变化，查看进行中的情况

强制升级后，worker的IP地址会变化，原来的worker 不可用，docker配置文件需要重新配置harbor地址


升级后ingress 不可用 ，ingress-nginx 命名空间无法删除 修复
```
kubectl  get ns  ingress-nginx  -o json > ingress.json

修改ingress.json文件， 把此字段finalizers的值置空

"finalizers": [
            "finalizers.kubesphere.io/namespaces"  #记得要完全删除字段中的内容
        ],

kubectl proxy  --address='127.0.0.1'   --port=8001   --accept-hosts='^localhost$,^127\.0\.0\.1$,^\[::1\]$'

curl -k -H "Content-Type: application/json" -X PUT --data-binary @ingress.json http://127.0.0.1:8001/api/v1/naespaces/ingress-nginx/finalize



升级kubectl
curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.18.8/2020-09-18/bin/linux/amd64/kubectl

参考文档：https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html

重新部署nginx ingress controller
kubectl apply  -f zuoguocai-nginx-ingress.yaml


更改 DNS A记录

```





2. EC2 上生成连接集群的 kubeconfig 
```



在web控制台 右上角自己账号里 我的安全凭证--AWS IAM凭证---创建访问密钥 里手动生成Access Key ID，Access Key

aws 命令 AWS 官方镜像 Amazon Linux 2 AMI (HVM) 默认已经安装，如果是其他镜像需要手动安装

aws --version

配置 AWS CLI 凭证,这样我们就可以用在命令行管理AWS 服务了

aws configure ，填入刚才在web里生成的 Acess Key ID，Access Key，  区域不填

AWS Access Key ID [None]: AKIAX2GTKKBQXT322
AWS Secret Access Key [None]: TATT9l4xJvg92Rk0ZZqyk1np
Default region name [None]: c
Default output format [None]: 


配置Kubeconfig，用于kubectl 通过apiserver操作EKS 集群，中国宁夏 对应 cn-northwest-1，zuoguocai-eks-master 为web控制台创建eks的集群名称

aws eks  --region   cn-northwest-1 update-kubeconfig   --name zuoguocai-eks-master
 
ls /root/.kube/config 


下载kubectl ，版本和集群k8s版本要一致
cd /usr/local/bin
curl -o kubectl https://amazon-eks.s3.cn-north-1.amazonaws.com.cn/1.17.9/2020-08-04/bin/linux/amd64/kubectl
chmod +x  kubectl 

配置kubectl 全局环境变量
vi /etc/profile
export KUBE=/usr/local/bin
export PATH=$PATH:$KUBE
source /etc/profile


参考文档：https://docs.aws.amazon.com/zh_cn/eks/latest/userguide/getting-started-console.html

```

## nginx  ingress  controller 

```


参考 https://github.com/kubernetes/ingress-nginx/blob/master/deploy/static/provider/baremetal/deploy.yaml

更改副本数或者使用DamonSet的方式、网络、默认端口、时区、image仓库地址

spec:
  replicas: 2

spec:
  hostNetwork: true
  dnsPolicy: ClusterFirstWithHostNet

image: harbor.zuoguocai.xyz:4443/devops/ingress-nginx/controller@sha256:51b3966f02453315e7b4cbd04f20b83be73f76aad02dc6207f8d9ffac6bf5c7b
         
- --http-port=11180
- --https-port=11443
          
- name: TZ
  value: Asia/Shanghai


80，443，8080 无法使用，所以这里把80 改为了11180，443 改为了11443，但通过kubectl get ingress -A  
查看的话还是显示为80，443，在eks上可以查看到11180，11443 已经监听端口，通过端口可以正常访问，
实际生产环境不建议改端口，备案就行了

更改后文件为zuoguocai-nginx-ingress.yaml


如果无法获取该镜像，可以使用我的镜像仓库里的镜像
k8s.gcr.io/ingress-nginx/controller:v0.35.0 改为 ghcr.io/zuoguocai/ingress-nginx/controller:v0.35.0


timedatectl set-timezone  Asia/Shanghai

kubectl create secret tls zuogucoai-secret   --cert=zuoguocai.xyz.cert  --key=zuoguocai.xyz.key   -n devops

kubectl  exec  -it  ingress-nginx-controller-798c579896-ctrg8    /bin/bash    -n ingress-nginx


参考文档：https://github.com/kubernetes/ingress-nginx/issues/4857

```

## 服务发布
```

https://devops2020.zuoguocai.xyz:11443/



基于header 灰度


ingress deployment 配置添加

  annotations:
    kubernetes.io/ingress.class: nginx  # nginx=nginx-ingress| qcloud=CLB ingress
    #kubernetes.io/ingress.subnetId: subnet-xxxxxxxx   # if qcloud, should give subnet
    nginx.ingress.kubernetes.io/canary: "true"
    nginx.ingress.kubernetes.io/canary-by-header: "location"
    nginx.ingress.kubernetes.io/canary-by-header-value: "shenzhen"



chrome 插件 Modheader 添加 location: shenzhen




参考文档: https://help.coding.net/docs/best-practices/cd/nginx-ingress.html

```

## 回滚

```
kubectl get deployment -n devops

kubectl -n devops rollout history deployment/kustomize-getrealip  --revision=1

kubectl -n devops rollout undo deployment/kustomize-getrealip --to-revision=1

```




## 扩容

```

- 手动

kubectl scale deployment nginx-deployment --replicas=10

- 自动

kubectl autoscale deployment nginx-deployment --min=10 --max=15 --cpu-percent=80


- Pause/Resume

当 Deployment 的 .spec.paused = true 时，任何更新都不会被触发 rollout。通过如下命令设置 Deployment 为 paused：

kubectl -n <namespace> rollout pause deployment/<deployment-name>

还原：

kubectl -n <namespace> rollout resume deploy/<deployment-name




参考
https://whypro.github.io/hexo-blog/20180301/Kubernetes-%E6%9C%8D%E5%8A%A1%E7%81%B0%E5%BA%A6%E5%8D%87%E7%BA%A7%E6%9C%80%E4%BD%B3%E5%AE%9E%E8%B7%B5/

```


## dnspod


这里没有负载均衡，所以 A 记录直接对应到两个worker的公网IP上


## git

```
https://github.com/ZuoGuocai/DevOps2020.git

```
## jenkins 容器化+动态构建


```
- master

安装yaml 文件 见 install_jenkins_in_k8s 文件夹

http://myci.zuoguocai.xyz:11180/
kubectl  exec -it    jenkins-0  cat /var/jenkins_home/secrets/initialAdminPassword  -n default

插件：
Build with Parameters，Blue Ocean，Dingtalk，Kubernetes plugin
jienkins 中的 Kubernetes 插件：Jenkins 在 Kubernetes 集群中运行动态代理 插件介绍：https://github.com/jenkinsci/kubernetes-plugin
http://myci.zuoguocai.xyz:11180/restart 

配置Kubernetes plugin 插件
系统管理--系统配置--Cloud--Add a new cloud--kubernetes

名称              kubernetes
Kubernetes 地址	https://kubernetes.default.svc.cluster.local 或https://kubernetes.default 
Jenkins 地址	http://172.31.14.217:30006
Jenkins 通道	172.31.14.217:31400

 测试连接


宽松的 RBAC 权限 
可以使用 RBAC 角色绑定在多个场合使用宽松的策略。
警告：
下面的策略允许 所有 服务帐户充当集群管理员。 容器中运行的所有应用程序都会自动收到服务帐户的凭据， 可以对 API 执行任何操作，包括查看 Secrets 和修改权限。 这个策略是不被推荐的。

kubectl create clusterrolebinding permissive-binding \
  --clusterrole=cluster-admin \
  --user=admin \
  --user=kubelet \
  --group=system:serviceaccounts






- slave

使用官方slave镜像 jenkinsci/jnlp-slave


测试pipeline

podTemplate(label: 'jenkins-slave', cloud: 'kubernetes', containers: [
    containerTemplate(
        name: 'jnlp', 
        image: "jenkinsci/jnlp-slave"
    ),
  ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
    hostPathVolume(mountPath: '/usr/bin/docker', hostPath: '/usr/bin/docker')
  ],
) 
{
  node("jenkins-slave"){
      // 第一步
      stage('拉取代码'){
        echo "123"
    }
  }
}






```

自定义slave镜像

```
wget  http://myci.zuoguocai.xyz:11180//jnlpJars/slave.jar
wget  http://myci.zuoguocai.xyz:11180//jnlpJars/agent.jar


[root@ip-172-31-39-226 jenkins-slave]# md5sum agent.jar
d866f0b482db94f38e49b26b465d5db5  agent.jar
[root@ip-172-31-39-226 jenkins-slave]# md5sum slave.jar 
d866f0b482db94f38e49b26b465d5db5  slave.jar

slave.jar 和agent.jar 应该是同样的一个包，使用其中一个就行


见 install_jenkins_in_k8s/jenkins-slave 文件夹

docker build -t harbor.zuoguocai.xyz:4443/devops/jenkins-slave-jdk:1.8  .
docker push harbor.zuoguocai.xyz:4443/devops/jenkins-slave-jdk:1.8


测试

podTemplate(label: 'jenkins-slave', cloud: 'kubernetes', containers: [
    containerTemplate(
        name: 'jnlp', 
        image: "harbor.zuoguocai.xyz:4443/devops/jenkins-slave-jdk:1.8"
    ),
  ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
    hostPathVolume(mountPath: '/usr/bin/docker', hostPath: '/usr/bin/docker')
  ],
) 
{
  node("jenkins-slave"){
      // 第一步
      stage('拉取代码'){
        echo "123"
    }
  }
}

```


参考文档： 

https://github.com/jenkinsci/kubernetes-plugin/blob/fc40c869edfd9e3904a9a56b0f80c5a25e988fa1/src/main/kubernetes/jenkins.yml

https://github.com/jenkinsci/docker-inbound-agent/blob/master/jenkins-agent

https://aws.amazon.com/cn/blogs/china/base-on-jenkins-create-kubernetes-on-aws-ci-cd-tube/

https://plugins.jenkins.io/kubernetes/

https://kubernetes.io/zh/docs/reference/access-authn-authz/rbac/#%E5%AE%BD%E6%9D%BE%E7%9A%84-rbac-%E6%9D%83%E9%99%90

https://github.com/jenkinsci/docker-inbound-agent

## jenkins 单机搭建


```

wget http://ftp-chi.osuosl.org/pub/jenkins/war/2.256/jenkins.war

yum  -y  install java  java-devel

java -jar  jenkins.war   --httpPort=9000  --daemon


http://52.82.121.46:9000/



```

## Kaniko

<!-- https://blog.ihypo.net/15487483292659.html >

## kustomization

<!-- argo-cd支持多种配置管理/模板工具(Kustomize、Helm、Ksonnet、Jsonnet、plain-YAML) -->


## argocd

```

kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

argocd rollout
kubectl apply -n argo-rollouts -f https://raw.githubusercontent.com/argoproj/argo-rollouts/stable/manifests/install.yaml



kubectl create secret tls zuogucoai-secret   --cert=zuoguocai.xyz.cert  --key=zuoguocai.xyz.key   -n argocd

kubectl apply -f argocd-ingress.yaml

 VERSION=$(curl --silent "https://api.github.com/repos/argoproj/argo-cd/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
 sudo curl --silent --location -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64
 curl -LO https://github.com/argoproj/argo-rollouts/releases/latest/download/kubectl-argo-rollouts-linux-amd64
 mv ./kubectl-argo-rollouts-linux-amd64 /usr/local/bin/kubectl-argo-rollouts
 kubectl argo rollouts


kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2


argocd login     myci.zuoguocai.xyz:11443  --grpc-web


argocd app create colorapi --repo https://github.com/particuleio/demo-concourse-flux.git --path deploy --dest-server https://kubernetes.default.svc --dest-namespace default


argocd app sync colorapi


https://myci.zuoguocai.xyz:11443/

参考文档：https://argoproj.github.io/argo-cd/operator-manual/ingress/
```

##  canary

https://blog.51cto.com/liuqunying/1925463 

https://www.cnblogs.com/xiaoqi/p/ingress-nginx-canary.html

https://troy.wang/docs/kubernetes/posts/argo-rollouts-support-traefik/

https://qiita.com/tomozo6/items/1bfc65a86a528f63d205



## tekon

```

Install Tekton Pipelines
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml

Install Tekton CLI (tkn)
curl -LO https://github.com/tektoncd/cli/releases/download/v0.7.1/tkn_0.7.1_Linux_x86_64.tar.gz
 
# Change destination directory as needed
tar xvzf tkn_0.7.1_Linux_x86_64.tar.gz -C ~/bin

Install Tekton dashboard 仪表盘
kubectl apply --filename https://github.com/tektoncd/dashboard/releases/download/v0.5.1/tekton-dashboard-release.yaml









```


## docker

```
yum  -y install docker

cat /etc/docker/daemon.json 
{
  "bridge": "none",
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "10"
  },
  "live-restore": true,
  "max-concurrent-downloads": 10，
  "registry-mirrors": [
        "https://docker.mirrors.ustc.edu.cn",
        "https://registry.docker-cn.com"
    ],
  "insecure-registries": [
    "https://harbor.zuoguocai.xyz:4443"
      ],
  "graph": "/data/docker"
}

 

systemctl start docker

systemctl enable docker


```


## docker-compose 安装

```

wget https://github.com/docker/compose/releases/download/1.27.2/docker-compose-Linux-x86_64


mv docker-compose-Linux-x86_64   docker-compose

chmod +x docker-compose 

mv docker-compose  /usr/bin

```

<!--ignore-preflight-errors=NumCPU  失败  https://juejin.im/entry/6844903781314887694 -->

## harbor


```
wget https://github.com/goharbor/harbor/releases/download/v2.1.0/harbor-offline-installer-v2.1.0.tgz

https://harbor.zuoguocai.xyz:4443/




docker login  https://harbor.zuoguocai.xyz:4443


docker tag  nginx:latest  harbor.zuoguocai.xyz:4443/devops/nginx:v1

docker push harbor.zuoguocai.xyz:4443/devops/nginx:v1

```
## JFog

使用 Artifactory docker https://www.kdocs.cn/l/s9epu3a2C?f=501  


使用 Artifactory Maven 仓库：https://www.kdocs.cn/l/sQpfp1M74?f=501



## Elastic cloud

https://cloud.elastic.co/home

Cloud/Deployment/集群名 里有CloudID

web控制台管理

> Manage--Reset password 里给elastic user 设置访问密码

> Management--Stack Management--Data--Index Management

> Management--Stack Management--Kibana--Index Patterns

> Management--Stack Management--Index Management--Templates--Create template



1. 日志

   修改filebeat 配置文件 把 cloud.id和cloud.auth（elastic 用户名和密码） 填入


- pod级别的采集(filebeat + sidecar)



```
通过 共享日志目录采集

见 filebeat-sidecar 文件夹


```
  参考文档：

   https://www.docker.elastic.co/r/beats/filebeat:7.9.2

   https://github.com/elastic/beats/blob/master/deploy/kubernetes/filebeat-kubernetes.yaml


- 节点级别的采集

 filebeat install

```
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-7.9.2-x86_64.rpm


yum  -y install filebeat-7.9.2-x86_64.rpm 

vim /etc/filebeat/filebeat.yml
# 6.3以前是 filebeat.prospectors  以后是 filebeat.inputs  
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/lib/docker/containers/*/*.log
cloud.id: "i-o-optimized-deployment:ZWFzdHVzMi5henVyZS5lbGFzdGljLWNsb3VkLmNvbSRhMWU2ZTNmN2I3ZmE0NGU4YTc4MDFjNjZmOGIzNGUxNSQ2MmM4OGI1YTQ5ZTM0NTFlYWFjYjMxMDY1OGI0ODNkMQ=="
cloud.auth: "elastic:PCR0smjgvI1PlfSwfODg5mny"

sudo filebeat modules enable elasticsearch
sudo filebeat setup
sudo systemctl start  filebeat 


```

  参考文档： https://www.elastic.co/guide/en/beats/filebeat/current/configure-cloud-id.html

2. APM

点击 Quick Link Kibana--- Observability（Add APM)---Go--Configure the agent 里可以找到ELASTIC_APM_SERVER_URL，ELASTIC_APM_SECRET_TOKEN

```

作为环境变量引入Dockerfile

ENV ELASTIC_APM_SERVICE_NAME=ipcat
# Set custom APM Server URL (default: http://localhost:8200)
ENV ELASTIC_APM_SERVER_URL=https://8b06fec588334601ba91e8ad7fe235c3.apm.eastus2.azure.elastic-cloud.com:443
# Use if APM Server requires a token
ENV ELASTIC_APM_SECRET_TOKEN=JOFwFHBYdXzAbIMUYP

ipcat.go import apm agent 的包

import (
	"net/http"

	"go.elastic.co/apm/module/apmhttp"
)

func main() {
	mux := http.NewServeMux()
	...
	http.ListenAndServe(":5000", apmhttp.Wrap(mux))
}

```




## 监控宝

https://monitoring.cloudwise.com/

## grafana tanka
```

https://tanka.dev/tutorial/jsonnet

sudo curl -Lo /usr/local/bin/tk https://github.com/grafana/tanka/releases/latest/download/tk-linux-amd64
sudo curl -Lo /usr/local/bin/jb https://github.com/jsonnet-bundler/jsonnet-bundler/releases/latest/download/jb-linux-amd64
```


<!--
## jenkins动态构建
https://mp.weixin.qq.com/s/ODLVHXRHCnNiWrSDfpoIqQ -->







<!--

## arm 
```
apiVersion: v1
kind: Service
metadata:
  name: my-service
  namespace: martin-demo
  labels:
    app: my-app
spec:
  selector:
    app: my-app
  ports:
    - protocol: TCP
      port: 8081
      targetPort: 80
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
  namespace: martin-demo
  labels:
    app: my-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: beta.kubernetes.io/arch
                operator: In
                values:
                - amd64
                - arm64
      containers:
      - name: nginx
        image: nginx:1.19.2
        ports:
        - containerPort: 80

```
-->
