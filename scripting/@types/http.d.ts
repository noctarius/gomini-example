declare module "http" {
    export enum RequestMethod {
        GET  = 1,
        POST = 2
    }

    export enum ResponseCode {
        OK                  = 200,
        NotFound            = 404,
        InternalServerError = 500
    }

    export interface Error {
        msg: string
    }

    export interface Request {
        header(key: string): string

        readonly method: RequestMethod
        readonly url: string
        readonly protocol: string
        readonly contentLength: number
        readonly host: string
    }

    export interface Response {
        respondWithString(responseCode: ResponseCode, content: string): Error

        respondWithError(responseCode: ResponseCode): Error
    }

    export interface RequestContext {
        pathParam(key: string): string

        queryParam(key: string): string

        formPram(key: string): string

        readonly request: Request
        readonly response: Response
        readonly path: string
    }

    export interface RequestHandler {
        (context: RequestContext): Error
    }

    export function registerRequestHandler(path: string, requestMethods: number, handler: RequestHandler): boolean
}
