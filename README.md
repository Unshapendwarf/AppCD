# AppCD  
>Feb.5.2020  
Feb.13.2020(updated)
---
## Concept  
AppCD를 통해 pipeline 관리 툴인 Argo Workflow와 git repo 기반 자동 배포 툴인 ArgoCD를 이용해 간단하게 정의한 config(App config)를 통해 어플리케이션을 지속적이고 선택적으로 배포한다.

---
## Based on Argo Workflow and Argo CD

You can study Argo Workflow at this url. See [Argo Workflow](https://argoproj.github.io/docs/argo/readme.html).
 Also, Argo CD at this url. see [Argo CD](https://argoproj.github.io/argo-cd/).

---
## Flow
argo workflow를 통해 작성한 이 CD 툴은 아래와 같은 흐름으로 진행된다.
1. argoCD credential: ArgoCD cluster 접근을 위해 ip, port, id, pw, token 등의 정보를 불러온다(2020년 2월 13일 기준으로 이 부분 미리 주어졌다 가정하고 진행)
2. Get "App config" file and current Application list(argoCD app).
    - App config is a newly defined config file for this AppCD, this is in a git repository which is periodically watched by AppCD.
    - We will get the current Application list from response of argocd API call.
3. Convert the data(from config, and application list) to json form. Also, get the difference between two given data.
4. Do one of these three actions;'create', 'delete', 'update' for each applications. It will be decided what kind of action should be performed by the difference found in above process.
5. Send ArgoCD API request to argoCD cluster, and do validation with response from cluster.

---
## How to use this pipeline(AppCD)
#### Requirement
AppCD를 이용하기 위해 아래와 같은 환경이 필요하다
- workflow를 정의한 파일(지금이 README.md가 작성된 tde repo의 workflow/workflow.yaml이다)
- argo workflow가 설치된 k8s 클러스터
- argoCD application을 배포할 argoCD k8s 클러스터
- application화 되어 배포될 서비스가 정의된 manifest(쉽게 말해 실제 배포를 원하는 helm chart, jsonnet 등등) 그리고 이 manifest를 포함하는 git repo  

결과적으로 최소한 **한개의 k8s 클러스터**(여기에 argocd, argo workflow를 설치하고 helm chart를 배포한다)와 **한개의 git repository**(여기에 config file, workflow.yaml, 배포할 helm chart들 있다)만 있어야 한다.

#### Pipeline 생성하고 실행하기
위 사항이 전부 준비되었다면 pipeline을 실행시키기 전에 배포하고자하는 helm chart에 대한 정보를 담은 config file을 수정해준다(configs/Autoconfig.yaml or Manualconfig.yaml참고)  
이후에 workflow 파일을 수정한다. 수정이 필요한 부분은 workflow전역변수로 사용가능한 parameter들이다(worklflow.yaml의 대문자 주석 참고)  
현재 config file들을 수정하지 않고 tde에 적힌 workflow를 그대로 돌려보면 192.168.48.12의 argoCD cluster에 3가지 application 을 배포한다. 직접 들어가서 확인해보자.  
접속 주소   
argocd -> 192.168.48.12:31410 (argoCD id: admin/ pw: 1222)  
argoworkflow -> 192.168.48.12:31367  

간단하게 argo workflow가 설치된 k8s 클러스터에서 배포해보자
```sh
$ cd ~
$ git clone https://tde.sktelecom.com/stash/scm/oreotools/appcd.git #id, pw 필요
$ cd appcd/
```
사용자에 맞게 configs/Autoconfig.yaml, configs/Manualconfig.yaml을 수정하자  
그리고 배포 환경에 맞게 workflow/workflow.yaml을 수정하자
이후 submit하면 파이프라인이 실행된다.
```sh
$ cd ~/appcd/
$ argo submit workflow/workflow.yaml
```
실행되는 과정, 결과는 위에 적힌 주소를 통해서 확인가능하다.
