# caleb-middle
## 해야 할 일 생각나는데로 적음.(일단 크게 이렇게 구성) - 5/10

- grpc server
- mesos 에서 resource offer 시 resource offer 저장 struct 구성
- jobscheduler 에서 해당 offer 를  accept 하거나 decline/suppress 해야함.
- DAG 사용해서 부모 자식 job 에 대한 순서를 구현할지 생각해 봐야 함.
- zk 의 경우... 고려해야함. 이걸 할지 말지를..
- 테스트 코드.....
- go-grpc-prometheus 참고
- TaskInfo 에서 shell script 작성에 대한 부분
- ContainerInfo?
- ec2 scale in/out, 기타 클러스터 구성(시간 오래되서 잊어버림. 젠장)
- mesos 설치, 관련 스크립트 제작 해야함. 이건 짜투리 시간을 활용하자.(11)
- aws ec2 계정 2개 살펴보자 -> 따로 문서로 정리 하자 -> ichthysgenomics@icg.ai
- upd 계정 2개 살리자. -> 문의후 aws 무료 계정에 2개 서버 살려두자.