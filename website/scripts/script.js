import App from "./classes/app/app.js";

const createNewRoom = document.querySelector('.createRoom')
const createRoomName = document.querySelector('.createRoomName')
const sendMessage = document.querySelector('.sendMessage')
const newMessage = document.querySelector('.newMessage')


let username = localStorage.getItem('ws-username')
while (!username) {
    username = prompt("What Is Your Username?")

    if (username) {
        localStorage.setItem('ws-username', username)
    }
}

const app = new App(username)

createNewRoom.addEventListener('submit', (e) => {
    e.preventDefault();

    if (createRoomName.value) {
        app.chat.actionHandlers.CreateRoom({
            action: 'CreateRoom',
            username: username,
            data: createRoomName.value
        })
    }

    createRoomName.value = ''
})

sendMessage.addEventListener('submit', (e) => {
    e.preventDefault()

    if (newMessage.value) {
        app.chat.actionHandlers.SendMessage({
            action: 'SendMessage',
            username: username,
            data: newMessage.value
        })
    }

    newMessage.value = ''
})