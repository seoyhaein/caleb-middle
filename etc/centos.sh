#!/bin/bash

#### SET VPN CONFIG
echo -n "Enter YOUR HAIIP ID :  "
read HAIIPID

echo -n "Enter YOUR HAIIP PASSWD :  "
read PASSWD

echo -n "Enter SET SERVER :  "
read SERVERID


#### chap-secrets SET
echo -e  "$HAIIPID	*	$PASSWD	*" > /etc/ppp/chap-secrets


touch /etc/ppp/peers/haiip
USERID=$(cat /etc/ppp/chap-secrets | awk {'print $1'})

echo "pty "\"pptp $SERVERID --nolaunchpppd"\"" > /etc/ppp/peers/haiip
echo "name "\"$USERID"\"" >> /etc/ppp/peers/haiip
echo "remotename PPTP" >> /etc/ppp/peers/haiip
echo "file /etc/ppp/options.pptp" >> /etc/ppp/peers/haiip
echo "ipparam HAIIP" >> /etc/ppp/peers/haiip


touch /etc/ppp/ip-up.local
echo "#!/bin/sh" >> /etc/ppp/ip-up.local
echo "route add -net 0.0.0.0/0 dev ppp0" >> /etc/ppp/ip-up.local

chmod 755 /etc/ppp/ip-up.local
