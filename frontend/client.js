const http = require('http')

let newPromise = function(logger, options) {
    return new Promise((resolve, reject) => {
        // Log
        const n = options.method + " request to " + options.path
        logger.debug("client.js: sending " + n)

        // Create request
        const r = http.request({
            host: "127.0.0.1",
            port: "6969",
            path: options.path,
            method: options.method,
            headers: options.headers
        }, function(resp) {
            // Get data
            let data = ''
            resp.on('data', (chunk) => {
                data += chunk
            })

            // Handle success
            resp.on('end', () => {
                // Parse JSON
                if (resp.headers["content-type"] === "application/json") {
                    data = JSON.parse(data)
                }

                // Check status code
                if (resp.statusCode < 200 || resp.statusCode > 299) {
                    reject("client.js: sending " + n + " failed with the invalid status code " + resp.statusCode + (typeof data.message !== "undefined" ? " and message " + data.message : ""))
                    return
                }
                resolve(data)
            })
        })

        // Handle error
        r.on("error", (err) => {
            reject("client.js: sending " + n + " failed with error " + err)
        })

        // Write body
        r.write(options.body)

        // Send
        r.end()
    })
}

class Client {

    init(logger) {
        this.l = logger
        return this
    }

    async send(options) {
        // Create result
        let r = {err: null}

        // Send
        try {
            // Create promise
            const p = newPromise(this.l, options)

            // Send
            r.data = await p
        }
        catch(err) {
            r.err = err
        }
        return r
    }

    async createBankAccount(label, bank) {
        return await this.send({
            path: "/bank-accounts",
            method: "POST",
            body: JSON.stringify({
                bank: bank,
                label: label,
            }),
        })
    }

    async openDB(path) {
        return await this.send({
            path: "/dbs",
            method: "POST",
            body: JSON.stringify({path: path}),
        })
    }

}

module.exports = new Client()