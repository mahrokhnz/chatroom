package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"slices"
	"sync"
	"time"
)

type User struct {
	Username string          `json:"username"`
	IsReady  bool            `json:"-"`
	Conn     *websocket.Conn `json:"-"`
	Room     *Room           `json:"-"`
	Mu       sync.Mutex      `json:"-"`
}

type Response struct {
	Data   interface{}
	Action string
}

var Users = make(map[string]*User)

func NewUser() *User {
	return &User{
		Username: "",
		IsReady:  false,
		Room:     nil,
		Conn:     nil,
		Mu:       sync.Mutex{},
	}
}

func (u *User) Create(username string, conn *websocket.Conn) {
	u.Username = username
	u.Conn = conn
	u.IsReady = true

	Users[username] = u

	fmt.Printf("New User \"%s\" :: Users: %d\n", username, len(Users))
}

func (u *User) Send(response Response) {
	u.Mu.Lock()

	defer u.Mu.Unlock()

	data := Data{
		Username: u.Username,
		Data:     response.Data,
		Action:   response.Action,
	}

	err := u.Conn.WriteJSON(data)
	if err != nil {
		fmt.Println(data.Username, " :: ", data.Action, " :: ", "Send to Client Err: ", err)
		u.Close()
	} else {
		fmt.Printf("Send Data '%s' to client: \"%s\" \n", data.Action, u.Username)
	}

	return
}

func (u *User) SendRaw(data Data) {
	u.Mu.Lock()

	defer u.Mu.Unlock()

	if u.Conn != nil {
		err := u.Conn.WriteJSON(data)
		if err != nil {
			fmt.Println(data.Username, " :: ", data.Action, " :: ", "Send Raw to Client Err: ", err)
			u.Close()
		} else {
			fmt.Printf("Send Raw Data '%s' to client: \"%s\" \n", data.Action, u.Username)
		}
	}

	return
}

func (u *User) Close() {
	if u.Conn != nil {
		closeErr := u.Conn.Close()
		if closeErr != nil {
			fmt.Println("User Close Conn Err: ", closeErr)
		} else {
			fmt.Printf("User Close Connection %s \n", u.Conn.RemoteAddr().String())
		}
	}

	u.Conn = nil
	u.IsReady = false
	Users[u.Username] = u

	return
}

func (u *User) Join(room *Room) {
	message := Message{
		Username:  u.Username,
		Text:      u.Username + " joined the group",
		CreatedAt: time.Now(),
		Type:      "Joined",
	}

	room.Users = append(room.Users, u)
	room.Usernames = append(room.Usernames, u.Username)
	room.Messages = append(room.Messages, message)

	Rooms[room.Name] = room
	u.Room = room
	Users[u.Username] = u

	fmt.Printf("User \"%s\" Join To Room \"%s\"\n", u.Username, room.Name)

	// send to all
	globalBroadcast <- Data{
		Username: u.Username,
		Action:   "JoinedRoom",
		Data:     RoomResponse{Room: *room, UserCount: len(room.Users)},
	}

	// send joined message to room clients
	room.BroadcastChan <- Data{
		Username: u.Username,
		Action:   "SentMessage",
		Data:     message,
		Excepted: []*User{u},
	}

	// send room info to room clients
	room.BroadcastChan <- Data{
		Username: u.Username,
		Action:   "UpdateRoomInfo",
		Data:     RoomResponse{Room: *room, UserCount: len(room.Users)},
	}
}

func (u *User) Leave() {
	room := u.Room

	if room == nil {
		return
	}

	message := Message{
		Username:  u.Username,
		Text:      u.Username + " left the group",
		CreatedAt: time.Now(),
		Type:      "Left",
	}

	// delete user
	for i, user := range room.Users {
		if u == user {
			room.Users = slices.Delete(room.Users, i, i+1)
			break
		}
	}

	// delete username
	for i, un := range room.Usernames {
		if un == u.Username {
			room.Usernames = slices.Delete(room.Usernames, i, i+1)
			break
		}
	}

	// add to room's messages
	room.Messages = append(room.Messages, message)

	Rooms[room.Name] = room
	u.Room = nil
	Users[u.Username] = u

	fmt.Printf("User \"%s\" Leave Room \"%s\"\n", u.Username, room.Name)

	// send to all
	globalBroadcast <- Data{
		Username: u.Username,
		Action:   "LeftRoom",
		Data:     RoomResponse{Room: *Rooms[room.Name], UserCount: len(Rooms[room.Name].Users)},
	}

	// send to room clients
	Rooms[room.Name].BroadcastChan <- Data{
		Username: u.Username,
		Action:   "SentMessage",
		Data:     message,
		Excepted: []*User{u},
	}

	// send room info to room clients
	Rooms[room.Name].BroadcastChan <- Data{
		Username: u.Username,
		Action:   "UpdateRoomInfo",
		Data:     RoomResponse{Room: *Rooms[room.Name], UserCount: len(Rooms[room.Name].Users)},
	}
}
