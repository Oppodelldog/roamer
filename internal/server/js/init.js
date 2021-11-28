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
            showSounds: false,
            currentPage: null,
            currentPageKey: null,
            connection: {isConnected: false},
            soundSettings: {Sessions: []},
            soundLoading: false,
        },
        methods: {
            selectPage: function (pageKey) {
                this.currentPage = this.config.Pages[pageKey];
                this.currentPageKey = pageKey;
                this.showPage = true;
            },
            startAction: function (actionIdx) {
                setConfigSequence(this.currentPageKey, actionIdx)
            },
            updateConfig: function (config) {
                this.config = config;
                this.hasConfig = true;
            },
            updateConnectionStatus: function (isConnected) {
                this.connection.isConnected = isConnected;
            },
            showPageSelection: function () {
                this.showPage = false;
                this.currentPage = null;
                this.currentPageKey = null;
            },
            showSoundSettings: function () {
                this.showSounds = !this.showSounds;
                if (this.showSounds) {
                    this.soundLoading = true;
                    querySoundSettings();
                }
            },
            changeSoundSessionValue: function (idx) {
                let id = this.soundSettings.Sessions[idx].Id;
                let volume = this.soundSettings.Sessions[idx].Value;
                setVolume(id, parseFloat(volume))
            },
            updateSoundSettings: function (soundSettings) {
                this.soundSettings = soundSettings;
                this.soundLoading = false;
            },
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

initApp();
connectToServer()