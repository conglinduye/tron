dsn="tron:tron@tcp(mine:3306)/tron"
startPos=0
step=50000
#round=0
run() {
    app=$1
    for round in `seq 5 10`;
    do
            echo $((startPos + step * round)) "   " $((startPos + step * (round + 1))) "    " fullnode_${round}_log_`date +'%Y%m%d%H%M%S'`
    done

            echo $((startPos + step * round)) "   " $((startPos + step * (round + 1))) "    " fullnode_${round}_log_`date +'%Y%m%d%H%M%S'`
}

run fullnode
