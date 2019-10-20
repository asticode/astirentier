const { ipcRenderer, remote } = require('electron')
const { client, logger, oauth2 } = remote.getGlobal("all")

let bank_accounts = {
    init: function() {
        // Init libs
        asticode.modaler.init()
        asticode.modaler.setWidth("30em")
        asticode.notifier.init()

        // Handle OAuth
        bank_accounts.handleOAuth()

        // Handle add
        bank_accounts.handleAdd()
    },
    handleOAuth: function() {
        // Listen to oauth2 finish
        ipcRenderer.on('oauth2.finish', (event, arg) => {
            logger.info("oauth2.finish: " + JSON.stringify(arg))
        })
    },
    handleAdd: function() {
        document.getElementById("add").addEventListener("click", function() {
            // Create form
            const f = asticode.modaler.newForm()
            f.addError()
            f.addField({
                label: "Label",
                name: "label",
                required: true,
                type: "text",
            })
            f.addField({
                label: "Bank",
                name: "bank",
                required: true,
                type: "select",
                values: {
                    "ing": "ING Direct",
                },
            })
            f.addField({
                className: "btn btn-success",
                label: "Add",
                success: async function(fs) {
                    // Create
                    const { data, err } = await client.createBankAccount(fs["label"], fs["bank"])
                    if (err !== null) {
                        logger.error('index.js: creating bank account failed: ' + err)
                        f.showError(err)
                        return
                    }

                    // OAuth2
                    if (typeof data.oauth2_start_url !== "undefined") {
                        oauth2.start(data.oauth2_start_url, function(data) {
                            // Switch on type
                            switch (data.type) {
                                case 'error':
                                    // Show error
                                    f.showError(data.message)
                                    break
                                case 'success':
                                    // Hide modal
                                    asticode.modaler.hide()

                                    // Notify
                                    asticode.notifier.success(data.message)

                                    // TODO Refresh account list
                            }
                        })
                    }
                },
                type: "submit",
            })

            // Show modal
            asticode.modaler.setContent(f)
            asticode.modaler.show()
            f.focus()
        })
    },
}