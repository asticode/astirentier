let write = function(startedAt, severity, msg) {
    process.stdout.write(severity + "[" + ((Date.now()-startedAt)/1000).toFixed(0) + "] " + msg + "\n")
}

class Logger {

    constructor() {
        this.startedAt = Date.now()
    }

    debug(msg) {
        write(this.startedAt, "DEBU", msg)
    }

    error(msg) {
        write(this.startedAt, "ERRO", msg)
    }

    info(msg) {
        write(this.startedAt, "INFO", msg)
    }

    warn(msg) {
        write(this.startedAt, "WARN", msg)
    }

}

module.exports = new Logger()