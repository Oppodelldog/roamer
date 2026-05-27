const wsUrl = 'ws://' + window.location.host + '/ws'
let ws = null
let wsMessageReader = new WebsocketMessageReader()
let app = null;

function initApp() {
    app = new Vue({
        el: '#app',
        data: {
            hasConfig: false,
            config: null,
            remoteMode: window.location.pathname === "/remote",
            remoteInfo: {Urls: []},
            remoteQr: {Open: false, Index: 0},
            remoteRecorder: {Open: false, Name: "", Saving: false, Error: ""},
            showPage: false,
            showVerticalSlide: false,
            pageEditor: false,
            macroEditor: false,
            showSounds: false,
            showDiagnostics: false,
            showCommandHelp: true,
            currentPageId: null,
            connection: {isConnected: false},
            inputMode: {Mode: "live", DryRun: false},
            recorder: {Active: false, Sequence: "", Count: 0, Error: ""},
            soundSettings: {Sessions: [], MainSession: {}},
            soundLoading: false,
            deletePageClicks: {},
            sequenceSaveStatus: {},
            sequenceValidationStatus: {},
            sequenceValidationTimers: {},
            sequenceExecutionErrors: {},
            playbackCommands: [],
            playbackCommandId: 0,
            playbackCommandTimers: {},
            recorderCommands: [],
            recorderCommandId: 0,
            playbackTimeline: {
                Visible: false,
                Current: "",
                Next: "",
                StepIndex: 0,
                TotalSteps: 0,
                Progress: 0,
                RemainingMs: 0,
                Cycle: 0,
                Looping: false
            },
            lastPlaybackStartKey: "",
            playbackState: {
                PageId: "",
                SequenceIndex: -1,
                Caption: "",
                State: "idle",
                Error: "",
                HasSequence: false,
                IsPlaying: false,
                QueuedPageId: "",
                QueuedSequenceIndex: -1,
                QueuedPageTitle: "",
                QueuedCaption: "",
            },
            iconOptions: [
                {Class: "ico-none", Label: "None"},
                {Class: "ico-get-mouse-pos", Label: "Mouse Pos"},
                {Class: "ico-click-left", Label: "Left Click"},
                {Class: "ico-walk", Label: "Walk"},
                {Class: "ico-run", Label: "Run"},
                {Class: "ico-run-stop", Label: "Run Stop"},
                {Class: "ico-jump", Label: "Jump"},
                {Class: "ico-arm", Label: "Arm"},
                {Class: "ico-unarm", Label: "Unarm"},
                {Class: "ico-gather-tree", Label: "Gather"},
                {Class: "ico-grill", Label: "Grill"},
                {Class: "ico-kayak-forward", Label: "Paddle"},
                {Class: "ico-kayak-backward", Label: "Paddle Back"},
                {Class: "icon-repair", Label: "Repair"},
                {Class: "ico-abort", Label: "Abort"},
                {Class: "ico-pause-button", Label: "Pause"}
            ],
            logMessages: []
        },
        methods: {
            currentPage: function () {
                return this.config.Pages[this.currentPageId];
            },
            pageActionCount: function (page) {
                if (!page || !page.Actions) {
                    return 0
                }

                return page.Actions.length
            },
            connectionLabel: function () {
                return this.connection.isConnected ? "Connected" : "Offline"
            },
            modeLabel: function () {
                return this.remoteMode ? "Remote" : "Local"
            },
            remoteUrl: function () {
                if (this.remoteInfo.Urls && this.remoteInfo.Urls.length > 0) {
                    return this.remoteInfo.Urls[0]
                }

                return window.location.origin + "/remote"
            },
            remoteUrls: function () {
                if (this.remoteInfo.Urls && this.remoteInfo.Urls.length > 0) {
                    return this.remoteInfo.Urls
                }

                return [this.remoteUrl()]
            },
            remoteTargets: function () {
                if (this.remoteInfo.Targets && this.remoteInfo.Targets.length > 0) {
                    return this.remoteInfo.Targets
                }

                return this.remoteUrls().map(function (url) {
                    return {Url: url, QrCode: ""}
                })
            },
            activeRemoteTarget: function () {
                let targets = this.remoteTargets()
                let idx = Math.min(Math.max(this.remoteQr.Index, 0), targets.length - 1)

                return targets[idx]
            },
            activeRemoteUrl: function () {
                return this.activeRemoteTarget().Url
            },
            remoteQrPositionLabel: function () {
                let targets = this.remoteTargets()
                let idx = Math.min(Math.max(this.remoteQr.Index, 0), targets.length - 1)

                return (idx + 1) + " / " + targets.length
            },
            remoteUrlHostLabel: function (url) {
                try {
                    return new URL(url).host
                } catch (e) {
                    return url
                }
            },
            inputModeLabel: function () {
                return this.inputMode.DryRun ? "Dry Run" : "Live"
            },
            playbackStatusLabel: function () {
                if (this.playbackState.State === "error") {
                    return "Execution error"
                }

                if (this.playbackState.State === "stopped") {
                    return "Stopped"
                }

                if (!this.playbackState.HasSequence) {
                    return this.showPage ? "Macro deck" : "Game profiles"
                }

                if (this.playbackState.State === "paused") {
                    return "Paused"
                }

                if (this.playbackState.State === "looping") {
                    return "Looping"
                }

                if (this.playbackState.QueuedCaption) {
                    return "Playing - queued next"
                }

                return "Playing"
            },
            activeMacroLabel: function () {
                if (this.playbackState.State === "stopped" && this.playbackState.Caption) {
                    return this.playbackState.PageTitle ? this.playbackState.PageTitle + " / " + this.playbackState.Caption : this.playbackState.Caption
                }

                if (!this.playbackState.HasSequence || !this.playbackState.Caption) {
                    return this.showPage ? "Choose a macro to run" : "Choose a profile"
                }

                if (!this.playbackState.PageTitle) {
                    return this.playbackState.Caption
                }

                return this.playbackState.PageTitle + " / " + this.playbackState.Caption
            },
            playbackPanelClass: function () {
                return {
                    'playback-panel': true,
                    'is-error': this.playbackState.State === "error"
                }
            },
            deckShellClass: function () {
                let classes = {
                    'deck-shell': true
                }

                classes[this.activeThemeClass()] = true

                return classes
            },
            activeThemeClass: function () {
                if (this.showPage && this.currentPageId && this.config.Pages[this.currentPageId]) {
                    return this.config.Pages[this.currentPageId].ThemeClass || "theme-default"
                }

                return "theme-default"
            },
            newPage: function () {
                createNewPage();
            },
            removePage: function (pageId) {
                Vue.set(this.deletePageClicks, pageId, this.getPageDeletionStatus(pageId) + 1)
                if (this.deletePageClicks[pageId] === 2) {
                    delete this.deletePageClicks[pageId]
                    deletePage(pageId);
                }
            },
            selectPage: function (pageKey) {
                this.currentPageId = pageKey
                this.showPage = true
                this.applyTheme()
            },
            newSequence: function () {
                createNewSequence(this.currentPageId);
            },
            startAction: function (actionIdx) {
                let action = this.config.Pages[this.currentPageId].Actions[actionIdx]
                if (this.isActionInvalid(action)) {
                    return
                }

                if (this.isActionActive(actionIdx)) {
                    abort()
                    return
                }

                this.clearSequenceExecutionError(this.currentPageId, actionIdx)
                setConfigSequence(this.currentPageId, actionIdx)
            },
            isActionActive: function (actionIdx) {
                return this.playbackState.PageId === this.currentPageId &&
                    this.playbackState.SequenceIndex === actionIdx &&
                    (this.playbackState.IsPlaying || this.playbackState.HasSequence)
            },
            isActionPaused: function (actionIdx) {
                return this.isActionActive(actionIdx) && this.playbackState.State === "paused"
            },
            isActionQueued: function (actionIdx) {
                return this.playbackState.QueuedPageId === this.currentPageId &&
                    this.playbackState.QueuedSequenceIndex === actionIdx
            },
            isActionStopped: function (actionIdx) {
                return this.playbackState.PageId === this.currentPageId &&
                    this.playbackState.SequenceIndex === actionIdx &&
                    this.playbackState.State === "stopped"
            },
            isActionInvalid: function (action) {
                return action && action.Meta && action.Meta.Valid === false
            },
            hasSequenceExecutionError: function (pageId, sequenceIndex) {
                return !!this.sequenceExecutionErrors[this.sequenceSaveKey(pageId, sequenceIndex)]
            },
            getSequenceExecutionError: function (pageId, sequenceIndex) {
                return this.sequenceExecutionErrors[this.sequenceSaveKey(pageId, sequenceIndex)] || ""
            },
            clearSequenceExecutionError: function (pageId, sequenceIndex) {
                Vue.delete(this.sequenceExecutionErrors, this.sequenceSaveKey(pageId, sequenceIndex))
            },
            macroStateLabel: function (actionIdx, action) {
                if (this.isActionInvalid(action)) {
                    return "Needs edit"
                }

                if (this.hasSequenceExecutionError(this.currentPageId, actionIdx)) {
                    return "Execution error"
                }

                if (this.isActionQueued(actionIdx)) {
                    return "Queued"
                }

                if (this.isActionStopped(actionIdx)) {
                    return "Stopped"
                }

                if (!this.isActionActive(actionIdx)) {
                    return "Tap to start"
                }

                if (this.isActionPaused(actionIdx)) {
                    return "Paused - tap to stop"
                }

                if (this.isActionLooping(actionIdx)) {
                    return "Looping - tap to stop"
                }

                return "Playing - tap to stop"
            },
            macroCardTitle: function (action) {
                if (this.isActionInvalid(action) && action.Meta.Error) {
                    return action.Meta.Error
                }

                let executionError = this.getSequenceExecutionError(this.currentPageId, this.currentPage().Actions.indexOf(action))
                if (executionError) {
                    return executionError
                }

                return action.Caption || ""
            },
            macroMeta: function (action) {
                if (!action) {
                    return []
                }

                if (this.isActionInvalid(action)) {
                    return ["Invalid"]
                }

                if (action.Meta && action.Meta.Labels) {
                    return action.Meta.Labels
                }

                if (!action.Sequence) {
                    return []
                }

                let meta = []
                let sequence = action.Sequence
                let normalized = ";" + sequence.replace(/\s+/g, " ").trim() + ";"

                if (/;\s*L\s*;/.test(normalized)) {
                    meta.push("Loop")
                }

                if (/\b(MP|MM|SM|LD|LU|RD|RU)\b/.test(sequence)) {
                    meta.push("Mouse")
                }

                let keyDowns = sequence.match(/\bKD\s+[A-Z0-9_]+/g) || []
                let keyUps = sequence.match(/\bKU\s+[A-Z0-9_]+/g) || []
                if (keyDowns.length > keyUps.length) {
                    meta.push("Holds keys")
                }

                let duration = this.estimatedDuration(sequence)
                if (duration !== "") {
                    meta.push(duration)
                }

                return meta
            },
            estimatedDuration: function (sequence) {
                let matches = sequence.match(/\bW\s+([0-9.]+)(ms|s|m)\b/g) || []
                let totalMs = 0

                for (let i = 0; i < matches.length; i++) {
                    let match = matches[i].match(/\bW\s+([0-9.]+)(ms|s|m)\b/)
                    if (!match) {
                        continue
                    }

                    let value = parseFloat(match[1])
                    if (Number.isNaN(value)) {
                        continue
                    }

                    if (match[2] === "ms") {
                        totalMs += value
                    } else if (match[2] === "s") {
                        totalMs += value * 1000
                    } else if (match[2] === "m") {
                        totalMs += value * 60000
                    }
                }

                if (totalMs <= 0) {
                    return ""
                }

                if (totalMs < 1000) {
                    return "~" + Math.round(totalMs) + "ms"
                }

                if (totalMs < 60000) {
                    return "~" + Math.round(totalMs / 1000) + "s"
                }

                return "~" + Math.round(totalMs / 60000) + "m"
            },
            pauseControlLabel: function () {
                return this.playbackState.State === "paused" ? "Resume" : "Pause"
            },
            isActionLooping: function (actionIdx) {
                return this.isActionActive(actionIdx) && this.playbackState.State === "looping"
            },
            macroButtonClass: function (action, actionIdx) {
                let classes = {
                    'macro-action': true,
                    'is-playing': this.isActionActive(actionIdx),
                    'is-paused': this.isActionPaused(actionIdx),
                    'is-looping': this.isActionLooping(actionIdx),
                    'is-queued': this.isActionQueued(actionIdx),
                    'is-stopped': this.isActionStopped(actionIdx),
                    'is-error': this.isActionInvalid(action),
                    'is-execution-error': this.hasSequenceExecutionError(this.currentPageId, actionIdx),
                    'is-idle': !this.isActionActive(actionIdx) && !this.isActionQueued(actionIdx) && !this.isActionStopped(actionIdx)
                }

                classes[action.Icon || 'ico-none'] = true

                return classes
            },
            iconLabel: function (iconClass) {
                for (let i = 0; i < this.iconOptions.length; i++) {
                    if (this.iconOptions[i].Class === iconClass) {
                        return this.iconOptions[i].Label
                    }
                }

                return "None"
            },
            updateConfig: function (config) {
                this.config = config
                this.hasConfig = true
                this.sequenceValidationStatus = {}
                this.sequenceValidationTimers = {}
                if (this.remoteRecorder.Open && this.remoteRecorder.Saving) {
                    this.remoteRecorder = {Open: false, Name: "", Saving: false, Error: ""}
                }
                this.applyTheme()
            },
            updateConnectionStatus: function (isConnected) {
                this.connection.isConnected = isConnected
            },
            updateInputMode: function (state) {
                this.inputMode = state
            },
            updateRemoteInfo: function (info) {
                this.remoteInfo = info
                if (this.remoteQr.Index >= this.remoteTargets().length) {
                    this.remoteQr.Index = 0
                }
            },
            toggleInputMode: function () {
                setInputExecutionMode(!this.inputMode.DryRun)
            },
            updateRecorderState: function (state) {
                this.recorder = state
                if (this.remoteMode && this.remoteRecorder.Open && !state.Active && state.Sequence && !this.remoteRecorder.Name) {
                    this.remoteRecorder.Name = "Recorded Macro"
                }
                if (state.Error) {
                    this.appendLogMessage("Recorder: " + state.Error)
                }
            },
            recorderStatusLabel: function () {
                if (this.recorder.Active) {
                    return "Recording input..."
                }

                if (this.recorder.Sequence) {
                    return this.recorder.Count + " recorded commands"
                }

                return "Recorder idle"
            },
            toggleRecorder: function (event) {
                if (event && event.currentTarget) {
                    event.currentTarget.blur()
                }

                if (this.recorder.Active) {
                    stopRecorder()
                    return
                }

                startRecorder()
            },
            recorderButtonLabel: function () {
                return this.recorder.Active ? "Stop Recording" : "Start Recording"
            },
            openRemoteRecorder: function () {
                this.remoteRecorder = {Open: true, Name: "", Saving: false, Error: ""}
            },
            closeRemoteRecorder: function () {
                if (this.recorder.Active) {
                    stopRecorder()
                }
                this.remoteRecorder = {Open: false, Name: "", Saving: false, Error: ""}
            },
            toggleRemoteRecording: function () {
                if (this.recorder.Active) {
                    stopRecorder()
                    return
                }

                this.remoteRecorder.Error = ""
                startRecorder()
            },
            remoteRecorderPrimaryLabel: function () {
                if (this.recorder.Active) {
                    return "Stop"
                }

                if (this.recorder.Sequence) {
                    return "Record Again"
                }

                return "Record"
            },
            saveRemoteRecording: function () {
                if (!this.currentPageId || !this.recorder.Sequence || this.remoteRecorder.Saving) {
                    return
                }

                this.remoteRecorder.Saving = true
                this.remoteRecorder.Error = ""
                saveRemoteMacro(this.currentPageId, this.remoteRecorder.Name || "Recorded Macro", this.recorder.Sequence)
            },
            showPageSelection: function () {
                this.pageEditor = false;
                this.macroEditor = false;
                this.showSounds = false;
                this.showVerticalSlide = false;
                this.showPage = false
                this.currentPageId = null
                this.applyTheme()
            },
            openRemoteQr: function () {
                this.remoteQr.Open = true
                if (this.remoteQr.Index >= this.remoteTargets().length) {
                    this.remoteQr.Index = 0
                }
            },
            closeRemoteQr: function () {
                this.remoteQr.Open = false
            },
            selectRemoteQr: function (idx) {
                this.remoteQr.Index = idx
            },
            previousRemoteQr: function () {
                let targets = this.remoteTargets()
                this.remoteQr.Index = (this.remoteQr.Index + targets.length - 1) % targets.length
            },
            nextRemoteQr: function () {
                let targets = this.remoteTargets()
                this.remoteQr.Index = (this.remoteQr.Index + 1) % targets.length
            },
            showSoundSettings: function () {
                if (this.remoteMode) {
                    return
                }

                if (this.showSounds) {
                    this.showSounds = false
                    this.soundLoading = false
                    return
                }

                if (this.soundLoading) {
                    return
                }

                if (this.pageEditor) {
                    savePages(this.config.Pages)
                }

                this.showVerticalSlide = false
                this.pageEditor = false
                this.macroEditor = false
                this.showSounds = true
                this.soundLoading = true
                querySoundSettings()
            },
            toggleDiagnostics: function () {
                if (this.remoteMode) {
                    return
                }

                this.showDiagnostics = !this.showDiagnostics
            },
            changeSoundSessionValue: function (idx) {
                let id = this.soundSettings.Sessions[idx].Id
                let volume = this.soundSettings.Sessions[idx].Value
                setVolume(id, parseFloat(volume))
            },
            changeMainSoundSessionValue: function () {
                let volume = this.soundSettings.MainSession.Value
                setMainVolume(parseFloat(volume))
            },
            updateSoundSettings: function (soundSettings) {
                this.soundSettings = soundSettings;
                this.soundLoading = false;
            },
            appendLogMessage: function (message) {
                this.logMessages.push(message)
                if (this.logMessages.length > 50) {
                    this.logMessages.shift()
                }
            },
            appendPlaybackCommand: function (event) {
                this.updatePlaybackTimeline(event)
                let id = ++this.playbackCommandId
                let waitMs = event.DurationMs || 0
                let durationMs = waitMs > 0 ? waitMs + 1800 : 1800
                this.pushPlaybackCommandsUp()
                let command = {
                    Id: id,
                    Type: "command",
                    Label: event.Label,
                    DisplayLabel: waitMs > 0 ? this.waitCountdownLabel(waitMs) : event.Label,
                    IsWait: waitMs > 0,
                    Slot: 0,
                    DurationMs: durationMs,
                    Style: this.playbackToastStyle(0, durationMs)
                }

                this.playbackCommands.push(command)
                this.prunePlaybackCommands()
                if (waitMs > 0) {
                    this.startWaitCountdown(command, waitMs)
                }

                let app = this
                setTimeout(function () {
                    app.removePlaybackCommand(id)
                }, durationMs + 120)
            },
            appendRecorderCommand: function (event) {
                let id = ++this.recorderCommandId
                let durationMs = 1500
                this.pushRecorderCommandsUp()
                let command = {
                    Id: id,
                    Label: event.Label,
                    Kind: event.Kind || "command",
                    Slot: 0,
                    Style: this.recorderToastStyle(0, durationMs)
                }

                this.recorderCommands.push(command)
                this.pruneRecorderCommands()

                let app = this
                setTimeout(function () {
                    app.removeRecorderCommand(id)
                }, durationMs + 120)
            },
            removeRecorderCommand: function (id) {
                this.recorderCommands = this.recorderCommands.filter(function (command) {
                    return command.Id !== id
                })
            },
            pushRecorderCommandsUp: function () {
                for (let i = 0; i < this.recorderCommands.length; i++) {
                    let command = this.recorderCommands[i]
                    command.Slot = (command.Slot || 0) + 1
                    command.Style = this.recorderToastStyle(command.Slot, 1500)
                }
            },
            recorderToastStyle: function (slot, durationMs) {
                return {
                    animationDuration: durationMs + "ms",
                    "--toast-offset": (slot * 34) + "px"
                }
            },
            pruneRecorderCommands: function () {
                if (this.recorderCommands.length <= 48) {
                    return
                }

                this.recorderCommands = this.recorderCommands.slice(this.recorderCommands.length - 48)
            },
            updatePlaybackTimeline: function (event) {
                if (event.IsLoop) {
                    this.playbackTimeline = {
                        Visible: true,
                        Current: "Loop",
                        Next: event.NextLabel || "",
                        StepIndex: event.StepIndex || event.TotalSteps || 0,
                        TotalSteps: event.TotalSteps || 0,
                        Progress: 100,
                        RemainingMs: 0,
                        Cycle: event.Cycle || 0,
                        Looping: true
                    }
                    return
                }

                this.playbackTimeline = {
                    Visible: true,
                    Current: event.Label || "",
                    Next: event.NextLabel || "",
                    StepIndex: event.StepIndex || 0,
                    TotalSteps: event.TotalSteps || 0,
                    Progress: Math.max(0, Math.min(event.Progress || 0, 100)),
                    RemainingMs: event.RemainingMs || 0,
                    Cycle: event.Cycle || 0,
                    Looping: this.playbackState.State === "looping"
                }
            },
            timelineProgressStyle: function () {
                return {
                    width: this.playbackTimeline.Progress + "%"
                }
            },
            timelineStepLabel: function () {
                if (!this.playbackTimeline.Visible) {
                    return ""
                }

                let prefix = this.playbackTimeline.Looping && this.playbackTimeline.Cycle > 0 ? "Cycle " + this.playbackTimeline.Cycle + " - " : ""
                if (this.playbackTimeline.TotalSteps <= 0) {
                    return prefix + "Step"
                }

                return prefix + this.playbackTimeline.StepIndex + " / " + this.playbackTimeline.TotalSteps
            },
            timelineRemainingLabel: function () {
                let ms = this.playbackTimeline.RemainingMs
                if (!ms || ms <= 0) {
                    return this.playbackTimeline.Looping ? "cycle ending" : "finishing"
                }

                if (ms < 1000) {
                    return "~" + Math.ceil(ms) + "ms left"
                }

                if (ms < 60000) {
                    return "~" + Math.ceil(ms / 1000) + "s left"
                }

                return "~" + Math.ceil(ms / 60000) + "m left"
            },
            appendPlaybackStart: function (state) {
                let id = ++this.playbackCommandId
                let durationMs = 2600
                let command = {
                    Id: id,
                    Type: "start",
                    Label: state.Caption || "Macro started",
                    DisplayLabel: state.Caption || "Macro started",
                    Context: state.PageTitle || "",
                    DurationMs: durationMs,
                    Style: {
                        animationDuration: durationMs + "ms",
                        "--toast-offset": "318px"
                    }
                }

                this.playbackCommands.push(command)
                this.prunePlaybackCommands()

                let app = this
                setTimeout(function () {
                    app.removePlaybackCommand(id)
                }, durationMs + 120)
            },
            removePlaybackCommand: function (id) {
                if (this.playbackCommandTimers[id]) {
                    clearInterval(this.playbackCommandTimers[id])
                    Vue.delete(this.playbackCommandTimers, id)
                }

                this.playbackCommands = this.playbackCommands.filter(function (command) {
                    return command.Id !== id
                })
            },
            startWaitCountdown: function (command, waitMs) {
                let app = this
                let endAt = Date.now() + waitMs
                let interval = setInterval(function () {
                    let remainingMs = Math.max(0, endAt - Date.now())
                    command.DisplayLabel = app.waitCountdownLabel(remainingMs)

                    if (remainingMs <= 0) {
                        command.DisplayLabel = command.Label
                        command.IsWait = false
                        clearInterval(interval)
                        Vue.delete(app.playbackCommandTimers, command.Id)
                    }
                }, 250)

                Vue.set(this.playbackCommandTimers, command.Id, interval)
            },
            waitCountdownLabel: function (remainingMs) {
                if (remainingMs >= 1000) {
                    return "W " + Math.ceil(remainingMs / 1000) + "s"
                }

                return "W " + Math.ceil(remainingMs) + "ms"
            },
            pushPlaybackCommandsUp: function () {
                for (let i = 0; i < this.playbackCommands.length; i++) {
                    let command = this.playbackCommands[i]
                    if (command.Type !== "command") {
                        continue
                    }

                    command.Slot = (command.Slot || 0) + 1
                    command.Style = this.playbackToastStyle(command.Slot, command.DurationMs)
                }
            },
            playbackToastStyle: function (slot, durationMs) {
                return {
                    animationDuration: durationMs + "ms",
                    "--toast-offset": (slot * 38) + "px"
                }
            },
            prunePlaybackCommands: function () {
                if (this.playbackCommands.length <= 80) {
                    return
                }

                this.playbackCommands = this.playbackCommands.slice(this.playbackCommands.length - 80)
            },
            showPagesEditor: function () {
                if (this.remoteMode) {
                    return
                }

                if (this.pageEditor && this.showVerticalSlide) {
                    this.closeVerticalSlide()
                    savePages(this.config.Pages)
                    return
                }

                this.openVerticalSlide("pages")
            },
            initPageDeleteStatus: function (key) {
                Vue.set(this.deletePageClicks, key, 0)
            },
            getPageDeletionStatus: function (key) {
                if (this.deletePageClicks[key] === undefined) {
                    this.initPageDeleteStatus(key)
                }

                return this.deletePageClicks[key];
            },
            pageDeleteLabel: function (key) {
                return this.getPageDeletionStatus(key) === 1 ? "Really?" : "Delete"
            },
            ensurePageTheme: function (page) {
                if (!page.Theme) {
                    Vue.set(page, "Theme", this.defaultTheme())
                }

                return page.Theme
            },
            defaultTheme: function () {
                return {
                    BackgroundImage: "/img/background/index.jpg",
                    BackgroundColor: "#141f35",
                    AccentColor: "#ffe36d",
                    CardColor: "#ecf1e8",
                    CardTextColor: "#141714"
                }
            },
            activeTheme: function () {
                if (this.showPage && this.currentPageId && this.config.Pages[this.currentPageId]) {
                    return this.ensurePageTheme(this.config.Pages[this.currentPageId])
                }

                return this.defaultTheme()
            },
            applyTheme: function () {
                let theme = this.activeTheme()
                let root = document.documentElement

                root.style.setProperty("--roamer-bg-color", theme.BackgroundColor || "#141f35")
                root.style.setProperty("--roamer-accent", theme.AccentColor || "#ffe36d")
                root.style.setProperty("--roamer-card-bg", theme.CardColor || "#ecf1e8")
                root.style.setProperty("--roamer-card-text", theme.CardTextColor || "#141714")

                if (theme.BackgroundImage) {
                    root.style.setProperty("--roamer-bg-image", "url('" + theme.BackgroundImage + "')")
                } else {
                    root.style.setProperty("--roamer-bg-image", "none")
                }
            },
            showMacroEditor: function () {
                if (this.remoteMode) {
                    return
                }

                if (this.macroEditor) {
                    this.macroEditor = false
                    return
                }

                if (this.pageEditor) {
                    savePages(this.config.Pages)
                }

                this.showVerticalSlide = false
                this.showSounds = false
                this.pageEditor = false
                this.macroEditor = true
            },
            clearVerticalSlide: function () {
                this.showSounds = false;
                this.macroEditor = false;
                this.pageEditor = false;
            },
            openVerticalSlide: function (target) {
                this.clearVerticalSlide()
                this.showVerticalSlide = true
                this.showSounds = target === "sounds"
                this.macroEditor = target === "macros"
                this.pageEditor = target === "pages"
            },
            closeVerticalSlide: function () {
                this.showVerticalSlide = false
                this.clearVerticalSlide()
            },
            saveSequence: function (idx) {
                let action = this.config.Pages[this.currentPageId].Actions[idx];
                const sequence = action.Sequence;
                const caption = action.Caption;
                const icon = action.Icon || "ico-none";

                this.setSequenceSaveStatus(this.currentPageId, idx, "saving", "")
                saveSequence(this.currentPageId, idx, caption, icon, sequence)
            },
            formatSequence: function (idx) {
                let action = this.config.Pages[this.currentPageId].Actions[idx];

                this.setSequenceSaveStatus(this.currentPageId, idx, "formatting", "")
                formatSequence(this.currentPageId, idx, action.Sequence)
            },
            scheduleValidateSequence: function (idx) {
                let pageId = this.currentPageId
                let action = this.config.Pages[pageId].Actions[idx]
                let key = this.sequenceSaveKey(pageId, idx)

                clearTimeout(this.sequenceValidationTimers[key])
                this.setSequenceValidationStatus(pageId, idx, "validating", "")

                this.sequenceValidationTimers[key] = setTimeout(function () {
                    validateSequence(pageId, idx, action.Sequence)
                }, 300)
            },
            deleteSequence: function (idx) {
                deleteSequence(this.currentPageId, idx)
            },
            duplicateSequence: function (idx) {
                duplicateSequence(this.currentPageId, idx)
            },
            moveSequence: function (idx, offset) {
                moveSequence(this.currentPageId, idx, offset)
            },
            applyRecording: function (idx) {
                if (!this.recorder.Sequence) {
                    return
                }

                let action = this.config.Pages[this.currentPageId].Actions[idx]
                action.Sequence = this.recorder.Sequence
                this.scheduleValidateSequence(idx)
            },
            canMoveSequence: function (idx, offset) {
                let page = this.currentPage()
                if (!page || !page.Actions) {
                    return false
                }

                let target = idx + offset
                return target >= 0 && target < page.Actions.length
            },
            respondSaveSequence: function (pageId, sequenceIndex, sequence, meta, success, error) {
                if (sequenceIndex < 0 && this.remoteRecorder.Open) {
                    this.remoteRecorder.Saving = false
                    this.remoteRecorder.Error = error || "Could not save recording"
                    return
                }

                this.config.Pages[pageId].Actions[sequenceIndex].Sequence = sequence;
                this.config.Pages[pageId].Actions[sequenceIndex].Meta = meta;
                if (success) {
                    this.clearSequenceExecutionError(pageId, sequenceIndex)
                    this.setSequenceSaveStatus(pageId, sequenceIndex, "saved", "")
                    this.setSequenceValidationStatus(pageId, sequenceIndex, "valid", "")
                    return
                }

                this.setSequenceSaveStatus(pageId, sequenceIndex, "error", error || "error saving sequence")
                this.setSequenceValidationStatus(pageId, sequenceIndex, "error", error || "error saving sequence")
                console.error(error || "error saving sequence")
            },
            respondFormatSequence: function (pageId, sequenceIndex, sequence, meta, success, error) {
                if (success) {
                    this.config.Pages[pageId].Actions[sequenceIndex].Sequence = sequence;
                    this.config.Pages[pageId].Actions[sequenceIndex].Meta = meta;
                    this.clearSequenceExecutionError(pageId, sequenceIndex)
                    this.setSequenceSaveStatus(pageId, sequenceIndex, "formatted", "")
                    this.setSequenceValidationStatus(pageId, sequenceIndex, "valid", "")
                    return
                }

                this.setSequenceSaveStatus(pageId, sequenceIndex, "error", error || "error formatting sequence")
                this.setSequenceValidationStatus(pageId, sequenceIndex, "error", error || "error formatting sequence")
                console.error(error || "error formatting sequence")
            },
            respondValidateSequence: function (pageId, sequenceIndex, sequence, meta, success, error) {
                if (!this.config.Pages[pageId] ||
                    !this.config.Pages[pageId].Actions[sequenceIndex] ||
                    this.config.Pages[pageId].Actions[sequenceIndex].Sequence !== sequence) {
                    return
                }

                if (success) {
                    this.config.Pages[pageId].Actions[sequenceIndex].Meta = meta;
                    this.setSequenceValidationStatus(pageId, sequenceIndex, "valid", "")
                    return
                }

                this.config.Pages[pageId].Actions[sequenceIndex].Meta = meta;
                this.setSequenceValidationStatus(pageId, sequenceIndex, "error", error || "invalid sequence")
            },
            sequenceSaveKey: function (pageId, sequenceIndex) {
                return pageId + ":" + sequenceIndex
            },
            setSequenceSaveStatus: function (pageId, sequenceIndex, state, error) {
                Vue.set(this.sequenceSaveStatus, this.sequenceSaveKey(pageId, sequenceIndex), {
                    State: state,
                    Error: error || ""
                })
            },
            getSequenceSaveStatus: function (pageId, sequenceIndex) {
                let key = this.sequenceSaveKey(pageId, sequenceIndex)
                if (!this.sequenceSaveStatus[key]) {
                    return {State: "", Error: ""}
                }

                return this.sequenceSaveStatus[key]
            },
            sequenceSaveStatusClass: function (pageId, sequenceIndex) {
                let state = this.getSequenceSaveStatus(pageId, sequenceIndex).State
                return {
                    'sequence-save-status': true,
                    'is-saving': state === "saving",
                    'is-formatting': state === "formatting",
                    'is-formatted': state === "formatted",
                    'is-saved': state === "saved",
                    'is-error': state === "error"
                }
            },
            sequenceSaveStatusText: function (pageId, sequenceIndex) {
                let status = this.getSequenceSaveStatus(pageId, sequenceIndex)
                if (status.State === "saving") {
                    return "Saving..."
                }

                if (status.State === "formatting") {
                    return "Formatting..."
                }

                if (status.State === "formatted") {
                    return "Formatted - not saved yet"
                }

                if (status.State === "saved") {
                    return "Saved"
                }

                if (status.State === "error") {
                    return status.Error
                }

                return ""
            },
            setSequenceValidationStatus: function (pageId, sequenceIndex, state, error) {
                Vue.set(this.sequenceValidationStatus, this.sequenceSaveKey(pageId, sequenceIndex), {
                    State: state,
                    Error: error || ""
                })
            },
            getSequenceValidationStatus: function (pageId, sequenceIndex) {
                let key = this.sequenceSaveKey(pageId, sequenceIndex)
                if (!this.sequenceValidationStatus[key]) {
                    return {State: "", Error: ""}
                }

                return this.sequenceValidationStatus[key]
            },
            sequenceValidationClass: function (pageId, sequenceIndex) {
                let state = this.getSequenceValidationStatus(pageId, sequenceIndex).State
                return {
                    'sequence-validation-status': true,
                    'is-validating': state === "validating",
                    'is-valid': state === "valid",
                    'is-error': state === "error"
                }
            },
            sequenceValidationText: function (pageId, sequenceIndex) {
                let status = this.getSequenceValidationStatus(pageId, sequenceIndex)
                if (status.State === "validating") {
                    return "Validating..."
                }

                if (status.State === "valid") {
                    return "Script valid"
                }

                if (status.State === "error") {
                    return status.Error
                }

                return ""
            },
            toggleCommandHelp: function () {
                this.showCommandHelp = !this.showCommandHelp
            },
            updatePlaybackState: function (state) {
                let isRunning = state.State === "playing" || state.State === "looping" || state.State === "paused"
                let playbackKey = state.PageId + ":" + state.SequenceIndex + ":" + state.Caption
                if (isRunning && state.HasSequence && playbackKey !== this.lastPlaybackStartKey) {
                    this.lastPlaybackStartKey = playbackKey
                    this.appendPlaybackStart(state)
                }
                if (!isRunning || !state.HasSequence) {
                    this.lastPlaybackStartKey = ""
                    this.playbackTimeline.Visible = false
                }

                this.playbackState = state;
                if (state.State === "error" && state.Error && state.PageId && state.SequenceIndex >= 0) {
                    Vue.set(this.sequenceExecutionErrors, this.sequenceSaveKey(state.PageId, state.SequenceIndex), state.Error)
                    this.appendLogMessage("Execution error in " + state.Caption + ": " + state.Error)
                }
            },
            togglePause: function () {
                pause()
            },
            abort:function(){
                abort()
            },
            releaseInputs: function () {
                releaseInputs()
            }
        }
    })
}

function updateAppConfig(config) {
    app.updateConfig(config)
}

function updateConnectionStatus(isConnected) {
    app.updateConnectionStatus(isConnected)
}

function updateInputMode(state) {
    app.updateInputMode(state)
}

function updateRemoteInfo(info) {
    app.updateRemoteInfo(info)
}

function updateRecorderState(state) {
    app.updateRecorderState(state)
}

function updateSoundSettings(soundSettings) {
    app.updateSoundSettings(soundSettings)
}

function appendLogMessage(message) {
    app.appendLogMessage(message)
}

function appendPlaybackCommand(event) {
    app.appendPlaybackCommand(event)
}

function appendRecorderCommand(event) {
    app.appendRecorderCommand(event)
}

function respondSaveSequence(pageId, sequenceIndex, sequence, meta, success, error) {
    app.respondSaveSequence(pageId, sequenceIndex, sequence, meta, success, error)
}

function respondFormatSequence(pageId, sequenceIndex, sequence, meta, success, error) {
    app.respondFormatSequence(pageId, sequenceIndex, sequence, meta, success, error)
}

function respondValidateSequence(pageId, sequenceIndex, sequence, meta, success, error) {
    app.respondValidateSequence(pageId, sequenceIndex, sequence, meta, success, error)
}

function updatePlaybackState(state) {
    app.updatePlaybackState(state)
}

Vue.config.devtools = true

initApp();
connectToServer()
