package controller

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:		0,
	WriteBufferSize:	0,

	// 取消跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	ID 				string
	BlackConn		*websocket.Conn
	WhiteConn		*websocket.Conn

	CurrentColor	byte
	GoMap			[15][15]byte
}

var rooms map[string]*Room = make(map[string]*Room)

const BlackColor = 1
const WhiteColor = 2

type SocketMessage struct {
	Option 		string 	`json:"option"`		// "step" "new"    "lose" "win"
	I 			int		`json:"i"`
	J 			int		`json:"j"`
	Color 		int 	`json:"color"`      	// 落子颜色 1:黑色 2：白色

	RoomId		string 	`json:"room_id"`	// room id

	CurrentColor		int		// 当前轮到颜色
	ID 			string	`json:"id"`			// 对局ID号
	LastId		string	`json:"last_id"`	// new game 时上局对局ID
}

func randomId() string {
	return "3122"
}


func judgeWin(r *Room, m int, n int) bool{
	c := r.GoMap[m][n]

	count := 1
	// 横
	i := m
	for j := n - 1 ; j >= 0; j -- {
		if r.GoMap[i][j] == c{
			count ++
		} else{
			break
		}
	}
	for j := n + 1; j < 15; j ++ {
		if r.GoMap[i][j] == c{
			count ++
		} else{
			break
		}
	}
	if count >= 5 {
		return true
	}

	// 竖
	count = 1
	j := n
	for i := m - 1; i >= 0; i -- {
		if r.GoMap[i][j] == c{
			count ++
		} else{
			break
		}
	}
	for i := m + 1; i < 15; i ++ {
		if r.GoMap[i][j] == c{
			count ++
		} else{
			break
		}
	}
	if count >= 5 {
		return true
	}

	// 斜
	count = 1
	i = 1
	for {
		if m - i < 0 || n - i < 0{
			break
		}
		if r.GoMap[m-i][n-i] == c {
			count ++
		} else{
			break
		}
	}
	for {
		if m + i < 0 || n + i < 0{
			break
		}
		if r.GoMap[m+i][n+i] == c {
			count ++
		} else{
			break
		}
	}
	if count >= 5 {
		return true
	}

	// 反斜
	count = 1
	i = 1
	for {
		if m - i < 0 || n + i < 0{
			break
		}
		if r.GoMap[m-i][n+i] == c {
			count ++
		} else{
			break
		}
	}
	for {
		if m + i < 0 || n - i < 0{
			break
		}
		if r.GoMap[m+i][n-i] == c {
			count ++
		} else{
			break
		}
	}
	if count >= 5 {
		return true
	}
	return false
}

func WsHandler(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	var cRoom *Room
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

		fmt.Println(message)

		switch message.Option {
		case "step":
			if cRoom.WhiteConn != nil && cRoom.BlackConn != nil {
				cRoom.GoMap[message.I][message.J] = byte(message.Color)
				// 发送 step 事件
				cRoom.WhiteConn.WriteJSON(gin.H{
					"option": "step",
					"i": message.I,
					"j": message.J,
					"color": message.Color,
					"current_color": 3 - message.Color,
				})
				cRoom.BlackConn.WriteJSON(gin.H{
					"option": "step",
					"i": message.I,
					"j": message.J,
					"color": message.Color,
					"current_color": 3 - message.Color,
				})
				//if judgeWin(cRoom, message.I, message.J){
				//	fmt.Println("some win", message.Color)
				//	if message.Color == BlackColor{
				//		cRoom.BlackConn.WriteJSON(gin.H{
				//			"option": "win",
				//		})
				//		cRoom.WhiteConn.WriteJSON(gin.H{
				//			"option": "lose",
				//		})
				//	} else{
				//		cRoom.WhiteConn.WriteJSON(gin.H{
				//			"option": "win",
				//		})
				//		cRoom.BlackConn.WriteJSON(gin.H{
				//			"option": "lose",
				//		})
				//	}
				//}
			} else{
				fmt.Println(*cRoom)
			}
			break
		case "new":
			if cRoom.WhiteConn != nil && cRoom.BlackConn != nil {
				// 初始化 map
				for i := 0; i < len(cRoom.GoMap); i ++ {
					for j := 0; j < len(cRoom.GoMap[i]); j ++{
						cRoom.GoMap[i][j] = 0
					}
				}

				// 发送 new 事件
				cRoom.WhiteConn.WriteJSON(gin.H{
					"option": "new",
					"color": WhiteColor,
					"current_color": BlackColor,
				})
				cRoom.BlackConn.WriteJSON(gin.H{
					"option": "new",
					"color": BlackColor,
					"current_color": BlackColor,
				})
			}
			break
		// 加入房间事件
		case "join":
			if room, ok := rooms[message.RoomId]; ok {
				if room.BlackConn == nil {
					room.BlackConn = conn
					cRoom = room
					conn.WriteJSON(gin.H{
						"option": "join",
						"room_id": room.ID,
					})
				} else if room.WhiteConn == nil {
					room.WhiteConn = conn
					cRoom = room
					conn.WriteJSON(gin.H{
						"option": "join",
						"room_id": room.ID,
					})
				}
			}
			break
		}

	}
}

// 开房
func NewRoom(c *gin.Context){
	var roomId = randomId()
	rooms[roomId] = &Room{
		ID:				roomId,
		CurrentColor: 	BlackColor,
	}
	c.JSON(http.StatusOK, gin.H{
		"room_id": roomId,
	})
}