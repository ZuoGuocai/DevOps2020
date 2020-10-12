/*
Author: zuoguocai@126.com
Description: Test successfully  in Jenkins 2.256 
*/





pipeline {
    agent any
	
   // tools {
      //  go 
    //}

       
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
                 
                 git url: "https://github.com/ZuoGuocai/DevOps2020.git"
                 
                
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
         

            steps {

                  dir('') {
                    // 删除之前构建镜像
                    sh "docker image prune  -a --force --filter  label='maintainer=zuoguocai@126.com'"
                    // build镜像
                    
                    //echo build_tag
                    
                   sh "docker build -t harbor.zuoguocai.xyz:4443/devops/ipcat:${build_tag} ."
                    // 登录镜像仓库
                    sh "docker login -u ${HARBOR_CREDS_USR} -p ${HARBOR_CREDS_PSW} harbor.zuoguocai.xyz:4443"
                    // 推送镜像到镜像仓库
                   sh "docker push harbor.zuoguocai.xyz:4443/devops/ipcat:${build_tag}"
                   
                   }
    
            }
          
        }
        
            
        stage('部署到k8s集群') {

                input {
                message "测试环境"
                ok "提交."
                submitter ""
                parameters {
                    string(name: 'PASSWD', defaultValue: '', description: '请输入密码开始部署')
                }
                }
              
               steps{

                        script{
                              if (PASSWD == HARBOR_CREDS_PSW) {
                                 echo "start release to test"
                        
                                 sh "sed  's/<IMG_TAG>/${build_tag}/g' ipcat-canary.tmpl   > ipcat.yaml"
                                 sh "kubectl apply  -f  ipcat.yaml"
                                 //sh "kubectl get pods -n devops"
                              } else {
                                 echo '密码错误,部署失败'
                              }
                        }
                 }
        }
           
	    
	
   // K8S紧急时回滚
    stage('Rollback to k8s') {
          when { environment name: 'action', value: 'rollback' }
          steps {
            echo "k8s images is rolled back! " 
            sh '''
             kubectl rollout undo deployment/tomcat-dpm  -n default
            '''
          } 
       }  

	    
	    
	    
       
        

    }
}
