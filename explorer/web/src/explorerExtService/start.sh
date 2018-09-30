#rm appLog*.log
logName=`date +%Y%m%d%H%M%S`
killall explorerExtService > /dev/null 2>&1
sleep 1
nohup ./explorerExtService -logLevel all > "appLog$logName.log" 2>&1 &
