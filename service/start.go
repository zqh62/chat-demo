package service

import (
	"chat/pkg/e"
	"encoding/json"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
)

func (manager *ClientManager) Start() {
	for true {
		logging.Info("--------Monitoring the channel-------")
		select {
		case conn := <-Manager.Register:
			logging.Info("Receive a new connection ", conn.ID, ".")
			Manager.Clients[conn.SendID] = conn //将用户加载到用户管理器上
			replaymsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: e.GetMsg(e.WebsocketSuccess),
			}
			msg, _ := json.Marshal(replaymsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		}
	}
}
