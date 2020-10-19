# DevOps2020

# 项目背景:


依托于 AWS EC2 基于 Dockerfile + Jenkinsfile + EKS   搭建golang web 应用 交付流水线


![image](http://processon.com/chart_image/5f3e59245653bb06f2dd122c.png)


# 整体说明：

- 需求、迭代、任务、缺陷管理，看板、测试计划

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

  脚本式Jenkinsfile 、Blue Ocean 、人工确认
  
  http://52.82.121.46:9000/

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

1. 升级k8s 版本

把k8s 从1.17 升级到1.18 

先升级master，再升级worker


强制升级后，worker的IP地址会变化，原来的worker 不可用，docker配置文件需要重新配置harbor地址


升级后ingress 不可用 ，ingress-nginx 命名空间无法删除 修复
```
kubectl  get ns  ingress-nginx  -o json > ingress.json

把此字段置空

"finalizers": [
            "finalizers.kubesphere.io/namespaces"  #记得要完全删除字段中的内容
        ],

kubectl proxy  --address='127.0.0.1'   --port=8001   --accept-hosts='^localhost$,^127\.0\.0\.1$,^\[::1\]$'

curl -k -H "Content-Type: application/json" -X PUT --data-binary @ingress.json http://127.0.0.1:8001/api/v1/naespaces/ingress-nginx/finalize



升级kubectl
curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.18.8/2020-09-18/bin/linux/amd64/kubectl

参考文档：https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html


kubectl apply  -f zuoguocai-nginx-ingress.yaml
```





2. 生成连接集群的 kubeconfig 
```
 2 worker 


aws --version

aws configure  先写自己账号里的手动生成的Access Key，区域不填

AWS Access Key ID [None]: AKIAX2GTKKBQXT322
AWS Secret Access Key [None]: TATT9l4xJvg92Rk0ZZqyk1np
Default region name [None]: 
Default output format [None]: 

aws eks  --region   cn-northwest-1 update-kubeconfig   --name zuoguocai-eks-master
 
ls /root/.kube/config 

cd /usr/local/bin
curl -o kubectl https://amazon-eks.s3.cn-north-1.amazonaws.com.cn/1.17.9/2020-08-04/bin/linux/amd64/kubectl
chmod +x  kubectl 

vi /etc/profile
export KUBE=/usr/local/bin
export PATH=$PATH:$KUBE
source /etc/profile

```


## nginx  ingress  controller 

```
k8s.gcr.io/ingress-nginx/controller:v0.35.0 改为 ghcr.io/zuoguocai/ingress-nginx/controller:v0.35.0

参考 https://github.com/kubernetes/ingress-nginx/blob/master/deploy/static/provider/baremetal/deploy.yaml

更改副本数、网络、默认端口、时区


timedatectl set-timezone  Asia/Shanghai

kubectl create secret tls zuogucoai-secret   --cert=zuoguocai.xyz.cert  --key=zuoguocai.xyz.key   -n devops

kubectl  exec  -it  ingress-nginx-controller-798c579896-ctrg8    /bin/bash    -n ingress-nginx


参考文档：https://github.com/kubernetes/ingress-nginx/issues/4857

```

## 服务发布
```

https://devops2020.zuoguocai.xyz:11443/



基于header 灰度

chrome 插件 Modheader


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

## jenkins 搭建


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


filebeat install

```
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-7.9.2-x86_64.rpm


yum  -y install filebeat-7.9.2-x86_64.rpm 
vim /etc/filebeat/filebeat.yml

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
