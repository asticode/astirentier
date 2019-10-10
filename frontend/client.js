const addr = "127.0.0.1:6969"

class Client {

    createSession(path) {
        logger.info("creating session with path=" + path)
    }

}

module.exports = {
    client: new Client(),
}