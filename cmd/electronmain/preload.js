const {contextBridge, ipcRenderer} = require('electron')

contextBridge.exposeInMainWorld("ipcRenderer", {
    "send": ipcRenderer.send,
    "on": (e,l) => {
        ipcRenderer.on(e,l)
    },
    "removeListener": ipcRenderer.removeListener,
})