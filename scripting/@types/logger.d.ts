declare module "logger" {

    export interface Logger {
        debug(msg: string, ...args: object[])

        info(msg: string, ...args: object[])

        warn(msg: string, ...args: object[])

        error(msg: string, ...args: object[])

        fatal(msg: string, ...args: object[])
    }

    export var log: Logger
}