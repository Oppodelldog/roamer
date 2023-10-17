function wsSend(msg) {
    ws.send(JSON.stringify(msg));
}

let connected = false;

// server -> client
const seqState = "SEQUENCE_STATE"
const seqSaveResult = "SEQUENCE_SAVE_RESULT"
const soundSettings = "SOUND_SETTINGS"
const logMessage = "LOG_MESSAGE"
const roamerConfig = "CONFIG"

// client -> server
const seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
const seqClearSequence = "SEQUENCE_CLEARSEQUENCE"
const seqPause = "SEQUENCE_PAUSE"
const seqAbort = "SEQUENCE_ABORT"
const seqSave = "SEQUENCE_SAVE"
const seqNew = "SEQUENCE_NEW"
const seqDelete = "SEQUENCE_DELETE"
const loadSoundSettings = "LOAD_SOUND_SETTINGS"
const setSoundVolume = "SET_SOUND_VOLUME"
const setMainSoundVolume = "SET_MAIN_SOUND_VOLUME"
const pageNew = "PAGE_NEW"
const pageDelete = "PAGE_DELETE"
const pagesSave = "PAGES_SAVE"

function connectToServer() {
    try {
        ws = new WebSocket(wsUrl);
        ws.onopen = function (evt) {
            console.log("websocket connected");
            connected = true;
            updateConnectionStatus(connected)
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
                    case roamerConfig:
                        updateAppConfig(data.Payload)
                        break;
                    case soundSettings:
                        updateSoundSettings(data.Payload)
                        break;
                    case logMessage:
                        appendLogMessage(data.Payload)
                        break;
                    case seqState:
                        updateState(data.Payload)
                        break
                    case seqSaveResult:
                        let payload = data.Payload;
                        respondSaveSequence(payload.PageId, payload.SequenceIndex, payload.Sequence, payload.Success)
                        break;
                }
            }
            ws.onclose = function (evt) {
                console.log("websocket closed");
                console.log(evt)
                connected = false;
                updateConnectionStatus(connected)
            };
            ws.onerror = function (evt) {
                console.log("websocket error");
                console.log(evt)
                connected = false;
                updateConnectionStatus(connected)
            };
        };
    } catch (err) {
        console.log("websocket cannot connect:", err)
        connected = false;
        updateConnectionStatus(connected)
    }
}

function reconnect() {
    if (connected) return;

    if (ws != null) {
        console.log("websocket reconnecting")
        ws.close();
        connectToServer();
    }
}

setInterval(reconnect, 2000);

function updateState(state) {
    updatePlaybackState(state)
}

function setVolume(id, volume) {
    wsSend({Type: setSoundVolume, Payload: {Id: id, Value: volume}})
}

function setMainVolume(volume) {
    wsSend({Type: setMainSoundVolume, Payload: {Value: volume}})
}

function querySoundSettings() {
    wsSend({Type: loadSoundSettings, Payload: {}})
}

function setConfigSequence(pageId, sequenceIndex) {
    wsSend({Type: seqSetConfigSequence, Payload: {PageId: pageId, SequenceIndex: sequenceIndex}})
}

function createNewPage() {
    wsSend({Type: pageNew, Payload: {}})
}

function deletePage(pageId) {
    wsSend({Type: pageDelete, Payload: {PageId: pageId}})
}

function createNewSequence(pageId) {
    wsSend({Type: seqNew, Payload: {PageId: pageId}})
}

function saveSequence(pageId, sequenceIndex, caption, sequence) {
    wsSend({
        Type: seqSave,
        Payload: {PageId: pageId, SequenceIndex: sequenceIndex, Caption: caption, Sequence: sequence}
    })
}

function clearSequence() {
    wsSend({Type: seqClearSequence, Payload: {}})
}

function deleteSequence(pageId, sequenceIndex) {
    wsSend({Type: seqDelete, Payload: {PageId: pageId, SequenceIndex: sequenceIndex}})
}

function pause() {
    wsSend({Type: seqPause, Payload: {}})
}


function abort() {
    wsSend({Type: seqAbort, Payload: {}})
}


function savePages(pages) {
    wsSend({Type: pagesSave, Payload: {Pages: pages}})
}