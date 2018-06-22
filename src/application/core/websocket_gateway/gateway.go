// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket_gateway

import (
	//"application/libraries/model"
	"application/libraries/utils"
	"sync"

	"github.com/robfig/cron"
)

type GatewayInfo struct {
	// Registered clients.
	clients sync.Map

	ClientIdToUid sync.Map

	UidToClientId sync.Map

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	idregister chan *CliendIdToUid
}

type CliendIdToUid struct {
	ClientId string
	UserId   string
}

var (
	Gateway *GatewayInfo
)

func GetGateWay() *GatewayInfo {

	if Gateway != nil {
		return Gateway
	}

	Gateway = &GatewayInfo{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		idregister: make(chan *CliendIdToUid),
		//clients:    make(map[*Client]bool),
	}

	return Gateway

}

func (h *GatewayInfo) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients.Store(client.Id, client)
			/*
				case client := <-h.unregister:
					if _, ok := h.clients.Load(client.Id); ok {

						h.clients.Delete(client.Id)
						if uid, ok := h.ClientIdToUid.Load(client.Id); ok {
							h.UidToClientId.Delete(uid.(string))
							h.ClientIdToUid.Delete(client.Id)
						}
						close(client.send)
						client.conn.Close()

					}

				case client := <-h.idregister:
					if id, ok := h.ClientIdToUid.Load(client.ClientId); ok {
						uid := id.(string)
						h.UidToClientId.Delete(uid)
					}

					h.ClientIdToUid.Store(client.ClientId, client.UserId)
					h.UidToClientId.Store(client.UserId, client.ClientId)
			*/
		}
	}
}

func Cron() {
	//获取管理员和
	c := cron.New()
	c.AddFunc("0 * * * * *", func() {
		//model.GetConfigMatchChannelAll() //同步配置
	})

	//心跳
	/*
		c.AddFunc("10 * * * * *", func() {
			GetGateWay().Heart()
		})
	*/
	c.Start()
}

func (h *GatewayInfo) Heart() {

	h.clients.Range(func(k, v interface{}) bool {

		if _, ok := h.ClientIdToUid.Load(v.(*Client).Id); !ok {
			h.unregister <- v.(*Client)
		}

		return false
	})
}

func (h *GatewayInfo) SendOne(id string, message []byte) {
	clientId, ok := h.UidToClientId.Load(id)
	if ok == false {
		return
	}

	if client, ok := h.clients.Load(clientId); ok {
		client.(*Client).send <- message
	}

}

func (h *GatewayInfo) SendAll(message []byte) {
	h.clients.Range(func(k, v interface{}) bool {

		v.(*Client).send <- message
		return true
	})
}

func (h *GatewayInfo) SendChange(cliendIds []string, message []byte) {

	for _, id := range cliendIds {

		client, err := h.clients.Load(id)
		if err == false {
			continue
		}

		client.(*Client).send <- message
	}
}

func (h *GatewayInfo) AddIdregister(client *CliendIdToUid) {

	if id, ok := h.ClientIdToUid.Load(client.ClientId); ok {
		uid := id.(string)
		h.UidToClientId.Delete(uid)
	}

	h.ClientIdToUid.Store(client.ClientId, client.UserId)
	h.UidToClientId.Store(client.UserId, client.ClientId)

	//h.idregister <- cliendIds
}

func (h *GatewayInfo) AddUnregister(client *Client) {

	if _, ok := h.clients.Load(client.Id); !ok {
		return
	}

	h.clients.Delete(client.Id)
	h.ClientIdToUid.Delete(client.Id)
	if uid, ok := h.ClientIdToUid.Load(client.Id); ok {
		h.UidToClientId.Delete(uid.(string))
	}

	utils.Logger().Println("AddUnregister close idfa:", client.Idfa)
	close(client.send)
	client.conn.Close()
	//h.unregister <- client

}

func (h *GatewayInfo) GetClient(id string) (*Client, bool) {

	if client, ok := h.clients.Load(id); ok {
		return client.(*Client), true
	}

	return nil, false
}
