const {contextBridge, ipcRenderer} = require('electron')

contextBridge.exposeInMainWorld("ipcRenderer", {
    "send": ipcRenderer.send,
    "on": ipcRenderer.on
})