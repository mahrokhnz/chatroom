package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Data struct {
	Username string      `json:"username"`
	Action   string      `json:"action"`
	Data     interface{} `json:"data"`
	Excepted []*User     `json:"-"`
}

// config ws
var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 5,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		_, err := w.Write([]byte("WS Error: " + reason.Error()))
		if err != nil {
			fmt.Println("websocket Upgrader err: ", err)
		}
	},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConnectionHandler(w http.ResponseWriter, req *http.Request) {
	// generate ws
	conn, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(conn *websocket.Conn) {
		connCloseErr := conn.Close()
		if connCloseErr != nil {
			fmt.Println("ConnectionHandler Close Connection err: ", connCloseErr)
		}
		fmt.Printf("Close Connection %s \n", conn.RemoteAddr().String())
	}(conn)

	fmt.Printf("Open New Connection %s \n", conn.RemoteAddr().String())

	// check if you get data from client
	for {
		var data Data

		err = conn.ReadJSON(&data)

		if err != nil {
			fmt.Println("ConnectionHandler Read JSON err: ", err)
			return
		}

		// handle data
		if data.Username != "" {
			data.Handler(conn)
		} else {
			err = conn.WriteJSON(Data{
				Username: "",
				Action:   "Fail",
				Data:     "Empty Username!",
			})
			if err != nil {
				fmt.Println("ConnectionHandler Empty Username WriteJSON err: ", err)

				err = conn.Close()
				if err != nil {
					fmt.Println("ConnectionHandler Empty Username Close Connection err: ", err)
				}
			}
		}
	}
}

func (d Data) Handler(conn *websocket.Conn) {
	fmt.Println("data", d)

	user, ok := Users[d.Username]

	if d.Action != "" && d.Action != "InitialMessage" {
		if !ok {
			err := conn.WriteJSON(Data{
				Username: d.Username,
				Action:   "Fail",
				Data:     "User not found!",
			})
			if err != nil {
				fmt.Println("Handler WriteJSON err: ", err)

				err = conn.Close()
				if err != nil {
					fmt.Println("Handler Close err: ", err)
				}
			}
			return
		}

		if !user.IsReady {
			user.Send(Response{
				Action: "Fail",
				Data:   "User is not ready!",
			})
		}
	}

	switch d.Action {
	case "InitialMessage":
		initialMessage(d.Username, conn)
	case "CreateRoom":
		createRoom(user, d.Data)
	case "JoinRoom":
		joinRoom(user, d.Data)
	case "SendMessage":
		sendMessage(user, d.Data)
	case "BeforeUnload":
		beforeUnload(user)
	default:
		user.Send(Response{
			Action: "Fail",
			Data:   "Action not found!",
		})
	}
}
