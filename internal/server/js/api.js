let isPaused = false;
let hasSequence = false;

function updateState(state) {
    let sequenceName = state.Sequence;
    isPaused = state.IsPaused;
    hasSequence = state.HasSequence;
    updatePauseButtonLabel();
    let pausedInfo = (state.IsPaused) ? " - (paused)" : "";
    let content = sequenceName + pausedInfo;
    document.getElementById("active-sequence").classList.remove("isActive");
    if (hasSequence) {
        document.getElementById("active-sequence").classList.add("isActive");
    }

    document.getElementById("active-sequence").innerHTML = content;
}

function updatePauseButtonLabel() {
    document.getElementById("pause-button").innerHTML = (isPaused) ? "RESUME" : "PAUSE";
}

function clearSequence() {
    wsSend({Type: "SEQUENCE_CLEARSEQUENCE", Payload: {}})
}

function setConfigSequence(pageId, sequenceIndex) {
    wsSend({Type: "SEQUENCE_SETCONFIGSEQUENCE", Payload: {PageId: pageId, SequenceIndex: sequenceIndex}})
}

function pause() {
    wsSend({Type: "SEQUENCE_PAUSE", Payload: {}})
}

function togglePause() {
    pause()
}

function abort() {
    wsSend({Type: "SEQUENCE_ABORT", Payload: {}})
}

function showError(text) {
    document.getElementById("latest-error").innerHTML = text;
}

function clearError() {
    document.getElementById("latest-error").innerHTML = "";
}
