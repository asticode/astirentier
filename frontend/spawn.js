const { spawn } = require('child_process')

class Spawn {

    init(logger) {
        this.l = logger
        return this
    }

    start() {
        this.s = spawn('./backend', ['-v'])
        this.s.stdout.on('data', (data) => {
            data.toString().split(/\r?\n/).forEach(function(s) {
                if (s === "") return
                try {
                    const j = JSON.parse(s)
                    const msg = j.msg + " source=" + j.source
                    switch (j.level) {
                        case "debug":
                            this.l.debug(msg)
                            break
                        case "error":
                            this.l.error(msg)
                            break
                        case "warn":
                            this.l.warn(msg)
                            break
                        default:
                            this.l.info(msg)
                            break
                    }
                } catch(e) {
                    this.l.error('JSON parsing spawn stdout "' + s + '" failed')
                }
            }.bind(this))
        })
        this.s.on('close', (code) => {
            this.l.info(`spawn exited with code ${code}`)
        })
    }

    stop() {
        this.s.kill('SIGTERM')
    }

}

module.exports = new Spawn()