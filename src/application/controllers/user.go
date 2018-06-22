package controllers

import (

	"application/core/websocket_gateway"
	"application/libraries/opcodes"
	"application/libraries/utils"
	"application/proto"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

func HandShake(client *websocket_gateway.Client, message []byte) {
	protocol := proto.RequestInit(&message)
	//utils.Logger().Println("connect22", &message)
	client.Idfa = protocol.Idfa
	fmt.Println("protokey",protocol.UserKey)
	//if sign(protocol) == false {
	//	fmt.Println("sing false")
	//	unregister(client.Id)
	//	return
	//}
	fmt.Println("info",protocol.UserKey)
	cid, ok := websocket_gateway.GetGateWay().UidToClientId.Load(protocol.UserKey)
	fmt.Println("cid",cid, ok)
	if ok && cid != client.Id {
		connId := cid.(string)
		oldClient, ok := websocket_gateway.GetGateWay().GetClient(connId)
		if ok && oldClient.Idfa != client.Idfa {
			data := map[string]string{"msg": "登陆提示：您的账号在另一台设备登录，已被迫下线。如果不是您本人所为，请尽快修改密码"}
			response := &proto.ClientProto{}
			response.Code = 2184
			response.Scene = 1
			response.Data = data
			response.Time = time.Now().Unix()

			message, _ := json.Marshal(response)
			utils.Logger().Println("client", oldClient)
			oldClient.Replay(message)
		}

		unregister(connId)
	}

	ids := &websocket_gateway.CliendIdToUid{}
	ids.ClientId = client.Id
	ids.UserId = protocol.UserKey
	fmt.Println(ids)
	websocket_gateway.GetGateWay().AddIdregister(ids)

	protos := proto.ClientProto{}
	protos.Code = opcodes.APP_HAND_SHAKE

	client.Info.Verson = protocol.Version
	client.Info.Platform = protocol.Platform

	data, _ := json.Marshal(protos)
	fmt.Println("data",data)
	fmt.Println("cccc",client.Info)
	client.Replay(data)

}

func sign(protocol *proto.ClientProtoRequest) bool {

	var token string
	var param map[string]interface{}

	idfa := protocol.Idfa
	uid := protocol.Uid
	version := protocol.Version
	platform := protocol.Platform
	if protocol.Parma != nil {
		param = protocol.Parma.(map[string]interface{})
		token = param["sign"].(string)
	}

	tokenStr := fmt.Sprintf("sign:%s:%s:%s:%d", version, idfa, platform, uid)
	md5Obj := md5.New()
	md5Obj.Write([]byte(tokenStr))
	serverToken := hex.EncodeToString(md5Obj.Sum(nil))
	if uid == 0 && idfa == "" {
		return false
	}

	if serverToken != token {
		return false
	}

	return true
}

func LoginOut(client *websocket_gateway.Client, message []byte) {
	unregister(client.Id)
}

func unregister(id string) {

	client, ok := websocket_gateway.GetGateWay().GetClient(id)
	if !ok {
		return
	}

	websocket_gateway.GetGateWay().AddUnregister(client)
}
