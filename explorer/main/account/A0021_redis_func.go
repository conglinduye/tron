package main

// redis key name
var (
	RedisSetAccountRefresh = "account:set:refresh" // 存放最近交易中出现的用户地址
)

// AddRefreshAddress ...
func AddRefreshAddress(addrs ...interface{}) (int64, error) {
	redisClient := getRedisClient()

	intRet := redisClient.SAdd(RedisSetAccountRefresh, addrs...)

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
