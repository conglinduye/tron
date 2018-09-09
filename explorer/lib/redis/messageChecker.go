package redis

import "fmt"
import "sync"

//redisMsgKeyChecker 缓存的redisMessageKey, 用于检查该key是否已经被使用
var redisMsgKeyChecker = make(map[string]string, 0)

//checkerLocker the locker for redisMsgKeyChecker
var checkerLocker sync.Mutex

//CheckRedisKey 检查消息是否可用，如果可用，返回true,否则返回false
func CheckRedisKey(msgkey string) bool {
	checkerLocker.Lock()
	defer checkerLocker.Unlock()

	if len(msgkey) == 0 {
		fmt.Println(fmt.Errorf("message key [%v] is null, please check it", msgkey))
		return false
	}

	if _, ok := redisMsgKeyChecker[msgkey]; ok {
		fmt.Println(fmt.Errorf("message key [%v] have used, please check it", msgkey))
		return false
	}

	//缓存这个消息，用于其他消息的检查
	redisMsgKeyChecker[msgkey] = msgkey
	return true
}
