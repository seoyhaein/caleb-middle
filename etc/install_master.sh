#!/bin/bash



##OS : centos 7.1
##Mesos version : 1.2.3
##Zookeeper version:3.4.14


set -e
set -x

## install default packages
## http://mesos.apache.org/documentation/latest/building/

yum install -y tar wget git
wget http://repos.fedorapeople.org/repos/dchen/apache-maven/epel-apache-maven.repo -O /etc/yum.repos.d/epel-apache-maven.repo
yum install -y epel-release
## yum install -y iptables-services

sudo bash -c "cat > /etc/yum.repos.d/wandisco-svn.repo <<EOF
[WANdiscoSVN]
name=WANdisco SVN Repo 1.9
enabled=1
baseurl=http://opensource.wandisco.com/centos/7/svn-1.9/RPMS/\$basearch/
gpgcheck=1
gpgkey=http://opensource.wandisco.com/RPM-GPG-KEY-WANdisco
EOF"

yum update -y systemd
yum groupinstall -y "Development Tools"
yum install -y apache-maven python-devel python-six python-virtualenv java-1.8.0-openjdk-devel zlib-devel libcurl-devel openssl-devel cyrus-sasl-devel cyrus-sasl-md5 apr-devel subversion-devel apr-util-devel

yum install -y yum-utils \
 device-mapper-persistent-data \
  lvm2

yum-config-manager \
  --add-repo \
  https://download.docker.com/linux/centos/docker-ce.repo

yum install -y docker-ce docker-ce-cli containerd.io



## create icg dir
CORE=${1}
ICGDIR=/opt/ichthysGenomics
ICGMESOSDIR=${ICGDIR}/mesos


mkdir -p ${ICGMESOSDIR}


##주키퍼 인스톨


cd ${ICGDIR}

wget http://apache.mirror.cdnetworks.com/zookeeper/zookeeper-3.4.14/zookeeper-3.4.14.tar.gz
tar -xzvf zookeeper-3.4.14.tar.gz
cd zookeeper-3.4.14 

## zookeeper configuration
cp conf/zoo_sample.cfg conf/zoo.cfg
sudo bash -c "cat > conf/zoo.cfg <<EOF
clientPort=2181
tickTime=2000
dataDir=/var/zookeeper
EOF"



##메소스 인스톨

cd ${ICGDIR}
git clone https://github.com/apache/mesos.git
cd ${ICGMESOSDIR}
git checkout 1.2.3
./bootstrap
mkdir build
cd ${ICGMESOSDIR}/build
../configure
cd ${ICGMESOSDIR}/build
make -j${CORE}
mkdir -p ${ICGDIR}/mesos/build/services
mkdir -p ${ICGDIR}/mesos/build/environments

## mesos master options
## http://mesos.apache.org/documentation/latest/configuration/master/
##MESOS_HOSTNAME=${HOSTNAME}
##MESOS_ZK=zk://${MASTERIP}:2181/mesos
sudo bash -c "cat > ${ICGMESOSDIR}/build/environments/mesos-master <<EOF
MESOS_QUORUM=1
MESOS_WORK_DIR=/var/lib/mesos
MESOS_AGENT_PING_TIMEOUT=15secs
MESOS_MAX_AGENT_PING_TIMEOUTS=1
EOF"

## mesos agent options
## http://mesos.apache.org/documentation/latest/configuration/agent/
## MESOS_MASTER=zk://${MASTERIP}:2181/mesos
# sudo bash -c "cat > ${ICGMESOSDIR}/build/environments/mesos-agent <<EOF
# MESOS_WORK_DIR=/var/lib/mesos
# MESOS_CONTAINERIZERS=docker,mesos
# MESOS_EXECUTOR_REGISTRATION_TIMEOUT=1mins
# EOF"

## mesos master environment
sudo bash -c "cat > ${ICGMESOSDIR}/build/environments/master-environment <<EOF
MESOS_WEBUI_DIR=${ICGMESOSDIR}/build/../src/webui
EOF"

## mesos agent environment
# sudo bash -c "cat > ${ICGMESOSDIR}/build/environments/agent-environment <<EOF
# MESOS_LAUNCHER_DIR=${ICGMESOSDIR}/build/src
# MESOS_NATIVE_JAVA_LIBRARY=${ICGMESOSDIR}/build/src/.libs/libmesos-1.10.0.so
# MESOS_NATIVE_LIBRARY=${ICGMESOSDIR}/build/src/.libs/libmesos-1.10.0.so
# EOF"


## mesos master service register in master
sudo bash -c "cat > ${ICGMESOSDIR}/build/services/mesos-master.service <<EOF
[Unit]
Description=Mesos Master: distributed systems kernel
After=network.target
[Service]
Restart=on-failure
StartLimitInterval=0
RestartSec=0s
LimitNOFILE=infinity
TasksMax=infinity
PermissionsStartOnly=True
SyslogIdentifier=mesos-master
EnvironmentFile=${ICGMESOSDIR}/build/environments/mesos-master
EnvironmentFile=${ICGMESOSDIR}/build/environments/master-environment
ExecStart=${ICGMESOSDIR}/build/src/mesos-master
EOF"


## zookeeper service register in master

sudo bash -c "cat > ${ICGMESOSDIR}/build/services/zookeeper.service <<EOF
[Unit]
Description=zookeeper
After=network.target
[Service]
Type=forking
User=root
Group=root
SyslogIdentifier=zookeeper
Restart=on-failure
RestartSec=0s
ExecStart=${ICGDIR}/zookeeper-3.4.14/bin/zkServer.sh start
EOF"
