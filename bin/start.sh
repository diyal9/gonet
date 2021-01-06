#nohup ./server "monitor" &
#sleep 1
nohup ./server "world" > /dev/null 2>logsys &
sleep 1
nohup ./server "account" > /dev/null 2>logsys &
sleep 1
nohup ./server "netgate" > /dev/null 2>logsys &
sleep 1
