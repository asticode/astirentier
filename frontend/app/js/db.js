const { ipcRenderer, remote } = require('electron')
const os = require('os')
const { client, logger } = remote.getGlobal("all")

let db = {
    init: function() {
        // Init libs
        asticode.notifier.init()

        // Handle new
        db.handleNew()

        // Handle open
        db.handleOpen()

        // Handle open last
        db.handleOpenLast()

        // Close
        db.handleClose()
    },
    handleNew: function() {
        // Handle click
        document.getElementById("new").addEventListener("click", function() {
            // Get path
            const p = remote.dialog.showSaveDialogSync({
                defaultPath: os.homedir() + "/" + os.userInfo().username +  ".rdbx"
            })

            // No path
            if (typeof p === "undefined") return

            // Open
            db.open(p)
        })
    },
    handleOpen: function() {
        // Handle click
        document.getElementById("open").addEventListener("click", function() {
            // Get path
            const ps = remote.dialog.showOpenDialogSync({
                defaultPath: os.homedir(),
            })

            // No path
            if (typeof ps === "undefined" || ps.length === 0) return

            // Open
            db.open(ps[0])
        })
    },
    handleOpenLast: function() {
        // Get element
        const e = document.getElementById("open-last")

        // Disable
        if (localStorage.getItem('db_path') === null) {
            e.style.color = "#cccccc"
            e.style.cursor = "auto"
            return
        }

        // Handle click
        e.addEventListener("click", function() {
            // Open
            db.open(localStorage.getItem('db_path'))
        })
    },
    handleClose: function() {
        // Handle click
        document.getElementById("close").addEventListener("click", function() { remote.app.quit() })
    },
    open: async function(path) {
        // Open
        const {err} = await client.openDB(path)
        if (err !== null) {
            logger.error('index.js: opening db failed: ' + err)
            asticode.notifier.error(err)
            return
        }

        // Store db path
        localStorage.setItem('db_path', path)

        // Notify the main process
        ipcRenderer.send('db.opened')
    }
}