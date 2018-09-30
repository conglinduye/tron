/**
 * @author [yanzheng]
 * @email [yan_zheng2018@163.com@mail.com]
 * @create date 2018-09-09 03:33:36
 * @modify date 2018-09-09 03:33:36
 * @desc [旧版浏览器其他服务入口]
 */

package main

import (
	"flag"
	"strings"

	"github.com/wlcy/tron/explorer/core/utils"

	"github.com/wlcy/tron/explorer/lib/config"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/ext/router"
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
	if strings.ToUpper(config.NetType) == "TESTNET" {
		utils.TestNet = true
	}

	//获取服务启动参数或其他参数
	var conf config.ConfigServer
	if 0 != conf.Populate(*configfile) {
		return
	}

	router.Start(conf.Address, conf.Objectpool)

}
