import * as http from "http"
import {Foo} from "./testdir/test123";

console.log("Register /external from my external test app :-)");
http.registerRequestHandler("/external", http.RequestMethod.GET, (context) => {
        return context.response.respondWithString(http.ResponseCode.OK, new Foo().func())
    }
);
