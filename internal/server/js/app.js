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
            showPage: false,
            showVerticalSlide: false,
            pageEditor: false,
            macroEditor: false,
            showSounds: false,
            currentPageId: null,
            connection: {isConnected: false},
            soundSettings: {Sessions: [], MainSession: {}},
            soundLoading: false,
            deletePageClicks: {},
            playbackState: {
                Caption: "",
                HasSequence: false,
                IsPlaying: false,
            },
            logMessages: []
        },
        methods: {
            currentPage: function () {
                return this.config.Pages[this.currentPageId];
            },
            newPage: function () {
                createNewPage();
            },
            removePage: function (pageId) {
                Vue.set(this.deletePageClicks, pageId, this.deletePageClicks[pageId] + 1)
                if (this.deletePageClicks[pageId] === 2) {
                    delete this.deletePageClicks[pageId]
                    deletePage(pageId);
                }
            },
            selectPage: function (pageKey) {
                this.currentPageId = pageKey
                this.showPage = true
            },
            newSequence: function () {
                createNewSequence(this.currentPageId);
            },
            startAction: function (actionIdx) {
                setConfigSequence(this.currentPageId, actionIdx)
            },
            updateConfig: function (config) {
                this.config = config
                this.hasConfig = true
            },
            updateConnectionStatus: function (isConnected) {
                this.connection.isConnected = isConnected
            },
            showPageSelection: function () {
                this.pageEditor = false;
                this.showPage = false
                this.currentPageId = null
            },
            showSoundSettings: function () {
                this.clearVerticalSlide()
                this.showVerticalSlide = !this.showVerticalSlide
                if (this.showVerticalSlide) {
                    this.showSounds = true
                    this.soundLoading = true
                    querySoundSettings()
                }
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
            showPagesEditor: function () {
                this.clearVerticalSlide()
                this.showVerticalSlide = !this.showVerticalSlide
                this.pageEditor = this.showVerticalSlide;
                if (!this.pageEditor) {
                    savePages(this.config.Pages)
                }
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
            showMacroEditor: function () {
                this.clearVerticalSlide()
                this.showVerticalSlide = !this.showVerticalSlide
                this.macroEditor = this.showVerticalSlide;
            },
            clearVerticalSlide: function () {
                this.showSounds = false;
                this.macroEditor = false;
            },
            saveSequence: function (idx) {
                let action = this.config.Pages[this.currentPageId].Actions[idx];
                const sequence = action.Sequence;
                const caption = action.Caption;

                saveSequence(this.currentPageId, idx, caption, sequence)
            },
            deleteSequence: function (idx) {
                deleteSequence(this.currentPageId, idx)
            },
            respondSaveSequence: function (pageId, sequenceIndex, sequence, success) {
                this.config.Pages[pageId].Actions[sequenceIndex].Sequence = sequence;
                if (!success) console.error("error saving sequence")
            },
            updatePlaybackState: function (state) {
                this.playbackState = state;
            },
            togglePause: function () {
                pause()
            },
            abort:function(){
                abort()
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

function updateSoundSettings(soundSettings) {
    app.updateSoundSettings(soundSettings)
}

function appendLogMessage(message) {
    app.appendLogMessage(message)
}

function respondSaveSequence(pageId, sequenceIndex, sequence, success) {
    app.respondSaveSequence(pageId, sequenceIndex, sequence, success)
}

function updatePlaybackState(state) {
    app.updatePlaybackState(state)
}

Vue.config.devtools = true

initApp();
connectToServer()