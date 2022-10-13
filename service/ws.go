package service

import (
	"chat/cache"
	"chat/pkg/e"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const month = 60 * 60 * 24 * 30 //默认一个月30天

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Manager 全局管理器
var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid //1->2
}

func Handler(context *gin.Context) {
	uid := context.Query("uid")
	toUid := context.Query("toUid")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(context.Writer, context.Request, nil) //升级ws协议
	if err != nil {
		http.NotFound(context.Writer, context.Request)
		return
	}

	//创建用户实例
	client := &Client{
		ID:     CreateID(uid, toUid),
		SendID: CreateID(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}

	//用户注册到管理器上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (client *Client) Read() {
	//让Read操作完成时，回收资源，Client回到Manger的Unregister中
	defer func() {
		Manager.Unregister <- client
		_ = client.Socket.Close()
	}()

	for true {
		client.Socket.PongHandler()
		sendMsg := new(SendMsg)
		err := client.Socket.ReadJSON(&sendMsg)
		if err != nil {
			logging.Info("Wrong data format.")
			Manager.Unregister <- client
			_ = client.Socket.Close()
			break
		}
		if sendMsg.Type == 1 { //发送消息
			r1, _ := cache.RedisClient.Get(client.ID).Result()     //1->2
			r2, _ := cache.RedisClient.Get(client.SendID).Result() //2->1
			if r1 > "3" && r2 == "" {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: e.GetMsg(e.WebsocketLimit),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(client.ID)
				_, _ = cache.RedisClient.Expire(client.ID, time.Hour*24*30*3).Result()
				//防止过快"分手"，设置三个月过期
			}
			Manager.Broadcast <- &Broadcast{
				Client:  client,
				Message: []byte(sendMsg.Content),
			}
		}
	}
}

func (client *Client) Write() {
	defer func() {
		_ = client.Socket.Close()

	}()

	for true {
		select {
		case _, ok := <-client.Send:
			if !ok {
				_ = client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replymsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: e.GetMsg(e.WebsocketSuccessMessage),
			}
			msg, _ := json.Marshal(replymsg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
