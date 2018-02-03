import * as http from "http"

console.log("Register /external from my external test app :-)");
http.registerRequestHandler("/external", http.RequestMethod.GET, (context) =>
    context.response.respondWithString(http.ResponseCode.OK, "Hello from the external world")
);