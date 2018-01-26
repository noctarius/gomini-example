import * as http from "http";

http.registerRequestHandler("/foo", http.RequestMethod.GET, (context) => {
    return context.response.respondWithString(http.ResponseCode.OK, "BAM!")
});