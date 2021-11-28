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
            soundSettings: {
                Sessions: [
                    {Name: "Chrome", Icon: "icon", Value: 20},
                    {Name: "Spotify", Icon: "icon", Value: 50},
                    {Name: "Diablo II", Icon: "icon", Value: 100},
                ]
            }
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

initApp();
connectToServer()