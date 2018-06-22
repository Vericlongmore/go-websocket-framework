package proto

import (
	"application/libraries/helpers"
	//"application/libraries/model/log"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
)

//客户端
type ClientProtoRequest struct {
	Code     uint16      `json:"code"`
	Uid      uint64      `json:"uid"`
	Idfa     string      `json:"idfa"`
	Platform string      `json:"platform"`
	Version  string      `json:"version"`
	Parma    interface{} `json:"param"`
	Source   string      `json:"source"`
	Ssid     uint64      `json:"ssid"`
	UserKey  string
	Content  []byte
}

//返回客户端协议
type RobotResultDagaProto struct {
	Code uint16      `json:"code"`
	Ssid string      `json:"ssid"`
	Data interface{} `json:"data"`
}

//返回客户端协议
type ClientProto struct {
	Code      uint16      `json:"code"`
	Ssid      uint64      `json:"ssid"`
	Broadcast int         `json:"broadcast"`
	SubCode   uint64      `json:"sub_code"`
	Msg       string      `json:"msg"`
	Seq       string      `json:"seq"`
	Question  string      `json:"question"`
	Power     int         `json:"power"`
	Time      int64       `json:"time"`
	Scene     uint8       `json:"scene"`
	Minute    string      `json:"minute"`
	Data      interface{} `json:"data"`
}

type ServerProto struct {
	Status    int8        `json:"status"`
	User      []string    `json:"user"`
	Idfa      []string    `json:"idfa"`
	Method    int         `json:"method"`
	MatchId   uint64      `json:"match_id"`
	Broadcast int         `json:"broadcast"`
	Channel   []uint64    `json:"channel"`
	Minute    string      `json:"minute"`
	Result    interface{} `json:"result"`
}

type GateWayContent struct {
	Clients []string    `json:"clients"`
	Content interface{} `json:"content"`
}

type ResponseAdminProto struct {
	Code uint64      `json:"code"`
	Seq  string      `json:"seq"`
	Data interface{} `json:"data"`
}

type ImgLink struct {
	Img         string `json:"img"`
	CompressImg string `json:"compress_img"`
}

type ImgData struct {
	Link  []ImgLink `json:"link"`
	Type  int       `json:"type"`
	Count int       `json:"count"`
}

func RequestInit(message *[]byte) *ClientProtoRequest {

	request := ClientProtoRequest{}
	err := json.Unmarshal(*message, &request)

	if err != nil {
		requestJson, _ := simplejson.NewJson(*message)
		ssid, _ := requestJson.Get("ssid").String()
		request.Ssid, _ = strconv.ParseUint(ssid, 10, 64)
	}

	if request.Code == 0 {
		return &request
	}

	if request.Uid > 0 {
		request.UserKey = fmt.Sprintf("%d", request.Uid)
	} else {
		request.UserKey = request.Idfa
	}

	if request.UserKey == "" {
		return &request
	}

	request.Content = *message
	//logs := formatLog(&request)
	//log.AddActionLog(*logs)
	requestJson1, _ := simplejson.NewJson(*message)
	fmt.Println(requestJson1)
	return &request
}

func formatLog(request *ClientProtoRequest) *[]byte {

	ssid := request.Ssid
	action := request.Code
	params := request.Parma
	userId := request.Uid
	idfa := request.Idfa
	platform := request.Platform
	source := request.Source
	log := make(map[string]interface{})
	log["action"] = fmt.Sprintf("websokect_code_%d", action)
	log["params"] = params
	log["ip"] = ""
	log["login_id"] = userId
	log["time"] = time.Now().Local().Format("2006-01-02 15:04:05")
	log["idfa"] = idfa
	log["server_ip"] = helpers.GetServerIp()
	log["user_agent"] = ""
	log["match_id"] = ssid
	log["platform"] = platform
	log["source"] = source
	log["message"] = request.Content

	logStr, _ := json.Marshal(log)

	return &logStr
}

func (response *ClientProto) SetPower(innerProto *ServerProto) {

	switch innerProto.Method {
	case 1000:
		response.Power = 1
	case 1001:
		response.Power = 2
	case 1002:
		response.Power = 3
	case 1003:
		response.Power = 4
	case 1004:
		response.Power = 5
	case 1005:
		response.Power = 6
	case 1006:
		response.Power = 7
	case 1007:
		response.Power = 8
	default:
		response.Power = 10
	}
}
