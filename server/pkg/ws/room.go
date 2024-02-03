package ws

import (
	"fmt"
	"slices"
	"time"
)

type Room struct {
	Name          string    `json:"name"`
	Creator       *User     `json:"creator"`
	Users         []*User   `json:"users"`
	Usernames     []string  `json:"usernames"`
	Messages      []Message `json:"messages"`
	CreatedAt     time.Time `json:"createdAt"`
	BroadcastChan chan Data `json:"-"`
}

var Rooms = make(map[string]*Room)

func NewRoom() *Room {
	return &Room{
		Name:          "",
		Creator:       nil,
		Users:         make([]*User, 0),
		Usernames:     make([]string, 0),
		Messages:      make([]Message, 0),
		CreatedAt:     time.Now(),
		BroadcastChan: make(chan Data),
	}
}

func (r *Room) Listener() {
	for {
		data := <-r.BroadcastChan
		fmt.Printf("Send Data '%s' to \"%s's\" clients\n", data.Action, r.Name)

		for _, client := range r.Users {
			if slices.Contains(data.Excepted, client) {
				fmt.Printf("Excepted \"%s\"\n", client.Username)

				continue
			}
			client.SendRaw(data)
		}
	}
}

func (r *Room) SendMessage(message string, user *User) {
	// create new message
	messageData := Message{
		Username:  user.Username,
		Text:      message,
		Type:      "Message",
		CreatedAt: time.Now(),
	}

	// add to room's messages
	r.Messages = append(r.Messages, messageData)

	Rooms[r.Name] = r

	// send to current room clients
	r.BroadcastChan <- Data{
		Username: user.Username,
		Action:   "SentMessage",
		Data:     messageData,
	}
}

func (r *Room) Create(name string, user *User) {
	// create room
	r.Name = name
	r.Creator = user

	// start gorouting
	go r.Listener()

	Rooms[r.Name] = r

	fmt.Printf("New Room \"%s\"\n", name)

	// send to all
	globalBroadcast <- Data{
		Username: user.Username,
		Action:   "CreatedRoom",
		Data:     []RoomResponse{{Room: *r, UserCount: len(r.Users)}},
	}
}
