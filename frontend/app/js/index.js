const { remote } = require('electron')

let index = {
    init: function() {
        // New
        document.getElementById("new").addEventListener("click", function() {
            // Get path
            const p = remote.dialog.showSaveDialogSync({
                defaultPath: require('os').homedir() + "/default.rnt"
            })

            // No path
            if (typeof p === "undefined") return

            // Create session
            remote.getGlobal("client").createSession("my path")
        })

        // Close
        document.getElementById("close").addEventListener("click", function() { remote.app.quit() })
    },
}