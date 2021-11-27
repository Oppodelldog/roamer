const wsUrl = 'ws://127.0.0.1:10982/ws'
let ws = null
let wsMessageReader = new WebsocketMessageReader()
connectToServer()