class Message {
    username;
    text;
    createdAt;

    constructor(username, text, type, createdAt) {
        this.username = username
        this.text = text
        this.createdAt = createdAt
    }

    render(wrapper) {
        const messageWrapper = document.createElement('div');
        messageWrapper.classList.add('messageWrapper')

        const messageUser = document.createElement('span');
        messageUser.classList.add('messageUser')
        const messageUserText = document.createTextNode(this.username)
        messageUser.appendChild(messageUserText)
        messageWrapper.appendChild(messageUser)

        const messageTextNode = document.createElement('span');
        messageTextNode.classList.add('messageText')
        const messageText = document.createTextNode(this.text)
        messageTextNode.appendChild(messageText)
        messageWrapper.appendChild(messageTextNode)

        const messageCreatedAt = document.createElement('span');
        messageCreatedAt.classList.add('messageCreatedAt')
        const messageCreatedAtText = document.createTextNode(luxon.DateTime.fromISO(this.createdAt).toFormat('cccc, d MMM yyyy, HH:mm'))
        messageCreatedAt.appendChild(messageCreatedAtText)
        messageWrapper.appendChild(messageCreatedAt)

        wrapper.appendChild(messageWrapper);

        wrapper.scrollTop = wrapper.scrollHeight + messageWrapper.scrollHeight

        wrapper.style.overflowY = 'scroll';
    }
}

export default Message;