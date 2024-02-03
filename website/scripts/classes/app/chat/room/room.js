import Message from "./message.js";

class Room {
    name = ''
    count = null

    element;
    nameElement;
    countElement;

    constructor(name = null, count = null) {
        this.name = name
        this.count = count
    }

    render(wrapper) {
        const roomWrapper = document.createElement('div');
        roomWrapper.classList.add('roomWrapper')

        const roomName = document.createElement('span');
        const roomNameText = document.createTextNode(this.name)
        roomName.appendChild(roomNameText)
        roomWrapper.appendChild(roomName)

        const roomCount = document.createElement('span');
        roomCount.classList.add('roomCount')
        const roomCountText = document.createTextNode(this.count)
        roomCount.appendChild(roomCountText)
        roomWrapper.appendChild(roomCount)

        this.element = roomWrapper;
        this.nameElement = roomName;
        this.countElement = roomCount;

        wrapper.appendChild(roomWrapper);

        return roomWrapper
    }

    renderUsers(wrapper, user) {
        const participantWrapper = document.createElement('span')
        const participantText = document.createTextNode(user)
        participantWrapper.appendChild(participantText)

        wrapper.appendChild(participantWrapper)

        wrapper.style.opacity = '1'
    }

    renderMessages(wrapper, messages) {
        messages.forEach((message) => {
            new Message(message.username, message.text, message.type, message.createdAt).render(wrapper)
        })
    }

}

export default Room