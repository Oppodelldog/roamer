class WebsocketMessageReader {
    constructor() {
        this.messageBuffer = "";
        this.readMessages = function (data) {
            const messages = [];
            let buffer = "";
            this.messageBuffer += data;
            for (let i = 0; i < this.messageBuffer.length; i++) {
                const char = this.messageBuffer[i];
                buffer += char;
                if (char === "\n") {
                    messages.push(buffer);
                    buffer = "";
                }
            }
            this.messageBuffer = buffer;

            return messages;
        }
    }
}