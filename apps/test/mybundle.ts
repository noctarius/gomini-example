import * as http from "http"
import {Foo, Test} from './testdir/test123';
import * as test2 from "./testdir/test321";

export class Test2 {
    f() {
        new Foo();
    }
}

console.log("Register /external from my external test app :-)");
http.registerRequestHandler("/external", http.RequestMethod.GET, (context) => {
        console.stackTrace();
        new test2.test.Obje();
        new Test();
        return context.response.respondWithString(http.ResponseCode.OK, new Foo().func())
    }
);

http.registerRequestHandler("/external/2", http.RequestMethod.GET, test);

function test(context: http.RequestContext): http.Error {
    return testHandler(context)
}

function testHandler(context: http.RequestContext): http.Error {
    console.stackTrace();
    return context.response.respondWithString(http.ResponseCode.OK, new Foo().func())
}
