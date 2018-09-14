package config

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

type (
	//ConfigServer 配置文件中的http服务配置及mysql\redis连接配置
	ConfigServer struct {
		configFile string //配置文件地址，不对外直接访问

		Address    string      //http服务监听地址
		Objectpool int         //http服务对象池大小
		Mysql      ConfigMysql //mysql连接信息
		Redis      ConfigRedis //redis连接信息

		confTree *toml.Tree
	}

	//ConfigMysql MYSQL的连接信息
	ConfigMysql struct {
		Host     string //mysql地址
		Port     string //端口
		User     string //登录用户名
		Pass     string //登录用户密码
		Protocol string //"tcp"
		Schema   string //默认的schema
		Charset  string //"utf8"
	}
	//ConfigRedis Redis的连接信息
	ConfigRedis struct {
		Host     string //redis地址，如"127.0.0.1:6379"
		Pass     string //登录密码
		Index    int    //登录的数据库，默认为0
		Poolsize int    //连接池的大小
	}
)

/*/Populate 加载数据，如果成功，则返回0，否则返回非0值
func (conf *ConfigServer) Populate(path string) (errno int) {
	conf.configFile = path //保存配置文件路径
	//加载配置文件
	tree, err := toml.LoadFile(path)
	if err != nil {
		fmt.Println("Load config failed")
		return -1
	}

	conf.Address = tree.Get("server.address").(string)
	if 0 == len(conf.Address) {
		fmt.Println("address is not valid in config")
		return -1
	}

	conf.Objectpool = int(tree.Get("server.objectpool").(int64))
	if 0 == conf.Objectpool {
		fmt.Println("objectpool is not valid in config")
		return -1
	}

	conf.Mysql.Host = tree.Get("mysql.host").(string)
	if 0 == len(conf.Mysql.Host) {
		fmt.Println("mysql.host is not valid in config")
		return -1
	}

	conf.Mysql.Port = tree.Get("mysql.port").(string)
	if 0 == len(conf.Mysql.Port) {
		fmt.Println("mysql.port is not valid in config")
		return -1
	}

	conf.Mysql.User = tree.Get("mysql.user").(string)
	if 0 == len(conf.Mysql.User) {
		fmt.Println("mysql.user is not valid in config")
		return -1
	}

	conf.Mysql.Pass = tree.Get("mysql.pass").(string)
	if 0 == len(conf.Mysql.Pass) {
		fmt.Println("mysql.pass is not valid in config")
		return -1
	}

	conf.Mysql.Protocol = tree.Get("mysql.protocol").(string)
	if 0 == len(conf.Mysql.Protocol) {
		fmt.Println("mysql.protocol is not valid in config")
		return -1
	}

	conf.Mysql.Schema = tree.Get("mysql.schema").(string)
	if 0 == len(conf.Mysql.Schema) {
		fmt.Println("mysql.schema is not valid in config")
		return -1
	}

	conf.Mysql.Charset = tree.Get("mysql.charset").(string)
	if 0 == len(conf.Mysql.Charset) {
		fmt.Println("mysql.charset is not valid in config")
		return -1
	}

	conf.Redis.Host = tree.Get("redis.host").(string)
	if 0 == len(conf.Redis.Host) {
		fmt.Println("redis.host is not valid in config")
		return -1
	}

	conf.Redis.Poolsize = int(tree.Get("redis.poolsize").(int64))
	if 0 == conf.Redis.Poolsize {
		fmt.Println("redis.poolsize is not valid in config")
		return -1
	}

	conf.Redis.Pass = tree.Get("redis.pass").(string)
	conf.Redis.Index = int(tree.Get("redis.index").(int64))

	return 0
}
//*/

//Populate 加载数据，如果成功，则返回0，否则返回非0值
func (conf *ConfigServer) Populate(path string) (errno int) {
	return conf.Reload(path)
}

//Reload 加载数据，如果成功，则返回0，否则返回非0值
func (conf *ConfigServer) Reload(path string) (errno int) {
	conf.configFile = path //保存配置文件路径
	//加载配置文件
	tree, err := toml.LoadFile(path)
	if err != nil {
		fmt.Println("Load config failed")
		return -1
	}

	conf.confTree = tree

	// load http server info
	if conf.Address, err = conf.getValueAsString(tree, "server.address"); err != nil {
		return -1
	}
	if conf.Objectpool, err = conf.getValueAsInt(tree, "server.objectpool"); err != nil {
		return -1
	}
	/*
		// load mysql info
		if conf.Mysql.Host, err = conf.getValueAsString(tree, "mysql.host"); err != nil {
			return -1
		}
		if conf.Mysql.Port, err = conf.getValueAsString(tree, "mysql.port"); err != nil {
			return -1
		}
		if conf.Mysql.User, err = conf.getValueAsString(tree, "mysql.user"); err != nil {
			return -1
		}
		if conf.Mysql.Pass, err = conf.getValueAsString(tree, "mysql.pass"); err != nil {
			return -1
		}
		if conf.Mysql.Protocol, err = conf.getValueAsString(tree, "mysql.protocol"); err != nil {
			return -1
		}
		if conf.Mysql.Schema, err = conf.getValueAsString(tree, "mysql.schema"); err != nil {
			return -1
		}
		if conf.Mysql.Charset, err = conf.getValueAsString(tree, "mysql.charset"); err != nil {
			return -1
		}
		// load redis info
		if conf.Redis.Host, err = conf.getValueAsString(tree, "redis.host"); err != nil {
			return -1
		}
		if conf.Redis.Pass, err = conf.getValueAsString(tree, "redis.pass"); err != nil {
			return -1
		}
		if conf.Redis.Index, err = conf.getValueAsInt(tree, "redis.index"); err != nil {
			return -1
		}
		if conf.Redis.Poolsize, err = conf.getValueAsInt(tree, "redis.poolsize"); err != nil {
			return -1
		}
	*/
	return 0
}

// GetValueByKey 获取扩展信息
func (conf *ConfigServer) GetValueByKey(configItem string) (string, error) {
	return conf.getValueAsString(conf.confTree, configItem)
}

//getValueAsString 返回string类型的配置值
func (conf *ConfigServer) getValueAsString(tree *toml.Tree, configItem string) (string, error) {
	if false == conf.isExist(tree, configItem) {
		return "", conf.genNotExistError(configItem)
	}
	return tree.Get(configItem).(string), nil
}

//getValueAsInt64 返回int64类型的配置值
func (conf *ConfigServer) getValueAsInt64(tree *toml.Tree, configItem string) (int64, error) {
	if false == conf.isExist(tree, configItem) {
		return 0, conf.genNotExistError(configItem)
	}
	return tree.Get(configItem).(int64), nil
}

//getValueAsInt 返回int类型的配置值
func (conf *ConfigServer) getValueAsInt(tree *toml.Tree, configItem string) (int, error) {
	if false == conf.isExist(tree, configItem) {
		return 0, conf.genNotExistError(configItem)
	}
	return int(tree.Get(configItem).(int64)), nil
}

//isExist 检查某个配置项是否存在
func (conf *ConfigServer) isExist(tree *toml.Tree, configItem string) bool {
	var value = tree.Get(configItem)
	if nil == value {
		fmt.Printf("can not load config :[%v] from [%v]\n", configItem, conf.configFile)
		return false
	}
	return true
}

//genNotExistError 生成配置项不存在的错误
func (conf *ConfigServer) genNotExistError(configItem string) error {
	return fmt.Errorf("config item:[%v] is not exist", configItem)
}
