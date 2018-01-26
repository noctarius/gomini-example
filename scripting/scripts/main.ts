import * as http from "http";
import * as mean from "mean";

console.log("main.ts: Initializing /foo handler...");
http.registerRequestHandler("/foo", http.RequestMethod.GET, (context) => {
    return context.response.respondWithString(http.ResponseCode.OK, "BAM!")
});

console.log("main.ts: Initializing /foo/bar handler...");
http.registerRequestHandler("/foo/bar", http.RequestMethod.GET, (context) => {
    return context.response.respondWithError(http.ResponseCode.NotFound)
});

console.log("main.ts: Initializing /error handler...");
http.registerRequestHandler("/error", http.RequestMethod.GET, () => {
    foo();
    return null;
});

const foo = () => {
    mean.fail(() => {
        throw "This is serious boy!"
    })
};