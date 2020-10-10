/*
Author: zuoguocai@126.com
*/

pipeline {
    agent any
       
    environment {
    
        HARBOR_CREDS = credentials('jenkins-harbor-creds')
        BUILD_USER_ID = ""
        BUILD_USER = ""
        BUILD_USER_EMAIL = ""
        
    }
    
    stages {
    
    
        stage('准备环境变量'){
        
              steps {
              // 由插件user build vars 提供
              
               wrap([$class: 'BuildUser']) {
                   script {
                       BUILD_USER_ID = "${env.BUILD_USER_ID}"
                       BUILD_USER = "${env.BUILD_USER}"
                       BUILD_USER_EMAIL = "${env.BUILD_USER_EMAIL}"
                   }
				}
				// Test out of wrap
				echo "Build User ID: ${BUILD_USER_ID}"
				echo "Build User: ${BUILD_USER}"
				echo "Build User Email: ${BUILD_USER_EMAIL}"
            }
        
        }    
    
    
        stage('拉取代码') { // for display purposes
        
            steps{
                 // 清理工作区
                 step([$class: 'WsCleanup'])
                 // 拉取代码
                              
                 script {
                 
                 git credentialsId: '', url: "https://github.com/ZuoGuocai/DevOps2020.git"
                 
                
                  build_tag = sh(returnStdout: true, script: 'git describe --tags --always').trim()
                 //echo build_tag
                 
                 }
                
                
            }
        }
        
        stage('编译代码') {
            steps {
                sh 'go version'
                sh 'go build .'
            }
        }
        
 
         
        
        
        stage('构建镜像'){
         
            input {
                message "测试环境"
                ok "提交."
                submitter ""
                parameters {
                    string(name: 'PASSWD', defaultValue: '', description: '请输入密码开始部署')
                }
            }
         
            steps {
            
              script{
              if (PASSWD == HARBOR_CREDS_PSW) {
                echo "start build image"
                
                  dir('') {
                    // 删除之前构建镜像
                    sh "docker image prune -a --force  --filter 'label=ZuoGuocai'"
                    // build镜像
                    
                    //echo build_tag
                    
                   sh "docker build -t harbor.zuoguocai.xyz:4443/devops/ipcat:${build_tag} ."
                    // 登录镜像仓库
                    sh "docker login -u ${HARBOR_CREDS_USR} -p ${HARBOR_CREDS_PSW} harbor.zuoguocai.xyz:4443"
                    // 推送镜像到镜像仓库
                   sh "docker push harbor.zuoguocai.xyz:4443/devops/ipcat:${build_tag}"
                   
                   }
                
                
                 } else {
                     echo '密码错误,部署失败'
                  }
                }
            }
          
        }
        
            
        stage('部署到k8s集群') {
              
               steps{
                    sh "sed  's/<IMG_TAG>/${build_tag}/g' /opt/deploy/ipcat.temp   > /opt/deploy/ipcat.yaml"
                    sh "kubectl apply  -f /opt/deploy/getrealip.yaml"
                    sh "kubectl get pods -n devops"
                                       
                    }
                    
                }
                
                post {
                 success{
                    dingtalk (
                    // robot 为插件DingTalk配置后自动生成的id,在系统管理--系统配置--钉钉里找
                        robot: '14205481-9c8d-40d3-a667-95a91a09b33f',
                        type: 'MARKDOWN',
                        title: 'Jenkins pipeline构建通知',
                        text: [
                            "# <font color=#66CDAA>${env.BUILD_DISPLAY_NAME}构建成功 </font>",
                           '---',
                           "- 执行人: ${BUILD_USER}",
                           "- 邮箱: ${BUILD_USER_EMAIL}",
                           "- 作业: ${env.WORKSPACE}",
                          
                        ],
                        at: [
                          '13020038138'
                        ]
                    )
                }
            }
        }
        
       
        

    }
}
