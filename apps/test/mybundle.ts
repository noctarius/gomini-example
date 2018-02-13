import * as http from "http"
import {Foo} from "./testdir/test123";

console.log("Register /external from my external test app :-)");
http.registerRequestHandler("/external", http.RequestMethod.GET, (context) => {
        console.stackTrace();
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
