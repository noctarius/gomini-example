import * as http from "http";
import * as mean from "mean";

console.log("main.ts: Initializing /foo handler...");
http.registerRequestHandler("/foo", http.RequestMethod.GET, (context) => {
    return context.response.respondWithString(http.ResponseCode.OK, "BAM!")
});

http.registerRequestHandler("/foo/bar", http.RequestMethod.GET, (context) => {
    return context.response.respondWithError(http.ResponseCode.NotFound)
});

http.registerRequestHandler("/error", http.RequestMethod.GET, () => {
    foo();
});

const foo = () => {
    mean.fail(() => {
        throw "This is serious boy!"
    })
};