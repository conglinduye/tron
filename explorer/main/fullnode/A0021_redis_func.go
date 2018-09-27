package main

import (
	"fmt"
	"sync"
	"time"
)

// redis key name
var (
	RedisSetAccountRefresh = "account:set:refresh" // 存放最近交易中出现的用户地址
)

var _bufCap = 10000
var _refresgAddrBuffer = make([]interface{}, 0, _bufCap)
var _refBufLock = sync.Mutex{}
var _latestPushTS = time.Now()

func setTestNetRedisKey() {
	RedisSetAccountRefresh = fmt.Sprintf("test_net:%s", RedisSetAccountRefresh)
}

// AddRefreshAddress ...
func AddRefreshAddress(addrs ...interface{}) (newLen int, err error) {
	if 0 == len(addrs) {
		return 0, nil
	}
	newLen = len(addrs)
	_refBufLock.Lock()
	bufLen := len(_refresgAddrBuffer)

	if needQuit() {
		redisSADD(addrs)
	} else {
		if bufLen+newLen > _bufCap {
			// fmt.Printf("%v\n%v-->%v\n\n", _refresgAddrBuffer, bufLen, cap(_refresgAddrBuffer))
			curAddrs, err := redisSADD(_refresgAddrBuffer)
			if nil != err {
				_ = curAddrs
				fmt.Printf("push account address to redis failed:%v\n", err)
			}

			// fmt.Printf("push %v address to redis for later synchronize, address count:%v, err:%v\n", bufLen, curAddrs, err)
			_refresgAddrBuffer = _refresgAddrBuffer[:0]
			_latestPushTS = time.Now()
		}
		_refresgAddrBuffer = append(_refresgAddrBuffer, addrs...)
	}
	_refBufLock.Unlock()

	return
}

func startRedisAccountRefreshPush() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_refBufLock.Lock()
			if time.Since(_latestPushTS) > time.Duration(30)*time.Second {
				// fmt.Printf("%v\n%v-->%v\n\n", _refresgAddrBuffer, bufLen, cap(_refresgAddrBuffer))
				// bufLen := len(_refresgAddrBuffer)
				curAddrs, err := redisSADD(_refresgAddrBuffer)
				// fmt.Printf("push %v address to redis for later synchronize, address count:%v, err:%v\n", bufLen, curAddrs, err)
				if nil != err {
					_ = curAddrs
					fmt.Printf("push account address to redis failed:%v\n", err)
				} else {
					_refresgAddrBuffer = _refresgAddrBuffer[:0]
				}
				_latestPushTS = time.Now()
			}
			_refBufLock.Unlock()

			time.Sleep(30 * time.Second)
			if needQuit() {
				break
			}
		}
		cleanAccountBuffer()
		fmt.Printf("Redis Account Refresh Push Daemon QUIT\n")
	}()
}

func cleanAccountBuffer() {
	_refBufLock.Lock()
	bufLen := len(_refresgAddrBuffer)
	curAddrs, err := redisSADD(_refresgAddrBuffer)

	_ = bufLen
	_ = curAddrs
	_ = err
	// fmt.Printf("push %v address to redis for later synchronize, address count:%v, err:%v\n", bufLen, curAddrs, err)

	_refresgAddrBuffer = _refresgAddrBuffer[:0]
	_latestPushTS = time.Now()

	_refBufLock.Unlock()
}

func redisSADD(val []interface{}) (int64, error) {
	redisClient := getRedisClient()

	intRet := redisClient.SAdd(RedisSetAccountRefresh, val...)

	if nil == intRet {
		return 0, ErrorRedisNilResult
	}
	return intRet.Val(), intRet.Err()
}

// ClearRefreshAddress ...
func ClearRefreshAddress() ([]string, error) {
	cnt := _redisCli.SCard(RedisSetAccountRefresh)
	if nil != cnt {
		if nil != cnt.Err() {
			return nil, cnt.Err()
		}

		strSliceRet := _redisCli.SPopN(RedisSetAccountRefresh, cnt.Val())
		if nil == strSliceRet {
			return nil, ErrorRedisNilResult
		}
		if nil != strSliceRet.Err() {
			return nil, strSliceRet.Err()
		}

		return strSliceRet.Val(), strSliceRet.Err()

	}
	return nil, ErrorRedisNilResult
}
