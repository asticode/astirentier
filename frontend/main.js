const { app, BrowserWindow } = require('electron')
const { logger } = require('./logger.js')
const { spawn } = require('./spawn.js')
const { client } = require('./client.js')

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win

// Make sure to clean up before quitting
app.on('before-quit', function() {
    // Stop backend
    spawn.stop()
})

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', function () {
    // Start backend
    spawn.start()

    // Create the browser window.
    win = new BrowserWindow({
        frame: false,
        height: 350,
        resizable: false,
        webPreferences: {
            nodeIntegration: true
        },
        width: 350,
    })

    // and load the index.html of the app.
    win.loadFile('app/index.html')

    // Emitted when the window is closed.
    win.on('closed', () => {
        // Dereference the window object, usually you would store windows
        // in an array if your app supports multi windows, this is the time
        // when you should delete the corresponding element.
        win = null
    })
})

// Share some variables with renderers
global.client = client
global.logger = logger