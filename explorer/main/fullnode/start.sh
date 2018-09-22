dsn="tron:tron@tcp(mine:3306)/tron"
startPos=0
step=50000
#round=0
run() {
    app=$1
	for round in `seq 5 10`;
	do
    	    nohup ./${app} -start_block $((startPos + step * round)) -end_block $((startPos + step * (round + 1))) -worker 64 -workload 500 -dsn "${dsn}"  > fullnode_${round}_log_`date +'%Y%m%d%H%M%S'`.log 2>&1 &
	done
    #nohup ./${app} -start_block 50000 -end_block 100000 -worker 64 -workload 500 -dsn "${dsn}"  > fullnode_2_log_`date +'%Y%m%d%H%M%S'`.log 2>&1 &
    #nohup ./${app} -start_block 100000 -end_block 150000 -worker 64 -workload 500 -dsn "${dsn}"  > fullnode_3_log_`date +'%Y%m%d%H%M%S'`.log 2>&1 &
    #nohup ./${app} -start_block 150000 -end_block 200000 -worker 64 -workload 500 -dsn "${dsn}"  > fullnode_4_log_`date +'%Y%m%d%H%M%S'`.log 2>&1 &
    #nohup ./${app} -start_block 200000 -end_block 250000 -worker 64 -workload 500 -dsn "${dsn}"  > fullnode_5_log_`date +'%Y%m%d%H%M%S'`.log 2>&1 &
    #nohup ./${app} -start_block 2500000 -worker 64 -workload 5000 -dsn "${dsn}"  > fullnode_now_log_`date +'%Y%m%d%H%M%S'`.log 2>&1 &
}

run fullnode
