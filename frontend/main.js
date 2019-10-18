const { app, BrowserWindow, ipcMain } = require('electron')
const { logger } = require('./logger.js')
const { spawn } = require('./spawn.js')
const { client } = require('./client.js')

// Make sure to clean up before quitting
app.on('before-quit', function() {
    // Stop backend
    spawn.stop()
})

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', function () {
    // Windows
    let dbWindow, mainWindow

    // Start backend
    spawn.start()

    // Listen to renderers
    ipcMain.on('db.opened', () => {
        // Create the main window.
        mainWindow = new BrowserWindow({
            height: 600,
            webPreferences: {
                nodeIntegration: true
            },
            width: 600,
        })

        // Load
        mainWindow.loadFile('app/index.html')

        // Emitted when the window is closed.
        mainWindow.on('closed', () => {
            // Dereference the window object, usually you would store windows
            // in an array if your app supports multi windows, this is the time
            // when you should delete the corresponding element.
            mainWindow = null
        })

        // Close db window
        dbWindow.close()
    })

    // Create the db window.
    dbWindow = new BrowserWindow({
        frame: false,
        height: 350,
        resizable: false,
        webPreferences: {
            nodeIntegration: true
        },
        width: 350,
    })

    // Load
    dbWindow.loadFile('app/db.html')

    // Emitted when the window is closed.
    dbWindow.on('closed', () => {
        // Dereference the window object, usually you would store windows
        // in an array if your app supports multi windows, this is the time
        // when you should delete the corresponding element.
        dbWindow = null
    })
})

// Share some variables with renderers
global.client = client
global.logger = logger