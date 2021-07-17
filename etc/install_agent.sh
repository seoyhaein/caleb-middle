#!/bin/bash



##OS : centos 7.1
##Mesos version : 1.2.3
##OpenSSL version: 1.1.1
##python version : 3.6


set -e
set -x


## install default packages
## http://mesos.apache.org/documentation/latest/building/

yum install -y tar wget git
wget http://repos.fedorapeople.org/repos/dchen/apache-maven/epel-apache-maven.repo -O /etc/yum.repos.d/epel-apache-maven.repo
yum install -y epel-release

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
yum install -y gcc zlib-devel apache-maven python-devel python-six python-virtualenv java-1.8.0-openjdk-devel zlib-devel libcurl-devel openssl-devel cyrus-sasl-devel cyrus-sasl-md5 apr-devel subversion-devel apr-util-devel

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

## 후에 awscli를 위한 파이썬, openssl, awscli(마운트시켜서 실행파일로 실행) 설치
SSL=${ICGDIR}/openssl
PYTHON=${ICGDIR}/python

cd ${ICGDIR}

wget https://github.com/openssl/openssl/archive/OpenSSL_1_1_1c.tar.gz
tar -xzvf OpenSSL_1_1_1c.tar.gz 
cd openssl-OpenSSL_1_1_1c
./config --prefix=$SSL
make -j${CORE} 
make install
cd ${ICGDIR}
rm -rf OpenSSL_1_1_1c.tar.gz



export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:${SSL}/lib

git clone https://github.com/python/cpython 

cd cpython && git checkout 3.6

echo -e "SSL=$SSL\n_ssl _ssl.c \\ \n \t-DUSE_SSL -I\$(SSL)/include -I\$(SSL)/include/openssl \\  \n \t-L\$(SSL)/lib -lssl -lcrypto" >> Modules/Setup.dist

./configure SSL="${SSL}" CPPFLAGS="-I${SSL}/include" LDFLAGS="-L${SSL}/lib" --prefix=${PYTHON} 


make -j${CORE}
make install 


${PYTHON}/bin/pip3 install awscli



##메소스 인스톨

cd ${ICGDIR}
git clone https://github.com/apache/mesos.git
cd ${ICGMESOSDIR}
git checkout 1.11.0
./bootstrap
mkdir -p ${ICGMESOSDIR}/build
cd ${ICGMESOSDIR}/build
../configure
cd ${ICGMESOSDIR}/build
make -j${CORE}
mkdir -p ${ICGDIR}/mesos/build/services
mkdir -p ${ICGDIR}/mesos/build/environments

## mesos agent options
## http://mesos.apache.org/documentation/latest/configuration/agent/
## MESOS_MASTER=zk://${MASTERIP}:2181/mesos
sudo bash -c "cat > ${ICGMESOSDIR}/build/environments/mesos-agent <<EOF
MESOS_WORK_DIR=/var/lib/mesos
MESOS_CONTAINERIZERS=docker,mesos
MESOS_EXECUTOR_REGISTRATION_TIMEOUT=1mins
EOF"


## mesos agent environment
sudo bash -c "cat > ${ICGMESOSDIR}/build/environments/agent-environment <<EOF
MESOS_LAUNCHER_DIR=${ICGMESOSDIR}/build/src
MESOS_NATIVE_JAVA_LIBRARY=${ICGMESOSDIR}/build/src/.libs/libmesos-1.10.0.so
MESOS_NATIVE_LIBRARY=${ICGMESOSDIR}/build/src/.libs/libmesos-1.10.0.so
EOF"

## mesos agent service register in agent
sudo bash -c "cat > ${ICGMESOSDIR}/build/services/mesos-agent.service <<EOF
[Unit]
Description=Mesos Agent: distributed systems kernel
[Service]
Restart=always
StartLimitInterval=0
RestartSec=15
LimitNOFILE=infinity
TasksMax=infinity
PermissionsStartOnly=True
SyslogIdentifier=mesos-agent
EnvironmentFile=${ICGMESOSDIR}/build/environments/mesos-agent
EnvironmentFile=${ICGMESOSDIR}/build/environments/agent-environment
ExecStart=${ICGMESOSDIR}/build/src/mesos-agent
EOF"


#systemctl daemon-reload
