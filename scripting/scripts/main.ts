import * as http from "http";
import * as mean from "mean";

console.log("main.ts: Initializing /foo handler...");
http.registerRequestHandler("/foo", http.RequestMethod.GET, (context) => {
    console.log("main.ts: Handling /foo");
    return context.response.respondWithString(http.ResponseCode.OK, "BAM!")
});

console.log("main.ts: Initializing /foo/bar handler...");
http.registerRequestHandler("/foo/bar", http.RequestMethod.GET, (context) => {
    console.log("main.ts: Handling /foo/bar");
    return context.response.respondWithError(http.ResponseCode.NotFound)
});

console.log("main.ts: Initializing /error handler...");
http.registerRequestHandler("/error", http.RequestMethod.GET, () => {
    console.log("main.ts: Gonna fail in /error");
    foo();
    return null;
});

const foo = () => {
    console.log("main.ts: We're almost there!");
    mean.fail(() => {
        console.log("main.ts: Failing...BOOM!");
        throw "This is serious boy!"
    })
};

mean.test(
    () => {
        return () => "test12345"
    }
);
