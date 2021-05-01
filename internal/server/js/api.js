let isPaused = false;
let hasSequence = false;

function updateState() {
    clearError();
    fetch('/state')
        .then(handleErrors)
        .then(function (response) {
            return response.json();
        })
        .then((state) => {
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
        })
        .catch((err) => {
                console.error(err)
            }
        );
}

function updatePauseButtonLabel() {
    document.getElementById("pause-button").innerHTML = (isPaused) ? "RESUME" : "PAUSE";
}

function setSequence(sequence) {
    clearError();
    fetch('/set/' + sequence, {method: 'POST'})
        .then(handleErrors)
        .then(updateState)
        .catch((err) => {
                console.error(err)
            }
        )
}

function pause() {
    clearError();
    fetch('/pause', {method: 'POST'})
        .then(handleErrors)
        .then(updateState)
        .catch((err) => {
                console.error(err)
            }
        )
}

function togglePause() {
    pause()
}

function abort() {
    fetch('/abort', {method: 'POST'})
        .then(handleErrors)
        .then(updateState)
        .catch((err) => {
                console.error(err)
            }
        )
}
function handleErrors(response) {
    if (!response.ok) {
        response.text().then(function (text) {
            document.getElementById("latest-error").innerHTML = text;
        });
        throw Error(response.statusText);
    }
    return response;
}

function clearError() {
    document.getElementById("latest-error").innerHTML = "";
}

window.addEventListener("load", function (event) {
    updateState();
});