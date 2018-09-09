package redis

import (
	"testing"
)

// 功能测试
// 功能测试函数以 Test 开头, Xxx结束(首字母必须大写,其他无要求)
// 入参为 *testing.T, 无返回值
func TestNewClient(t *testing.T) {
	client := NewClient("localhost:6379", "", 0, 0)
	t.Logf("create new client:[%#v]\n", client)
	t.Logf("client.Ping() ret:%v\n", client.Ping())
}

// 性能测试
// 性能测试函数以 Benchmark 开头, Xxx结束(首字母必须大写,其他无要求)
// 入参为 *testing.B, 无返回值
func BenchmarkNewClient(b *testing.B) {
}
