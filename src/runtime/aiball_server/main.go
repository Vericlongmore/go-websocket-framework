package main

import (
	"application/core/websocket_gateway"
	"application/libraries/opcodes"
	"application/libraries/toml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"application/controllers"
)

var addr *string

func init() {
	_, err := toml.LoadTomlConfig("./etc/config.toml")
	if err != nil {
		panic(err)
	}
	OuterAddr := fmt.Sprintf("%s:%d", toml.GlobalTomlConfig.Pillx.GatewayOuterHost, toml.GlobalTomlConfig.Pillx.GatewayOuterPort)
	addr = flag.String("addr", OuterAddr, "http service address")
	fmt.Println(OuterAddr)
}

/**

js websocket post data：

{"uid":579501,"platform":"iOS","code":10001,"ssid":"","version":"1.30",
"param":{"ack":"c121c230-eba5-11e7-b63c-00163e0ffda7","idx":"3","type":"人员"},
"idfa":"2f3c365736ae9efce5ab2974806a8c3999daded3"}

 */

func main() {

	flag.Parse()
	hub := websocket_gateway.GetGateWay()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket_gateway.ServeWs(hub, w, r)
	})

	websocket_gateway.GetRuteInstance().HandleFunc(opcodes.APP_LOGIN_OUT, controllers.LoginOut)   //退出
	websocket_gateway.GetRuteInstance().HandleFunc(opcodes.APP_HAND_SHAKE, controllers.HandShake) //握手
	go hub.Run()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
