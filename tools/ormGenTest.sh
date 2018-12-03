#! /bin/bash
dbName=winmode
tableExp=questions
$GOPATH/bin/xorm reverse \
mysql devuser:devuser@tcp\(192.168.1.174:3306\)/${dbName}?charset=utf8mb4 \
$GOPATH/src/github.com/qinhao/botKit/xorm/tpl \
$GOPATH/src/github.com/qinhao/botKit/xorm/models ${tableExp}