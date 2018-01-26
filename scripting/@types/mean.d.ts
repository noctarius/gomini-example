declare module "mean" {
    export interface Callback {
        ()
    }

    export function fail(callback:Callback)
}