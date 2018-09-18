/**
 * @author [yanzheng]
 * @email [yan_zheng2018@163.com@mail.com]
 * @create date 2018-09-09 03:33:36
 * @modify date 2018-09-09 03:33:36
 * @desc [区块链浏览器服务入口]
 */

package main

import (
	"flag"

	"github.com/wlcy/tron/explorer/lib/config"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/router"
	"github.com/wlcy/tron/explorer/web/task"
	"github.com/wlcy/tron/explorer/web/buffer"
)

// config file
var configfile = flag.String("cfgfile", "config.toml", "the config file path when running.")
var gLogFile = flag.String("log", "appLog", "set log base file, default is \"appLog\"")
var gDebug = flag.String("debug", "false", "debug flag default is \"false\"")
var gLogLevel = flag.String("logLevel", "info", "debug level default is Debug")

func main() {

	flag.Parse()
	log.ChangeLogLevel(log.Str2Level(*gLogLevel))

	//初始化db redis
	config.LoadConfig(*configfile)

	//获取服务启动参数或其他参数
	var conf config.ConfigServer
	if 0 != conf.Populate(*configfile) {
		return
	}

	//初始化buffer
	buffer.GetBlockBuffer()
	buffer.GetWitnessBuffer()
	buffer.GetMarketBuffer()
	buffer.GetVoteBuffer()
	buffer.GetAccountTokenBuffer()

	buffer.GetTokenBuffer()


	/* 数据库和redis初始化也可以用这种方式， but i don't like it
	redisCli = redis.NewClient(conf.Redis.Host, conf.Redis.Pass, conf.Redis.Index, conf.Redis.Poolsize)
	mysql.Initialize(conf.Mysql.Host, conf.Mysql.Port, conf.Mysql.Schema,
		conf.Mysql.User, conf.Mysql.Pass)
	*/

	go task.SyncCacheTodayReport()

	go task.SyncPersistYesterdayReport()

	go task.SyncAssetIssueParticipated()

	router.Start(conf.Address, conf.Objectpool)

}
