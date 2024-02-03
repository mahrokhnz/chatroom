package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type RoomResponse struct {
	Room
	UserCount int `json:"userCount"`
}

func initialMessage(username string, conn *websocket.Conn) {

	user, exist := Users[username]

	if exist {
		user.Conn = conn
		user.IsReady = true
		Users[username] = user

		if user.Room != nil {
			userRoom := Rooms[user.Room.Name]

			for i, currentUser := range userRoom.Users {
				if currentUser == user {
					userRoom.Users[i] = user
				}
			}

			Rooms[user.Room.Name] = userRoom
		}

		fmt.Printf("User Found \"%s\" :: Users: %d\n", username, len(Users))
	} else {
		user = NewUser()
		user.Create(username, conn)
	}

	rooms := make([]RoomResponse, 0)

	for _, room := range Rooms {
		rooms = append(rooms, RoomResponse{
			Room:      *room,
			UserCount: len(room.Users),
		})
	}

	user.Send(Response{
		Action: "Rooms",
		Data:   rooms,
	})

}

func createRoom(user *User, roomName interface{}) {
	// check if room name is not empty
	if roomName.(string) == "" {
		user.Send(Response{Action: "Fail", Data: "Room name is required!"})
		return
	}

	// check if room exists
	if _, exist := Rooms[roomName.(string)]; exist {
		user.Send(Response{
			Action: "Fail",
			Data:   "Room already created!",
		})
		return
	}

	room := NewRoom()
	room.Create(roomName.(string), user)
}

func joinRoom(user *User, roomName interface{}) {
	newRoom, ok := Rooms[roomName.(string)]

	// check if room exists
	if !ok {
		user.Send(Response{
			Action: "Fail",
			Data:   "Room not found!",
		})
		return
	}

	// check if joined already
	if user.Room != nil {
		if roomName == user.Room.Name {
			user.Send(Response{
				Action: "Fail",
				Data:   "You already Joined!",
			})
			return
		}

		user.Leave()
	}

	user.Join(newRoom)
}

func sendMessage(user *User, message interface{}) {
	//check if room exists
	if user.Room == nil {
		user.Send(Response{
			Action: "Fail",
			Data:   "Room not found! - You have not joined any room",
		})
		return
	}

	user.Room.SendMessage(message.(string), user)
}

func beforeUnload(user *User) {
	user.Leave()
	user.Close()
}
