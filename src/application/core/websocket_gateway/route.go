package websocket_gateway

import (
	"application/libraries/utils"
	"sync"

	"github.com/bitly/go-simplejson"
	"fmt"
)

var (
	route *ServeRouter
)

type ServeRouter struct {
	opcode_list sync.Map
}

type Handler interface {
	Serve(*Client, []byte)
}

// 这里将HandlerFunc定义为一个函数类型，因此以后当调用a = HandlerFunc(f)之后, 调用a的serve实际上就是调用f的对应方法, 拥有相同参数和相同返回值的函数属于同一种类型。
type HandlerFunc func(*Client, []byte)

func GetRuteInstance() *ServeRouter {

	if route == nil {
		route = &ServeRouter{}
	}

	return route
}

// Serve calls f(w, r).
func (f HandlerFunc) Serve(client *Client, message []byte) {
	f(client, message)
}

//将router对应的opcode,方法存储
func (router *ServeRouter) Handle(name uint16, handler Handler) {
	router.opcode_list.Store(name, handler)
}

// HandleFunc registers the handler function for the given opcode.
func (rounter *ServeRouter) HandleFunc(name uint16, handler func(*Client, []byte)) {
	rounter.Handle(name, HandlerFunc(handler))
}

// 取出opcode对应的操作方法,然后回调
func (router *ServeRouter) Serve(client *Client, message []byte) {

	var handler Handler

	fmt.Println("serve")

	j, _ := simplejson.NewJson(message)
	if j == nil {
		utils.Logger().Println("error :", message)
		return
	}
	code, _ := j.Get("code").Int()
	fmt.Println("code",code)

	if code == 0 {
		fmt.Println("ii",string(message))
		return
	}

	rute, ok := router.opcode_list.Load(uint16(code))
	fmt.Println(rute,ok)
	if ok {
		handler = rute.(Handler)
		fmt.Println(handler,client)
		handler.Serve(client, message)
	}
}
