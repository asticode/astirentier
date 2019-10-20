const { BrowserWindow } = require('electron')

class OAuth2 {

    init(logger) {
        this.l = logger
        return this
    }

    start(url, done) {
        // Log
        this.l.debug("oauth2.js: starting oauth2 at " + url)

        // Reset data
        this.d = null

        // Create the window.
        this.w = new BrowserWindow({
            height: 600,
            webPreferences: {
                nodeIntegration: true
            },
            width: 600,
        })

        // Load
        this.w.loadURL(url)

        // Emitted when the window is closed.
        this.w.on('closed', () => {
            // Dereference
            this.w = null

            // Execute callback
            done(this.d)
        })
    }

    finish(data) {
        // Log
        this.l.debug("oauth2.js: finishing oauth2 with data " + JSON.stringify(data))

        // Store data
        this.d = data

        // Close window
        this.w.close()
    }

}

module.exports = new OAuth2()