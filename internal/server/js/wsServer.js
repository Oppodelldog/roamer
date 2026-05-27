function wsSend(msg) {
    ws.send(JSON.stringify(msg));
}

let connected = false;

// server -> client
const seqState = "SEQUENCE_STATE"
const seqSaveResult = "SEQUENCE_SAVE_RESULT"
const seqFormatResult = "SEQUENCE_FORMAT_RESULT"
const seqValidateResult = "SEQUENCE_VALIDATE_RESULT"
const seqElementEvent = "SEQUENCE_ELEMENT_EVENT"
const recorderState = "RECORDER_STATE"
const recorderInputEvent = "RECORDER_INPUT_EVENT"
const remoteInfo = "REMOTE_INFO"
const soundSettings = "SOUND_SETTINGS"
const inputModeState = "INPUT_MODE_STATE"
const logMessage = "LOG_MESSAGE"
const roamerConfig = "CONFIG"

// client -> server
const seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
const seqClearSequence = "SEQUENCE_CLEARSEQUENCE"
const seqPause = "SEQUENCE_PAUSE"
const seqAbort = "SEQUENCE_ABORT"
const seqReleaseInputs = "SEQUENCE_RELEASE_INPUTS"
const seqSave = "SEQUENCE_SAVE"
const seqFormat = "SEQUENCE_FORMAT"
const seqValidate = "SEQUENCE_VALIDATE"
const seqNew = "SEQUENCE_NEW"
const seqDelete = "SEQUENCE_DELETE"
const seqDuplicate = "SEQUENCE_DUPLICATE"
const seqMove = "SEQUENCE_MOVE"
const loadSoundSettings = "LOAD_SOUND_SETTINGS"
const setSoundVolume = "SET_SOUND_VOLUME"
const setMainSoundVolume = "SET_MAIN_SOUND_VOLUME"
const setInputMode = "SET_INPUT_MODE"
const recorderStart = "RECORDER_START"
const recorderStop = "RECORDER_STOP"
const remoteMacroSave = "REMOTE_MACRO_SAVE"
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
                    case inputModeState:
                        updateInputMode(data.Payload)
                        break;
                    case logMessage:
                        appendLogMessage(data.Payload)
                        break;
                    case seqState:
                        updateState(data.Payload)
                        break
                    case seqElementEvent:
                        appendPlaybackCommand(data.Payload)
                        break
                    case recorderState:
                        updateRecorderState(data.Payload)
                        break
                    case recorderInputEvent:
                        appendRecorderCommand(data.Payload)
                        break
                    case remoteInfo:
                        updateRemoteInfo(data.Payload)
                        break
                    case seqSaveResult:
                        let payload = data.Payload;
                        respondSaveSequence(payload.PageId, payload.SequenceIndex, payload.Sequence, payload.Meta, payload.Success, payload.Error)
                        break;
                    case seqFormatResult:
                        let formatPayload = data.Payload;
                        respondFormatSequence(formatPayload.PageId, formatPayload.SequenceIndex, formatPayload.Sequence, formatPayload.Meta, formatPayload.Success, formatPayload.Error)
                        break;
                    case seqValidateResult:
                        let validatePayload = data.Payload;
                        respondValidateSequence(validatePayload.PageId, validatePayload.SequenceIndex, validatePayload.Sequence, validatePayload.Meta, validatePayload.Success, validatePayload.Error)
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

function setInputExecutionMode(dryRun) {
    wsSend({Type: setInputMode, Payload: {DryRun: dryRun}})
}

function startRecorder() {
    wsSend({Type: recorderStart, Payload: {}})
}

function stopRecorder() {
    wsSend({Type: recorderStop, Payload: {}})
}

function saveRemoteMacro(pageId, caption, sequence) {
    wsSend({Type: remoteMacroSave, Payload: {PageId: pageId, Caption: caption, Sequence: sequence}})
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

function saveSequence(pageId, sequenceIndex, caption, icon, sequence) {
    wsSend({
        Type: seqSave,
        Payload: {PageId: pageId, SequenceIndex: sequenceIndex, Caption: caption, Icon: icon, Sequence: sequence}
    })
}

function formatSequence(pageId, sequenceIndex, sequence) {
    wsSend({
        Type: seqFormat,
        Payload: {PageId: pageId, SequenceIndex: sequenceIndex, Sequence: sequence}
    })
}

function validateSequence(pageId, sequenceIndex, sequence) {
    wsSend({
        Type: seqValidate,
        Payload: {PageId: pageId, SequenceIndex: sequenceIndex, Sequence: sequence}
    })
}

function clearSequence() {
    wsSend({Type: seqClearSequence, Payload: {}})
}

function deleteSequence(pageId, sequenceIndex) {
    wsSend({Type: seqDelete, Payload: {PageId: pageId, SequenceIndex: sequenceIndex}})
}

function duplicateSequence(pageId, sequenceIndex) {
    wsSend({Type: seqDuplicate, Payload: {PageId: pageId, SequenceIndex: sequenceIndex}})
}

function moveSequence(pageId, sequenceIndex, offset) {
    wsSend({Type: seqMove, Payload: {PageId: pageId, SequenceIndex: sequenceIndex, Offset: offset}})
}

function pause() {
    wsSend({Type: seqPause, Payload: {}})
}


function abort() {
    wsSend({Type: seqAbort, Payload: {}})
}

function releaseInputs() {
    wsSend({Type: seqReleaseInputs, Payload: {}})
}

function savePages(pages) {
    wsSend({Type: pagesSave, Payload: {Pages: pages}})
}
