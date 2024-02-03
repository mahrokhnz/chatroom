import Socket from "./chat/socket.js";
import Chat from "./chat/chat.js";

class App {
    username;
    socket;
    chat;

    constructor(username) {
        this.username = username;

        this.socket = new Socket("ws://192.168.3.4:8080/ws");

        this.chat = new Chat(this.socket, this.username);
    }
}

export default App;