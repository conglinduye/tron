package util

import (
	"net/http"
	"reflect"
	"sync"

	"github.com/wlcy/tron/explorer/lib/log"
)

//ServicePool 服务容器, 提供
type ServicePool struct {
	sList []http.Server // 池

	lock sync.RWMutex // 锁

	sCap int // 服务池大小

	addr string // 服务地址, 默认 ":8080"

}

//NewServicePool 创建新的服务对象池, 默认大小为0, 默认路由为 "/", 默认服务地址为 ":8080"
func NewServicePool(s http.Server, sCap int, addr string) *ServicePool {
	if sCap == 0 {
		sCap = 10
	}

	if len(addr) == 0 {
		addr = ":8080"
	}

	st := reflect.SliceOf(reflect.TypeOf(s))
	// []*base.OrderService 不能转换成 []base.Server, 需要手动一个一个转
	ssList := reflect.MakeSlice(st, 0, sCap)
	sList := make([]http.Server, 0, sCap)
	// as len(sLis) == 0, following for will not execute
	for i := range sList {
		tmp := ssList.Index(i).Interface().(http.Server)
		log.Debugf("tmp [%v] == [%v] [%T]", i, tmp, tmp)
	}
	// panic: interface is []*base.OrderService, not []base.Server
	// panic: interface conversion: interface is []*base.OrderService, not []base.Server
	//sList := reflect.MakeSlice(st, 0, sCap).Interface().([]Server)
	// construct service
	objType := reflect.TypeOf(s).Elem() // s Server 应为指针类型, 我们使用 Elem() 得到实体类型
	for i := 0; i < sCap-1; i++ {       // 实体类型使用 reflect.New 得到指针类型, 刚好满足 Server
		sList = append(sList, reflect.New(objType).Interface().(http.Server))
	}
	sList = append(sList, s)

	ret := &ServicePool{
		sList: sList,
		sCap:  sCap,
		addr:  addr,
	}
	return ret
}

//Start 启动服务, 阻塞方法, 启动失败返回错误
func (sp *ServicePool) Start() error {
	// 启动服务
	log.Debugf("Start service, address:[%v], Server Type:[%T], PoolSize:[%v]",
		sp.addr, sp.sList, cap(sp.sList))

	return http.ListenAndServe(sp.addr, nil)
}
