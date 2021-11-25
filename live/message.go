package live

import (
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	clients = make(map[uint]map[*websocket.Conn]bool)
	mux     sync.Mutex
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 20 * time.Second,
	// 取消 ws 跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler 处理ws请求
func WsHandler(c *gin.Context) {
	var conn *websocket.Conn
	var err error

	conn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	var json struct{
		Jwt string `json:"jwt"`
	}

	err = conn.ReadJSON(&json)

	if err != nil {
		log.Println("read json err:", err)
		_ = conn.Close()
		return
	}

	// 验证
	var userId uint
	userId, err = model.ValidateJWT(json.Jwt)

	if err != nil {
		log.Println("validate jwt err:", err)
		_ = conn.Close()
		return
	}

	// 将上下文保存到字典
	addClient(userId, conn)
}

func SendMessage(fromId, toId uint, title, content string) {
	model.SendMessage(fromId, toId, title, content)
	setMessage(toId, gin.H{
		"data": gin.H{
			"title":   title,
			"content": content,
		},
	})
}

func BatchSendMessage(fromId uint, toId []uint, title, content string) {
	model.BatchSendMessage(fromId, toId, title, content)
	for i := range toId {
		setMessage(toId[i], gin.H{
			"data": gin.H{
				"title":   title,
				"content": content,
			},
		})
	}
}

func addClient(id uint, conn *websocket.Conn) {
	mux.Lock()
	if clients[id] == nil {
		clients[id] = make(map[*websocket.Conn]bool)
	}
	clients[id][conn] = true
	mux.Unlock()
}

func getClients(id uint) (conns []*websocket.Conn) {
	mux.Lock()
	_conns, ok := clients[id]
	if ok {
		for k := range _conns {
			conns = append(conns, k)
		}
	}
	mux.Unlock()
	return
}

func deleteClient(id uint, conn *websocket.Conn) {
	mux.Lock()
	_ = conn.Close()
	delete(clients[id], conn)
	mux.Unlock()
}

func SetMessage(userId uint, content interface{}) {
	setMessage(userId, content)
}

func setMessage(userId uint, content interface{}) {
	conns := getClients(userId)
	for i := range conns {
		i := i
		err := conns[i].WriteJSON(content)
		if err != nil {
			log.Println("write json err:", err)
			deleteClient(userId, conns[i])
		}
	}
}
