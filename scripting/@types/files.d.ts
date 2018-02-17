declare module "files" {

    export enum FileMode {

    }

    export class File {
        constructor(path: string);
        constructor(path: string, mode )

        readonly fileName: string;
        readonly path: string;
    }
}