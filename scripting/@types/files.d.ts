declare module "files" {

    export enum FileMode {
        ReadOnly    = 1,
        WriteOnly   = 2,
        ReadWrite   = 4,
        Append      = 8,
        Create      = 16,
        Truncate    = 32,
        Synchronous = 64
    }

    export enum FileType {
        Unknown   = 1,
        Kernel    = 2,
        Directory = 4,
        File      = 8
    }

    export interface ReadBuffer {
        readUint8(): number;
        readInt8(): number;
        readUint16(): number;
        readInt16(): number;
        readUint32(): number;
        readInt32(): number;
        readUint64(): number;
        readInt64(): number;
        readFloat32(): number;
        readFloat64(): number;
        readString(): number;
        readBoolean(): boolean;
        readAny<T>(): T;
    }

    export interface WriteBuffer {
        writeUint8(val: number): void;
        writeInt8(val: number): void;
        writeUint16(val: number): void;
        writeInt16(val: number): void;
        writeUint32(val: number): void;
        writeInt32(val: number): void;
        writeUint64(val: number): void;
        writeInt64(val: number): void;
        writeFloat32(val: number): void;
        writeFloat64(val: number): void;
        writeString(val: string): void;
        writeBoolean(val: boolean): void;
        writeAny(val: any): void;
    }

    export function resolvePath(path: string): Path;

    export interface Path {
        readonly name: string;
        readonly path: string;
        readonly type: FileType;

        exists(): boolean;
        mkdir(createParents: boolean);
        resolve(subpath: string): Path;

        toFile(...modes: FileMode[]): File;
        toPipe(): Pipe;
    }

    export interface BaseFile {
        readonly name: string;
        readonly path: string;
        readonly length: number;
        readonly type: FileType;

        close(): void;
    }

    export interface File extends BaseFile {
        readBuffer(): ReadBuffer;
        writeBuffer(): WriteBuffer;
    }

    export interface Pipe extends BaseFile {
        write(val: any): void;
        read<T>(): T;
    }
}