declare module "mean" {
    export interface Callback {
        ()
    }

    export function fail(callback: Callback)

    export function test(f: () => () => string)
}