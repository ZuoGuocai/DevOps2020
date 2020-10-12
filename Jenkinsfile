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
           
	    
	
stage('RollOut') {
      
      input {
        id 'ROLLOUT'
        message "是否快速回滚？"
        ok "确认"
        submitter ""
        parameters {
          choice(name: 'UNDO', choices: ['NO', 'YES'], description: '是否快速回滚？')
        }
      }
      
      
        steps {
    
          echo "Kubernetes快速回滚"
               script {
                 if ("${UNDO}" == 'YES') {
                   sh '''
                   # 快速回滚 - 回滚到最近版本
                   kubectl  rollout undo deployment ipcat -n devops
                   # 回滚到指定版本
                   # kubectl -n ${NAMESPACE} rollout undo deployment consume-deployment --to-revision=$(kubectl -n ${NAMESPACE} rollout history deployment consume-deployment | grep ${COMMIT_ID} | awk '{print $1}')
                   # kubectl -n ${NAMESPACE} rollout status deployment consume-deployment
                   '''
                 }
               }
        }
      
}
	    
	    
	    
       
        

    }
}
