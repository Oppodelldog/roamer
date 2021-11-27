function wsSend(msg) {
    ws.send(JSON.stringify(msg));
}

function connectToServer() {
    ws = new WebSocket(wsUrl);
    ws.onopen = function (evt) {
        console.log("opened channel");
    };
    ws.onmessage = function (evt) {
        const messages = wsMessageReader.readMessages(evt.data);
        for (const k in messages) {
            if (!messages.hasOwnProperty(k)) {
                continue;
            }
            const message = messages[k];
            const data = JSON.parse(message);
            switch (data.Type) {
                case "SERVER_TEXT":
                    const serverMessage = data.Payload.Text;
                    console.log("server says:" + serverMessage)
                    break
                case "SEQUENCE_STATE":
                    updateState(data.Payload)
                    break
            }
        }
        ws.onclose = function (evt) {
            console.log("closed channel");
            console.log(evt)
        };
        ws.onerror = function (evt) {
            console.log("error");
            console.log(evt)
        };
    };
}
