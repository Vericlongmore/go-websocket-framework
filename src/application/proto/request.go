package proto

import (
	//"application/libraries/helpers"
	//"application/libraries/opcodes"
	//"encoding/json"
)

type RequestRobotProto struct {
	User string      `json:"user"`
	Idfa string      `json:"idfa"`
	Data interface{} `json:"data"`
}

