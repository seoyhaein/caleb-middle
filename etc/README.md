# Caleb Mesos Environment

## 목적
메소스 (master, agent)를 빌드 후, 생성 삭제를 쉽게 할 수 있도록 AWS AMI 생성.

생성 된 master, agent ami를 통해 인스턴스 생성 후, 
1. master -> 
    1. MESOS_HOSTNAME=${HOSTNAME}, MESOS_ZK=zk://${MASTERIP}:2181/mesos 옵션 등록
    2. zookeeper.service, mesos-master.service 서비스 등록 후 docker, zookeeper, mesos-master 런
2. agent -> 
    1. MESOS_MASTER=zk://${MASTERIP}:2181/mesos 옵션 등록
    2. mesos-agent.service 서비스 등록 후 docker, mesos-agent 런 

## AMI 만들기
### centos 이미지 선택
![centos.png](./static/centos.png)
---
### 빠른 빌드를 위해 c5.4xlarge 선택(16코어)
![c5.4xlarge.png](./static/c5.4xlarge.png)
---
### master, agent 각각 빌드하기 위해 2개의 인스턴스 설정 
![instanceInfo.png](./static/instanceInfo.png)
### 종료 시 삭제 체크
![block2.png](./static/block2.png)
---
---

### 한 대의 노드에서 Master 빌드 스크립트 실행
* [메소스빌드 참조 URL](http://mesos.apache.org/documentation/latest/building/)
> ./install_master.sh  코어수
0. ICGPATH=/opt/ichthysGenomics, 마스터 런타임 환경파일 경로=/opt/ichthysGenomics/build/environments, systemd서비스파일 =/opt/ichthysGenomics/build/services
1. 기본 빌드 환경 패키지  설치
2. 주키퍼 설치
3. 메소스 설치
---
### 나머지 노드에서 Agent 빌드 스크립트 실행
* [메소스빌드 참조 URL](http://mesos.apache.org/documentation/latest/building/)
> ./install_agent.sh 코어수
0. ICGPATH=/opt/ichthysGenomics, 마스터 런타임 환경파일 경로=/opt/ichthysGenomics/build/environments, systemd서비스파일 =/opt/ichthysGenomics/build/services
1. awscli 를 위한 openssl, python 설치 
2. 메소스 설치

### 빌드 성공 후 AMI 생성
* [awscli URL](https://aws.amazon.com/ko/cli/)
> 외부에서 awscli 사용 
>> 각각 빌드 노드에서 해당 인스턴스ID(AWS에서 확인가능) 사용
>>> aws ec2 create-image --instance-id 인스턴스ID --name "[master|agent]" --no-reboot
>>>> aws ec2 콘솔에서 AMI 상태 available이면 해당 빌드 인스턴스 삭제, AMI로 인스턴스 생성 가능


## AMI로 서버실행하기 
### master ami 로 이미지 선택
![masterAMI.png](./static/masterAMI.png)
---
### master 인스턴스 설정(public subnet, 외부 IP 활성화)
![masterSubnet.png](./static/masterSubnet.png)
---
### master disk 설정(AMI 만든 인스턴스 디스크크기(50G로 설정했음) 이상)
![masterDisk.png](./static/masterDisk.png)
---
---
### agent ami 로 이미지 선택
![agentAMI.png](./static/agentAMI.png)
---
### agent 인스턴스 설정(private subnet, 외부 IP 비활성화)
![agentSubnet.png](./static/agentSubnet.png)
---
### agent disk 설정(AMI 만든 인스턴스 디스크크기(50G로 설정했음) 이상)
![masterDisk.png](./static/masterDisk.png)


## Master 실행하기 
* [master|agent 공통옵션 URL](http://mesos.apache.org/documentation/latest/configuration/master-and-agent/)
* [master 옵션 URL](http://mesos.apache.org/documentation/latest/configuration/master/)
> ./run_master.sh 마스터내부IP
0. ICGPATH=/opt/ichthysGenomics, 옵션파일경로=ICGPATH=/opt/ichthysGenomics/build/environments/mesos-master(MESOS_**형식)
1. 옵션파일에 필수 옵션 추가(masterip)
2. 서비스파일 systemd 등록
3. docker, zookeeper, mesos-master 실행 

## Agent 실행하기
* [master|agent 공통옵션 URL](http://mesos.apache.org/documentation/latest/configuration/master-and-agent/)
* [agent 옵션 URL](http://mesos.apache.org/documentation/latest/configuration/agent/)
> ./run_agent.sh 마스터내부IP
0. ICGPATH=/opt/ichthysGenomics, 옵션파일경로=ICGPATH=/opt/ichthysGenomics/build/environments/mesos-agent(MESOS_**형식)
1. 옵션파일에 필수 옵션 추가(masterip)
2. 서비스파일 systemd 등록
3. docker, mesos-agent 실행 
