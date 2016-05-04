#!/bin/sh

SSH=/usr/bin/ssh #ssh位置
KEY=/Users/lalo/kp-lalo-im #私钥位置
RUSER=root #主服务器帐号
RHOST=lalo.im #主服务器IP
RDIR=/root/doc
RFILE=/root/doc/godeploy #主服务器文件
LFILE=$GOPATH/bin/godeploy-amd64 #本地部署文件

echo "Start build amd64 version for linux ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -v -o $LFILE main.go

# 判断执行结果
if [[ $? -eq 0 ]]
then
	# 停止服务
	ssh -i $KEY $RUSER@$RHOST "supervisorctl stop godeploy"
	# 上传
	rsync -aSvHu --progress $LFILE -e "$SSH -i $KEY" $RUSER@$RHOST:$RFILE
	echo "Deploy ..."
	# 开始服务
	ssh -i $KEY $RUSER@$RHOST "supervisorctl start godeploy"
else
 	echo "Error!";
 fi
