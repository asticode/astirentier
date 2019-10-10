const { spawn } = require('child_process')

class Spawn {

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
                            logger.debug(msg)
                            break
                        case "error":
                            logger.error(msg)
                            break
                        case "warn":
                            logger.warn(msg)
                            break
                        default:
                            logger.info(msg)
                            break
                    }
                } catch(e) {
                    logger.error('JSON parsing spawn stdout "' + s + '" failed')
                }
            })
        })
        this.s.on('close', (code) => {
            logger.info(`spawn exited with code ${code}`)
        })
    }

    stop() {
        this.s.kill('SIGTERM')
    }

}

module.exports = {
    spawn: new Spawn(),
}