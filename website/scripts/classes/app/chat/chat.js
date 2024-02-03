import Room from "./room/room.js";

class Chat {
    socket;
    username;

    roomsElement;
    messagesElement;
    roomElement;

    room;

    actionHandlers = {
        Fail: this.#failHandler,
        Rooms: this.#roomsHandler,
        JoinRoom: this.#joinRoomHandler,
        JoinedRoom: this.#joinedRoomHandler,
        SendMessage: this.#sendMessageHandler.bind(this),
        SentMessage: this.#sentMessageHandler,
        LeftRoom: this.#leftRoomHandler,
        CreateRoom: this.#createRoomHandler.bind(this),
        CreatedRoom: this.#createdRoomHandler.bind(this),
        UpdateRoomInfo: this.#updateRoomInfoHandler.bind(this),
    };

    constructor(socket, username) {
        this.socket = socket;
        this.username = username;

        this.roomsElement = document.querySelector('.roomsWrapper');
        this.messagesElement = document.querySelector('.messagesWrapper');
        this.roomElement = document.querySelector('.roomInfo')

        this.socket.addHandler('open', this.#openHandler.bind(this));
        this.#beforeunloadHandler.call(this);
        this.socket.addHandler('message', this.#messageHandler.bind(this));
    }

    #openHandler(event) {
        this.socket.send({
            action: 'InitialMessage',
            username: this.username,
            data: null,
        });
    }

    #beforeunloadHandler() {
        window.addEventListener('beforeunload', (e) => {
            this.socket.send({
                action: 'BeforeUnload',
                username: this.username,
                data: null,
            })
        })
    }

    #messageHandler(event, data) {
        if (this.actionHandlers[data.action]) {
            this.actionHandlers[data.action].call(this, data);
        } else {
            console.error('Invalid action received: ', data);
        }
    }

    #failHandler(data) {
        console.log('#failHandler', data);
    }

    #roomsHandler(data) {
        data.data.forEach((roomData) => {
            this.room = new Room(roomData.name, roomData.userCount);

            const roomNode = this.room.render(this.roomsElement);

            roomNode.addEventListener('click', () => {
                this.#joinRoomHandler({
                    action: 'JoinRoom',
                    username: this.username,
                    data: roomData.name
                })
            })
        });
    }

    #joinRoomHandler(data) {
        this.socket.send(data)
    }

    #joinedRoomHandler(data) {
        const joinedRoom = Array.from(this.roomsElement.children).find((element) => element.firstChild.textContent === data.data.name)

        if (joinedRoom) {
            joinedRoom.firstChild.nextSibling.textContent ++
        }

        if (data.username === this.username) {
            this.room = new Room(data.data.name, data.data.userCount);
            this.messagesElement.innerHTML = ''
            this.room.renderMessages(this.messagesElement, data.data.messages)
        }
    }

    #sentMessageHandler(data) {
        this.room = new Room();
        this.room.renderMessages(this.messagesElement, [data.data])
    }

    #sendMessageHandler(data) {
        this.socket.send(data)
    }

    #leftRoomHandler(data) {
        const leftRoom = Array.from(this.roomsElement.children).find((element) => element.firstChild.textContent === data.data.name)

        if (leftRoom) {
            leftRoom.firstChild.nextSibling.textContent --
        }
    }

    #updateRoomInfoHandler(data) {
        this.roomElement.innerHTML = ''

        this.room = new Room();
        data.data.usernames.forEach((user) => {
            this.room.renderUsers(this.roomElement, user)
        })
    }

    #createRoomHandler(data) {
        this.socket.send(data)
    }

    #createdRoomHandler(data) {
       this.#roomsHandler(data)
    }
 }

export default Chat;