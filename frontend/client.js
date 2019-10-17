const http = require('http')

let newPromise = function(options) {
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
            // Check status code
            if (resp.statusCode < 200 || resp.statusCode > 299) {
                reject("client.js: sending " + n + " failed with the invalid status code " + resp.statusCode)
                return
            }

            // Get data
            let data = ''
            resp.on('data', (chunk) => {
                data += chunk
            })

            // Handle success
            resp.on('end', () => {
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

let send = async function(options) {
    // Create result
    let r = {err: null}

    // Send
    try {
        // Create promise
        const p = newPromise(options)

        // Send
        r.data = await p
    }
    catch(err) {
        r.err = err
    }
    return r
}

class Client {

    async openDB(path) {
        return await send({
            path: "/db/open",
            method: "POST",
            body: JSON.stringify({path: path}),
        })
    }

}

module.exports = {
    client: new Client(),
}