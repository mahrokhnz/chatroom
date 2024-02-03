class Socket {
    socket;

    handlers = {
        open: [],
        close: [],
        error: [],
        message: [],
    };

    constructor(url) {
        this.socket = new WebSocket(url);

        this.socket.onopen = this.#openHandler.bind(this);
        this.socket.onclose = this.#closeHandler.bind(this);
        this.socket.onerror = this.#errorHandler.bind(this);
        this.socket.onmessage = this.#messageHandler.bind(this);
    }

    #openHandler(event) {
        this.handlers.open.forEach((handler) => {
            handler(event);
        });
    }

    #closeHandler(event) {
        if (event.wasClean) {
            console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`)
        } else {
            console.error('[close] Connection died')
        }

        this.handlers.close.forEach((handler) => {
            handler(event);
        });
    }

    #errorHandler(error) {
        console.error(`[error], ${error}`);

        this.handlers.error.forEach((handler) => {
            handler(error);
        });
    }

    #messageHandler(event) {
        try {
            const data = JSON.parse(event.data)

            this.handlers.message.forEach((handler) => {
                handler(event, data);
            });
        } catch (error) {
            this.handlers.error.forEach((handler) => {
                handler(error);
            });
        }
    }

    send(data) {
        this.socket.send(JSON.stringify(data));
    }

    addHandler(event, handler) {
        this.handlers[event].push(handler);
    }

    removeHandler(event, handler) {
        this.handlers[event] = this.handlers[event].filter((h) => h !== handler);
    }
}

export default Socket;