package controller

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:		1024,
	WriteBufferSize:	1024,

	// 取消跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Game struct {
	Clients 	[2]*websocket.Conn			// 0 黑方  1 白方
	Table 		[15][15]byte				// 当前棋盘状态
	Current		byte						// 当前颜色
}

var games map[string]Game

const BlackColor = 1
const WhiteColor = 1

type SocketMessage struct {
	Option 		string 	`json:"option"`		// "step" "new"    "lose" "win"
	I 			int		`json:"i"`
	J 			int		`json:"j"`
	C 			int 	`json:"c"`      	// 落子颜色 1:黑色 2：白色

	Current		int		// 当前轮到颜色
	SetColor	int		// 发送给客户端 new 事件会带上 set color
	ID 			string	`json:"id"`			// 对局ID号
	LastId		string	`json:"last_id"`	// new game 时上局对局ID
}

func randomId() string {
	return "3122"
}

func newGame(message *SocketMessage, conn *websocket.Conn) {
	if "" == message.LastId {
		game := Game{
			Current:	BlackColor,
		}

		// 默认是黑色
		game.Clients[0] = conn


		// 生成ID
		id := randomId()

		// 存储 game
		games[id] = game
	} else{
		// 交换双方 conn 颜色


	}
}

func step (message *SocketMessage, conn *websocket.Conn){

}

func WsHandler(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var message SocketMessage

		if err := conn.ReadJSON(&message); err != nil {
			fmt.Println(err)
			return
		}

		switch message.Option {
		case "step":
			go step(&message, conn)
			break
		case "new":
			go newGame(&message, conn)
			break
		}
		//conn.WriteJSON(message)
	}
}

func NewRoom(c *gin.Context){
	var room_id = 3241
	c.JSON(http.StatusOK, gin.H{
		"room_id": room_id,
	})
}