package example

import (
	"github.com/relationsone/gomini"
	"github.com/labstack/echo"
	"github.com/dop251/goja"
	"net/http"
	"strconv"
	"github.com/go-errors/errors"
)

type requestHandler func(context gomini.Object) error

type RequestMethod int

const (
	REQUEST_METHOD_GET  RequestMethod = 1
	REQUEST_METHOD_POST RequestMethod = 2
)

func (rm RequestMethod) handling(requestMethod int) bool {
	method := int(rm)
	return requestMethod&method == method
}

type ResponseCode int

const (
	OK                  ResponseCode = 200
	NotFound            ResponseCode = 404
	InternalServerError ResponseCode = 500
)

type httpKernelModule struct {
	e *echo.Echo
}

func NewHttpKernelModule(e *echo.Echo) gomini.KernelModule {
	return &httpKernelModule{e}
}

func (hkm *httpKernelModule) ID() string {
	return "4b0f1733-d7e8-4a95-a671-6ab0c1600373"
}

func (hkm *httpKernelModule) Name() string {
	return "http"
}

func (hkm *httpKernelModule) ApiDefinitionFile() string {
	return "/kernel/@types/http"
}

func (hkm *httpKernelModule) SecurityInterceptor() gomini.SecurityInterceptor {
	return func(caller gomini.Bundle, property string) (accessGranted bool) {
		// UBER-HACKER!!! :-)
		return true
	}
}

func (hkm *httpKernelModule) KernelModuleBinder() gomini.KernelModuleBinder {
	return func(bundle gomini.Bundle, builder gomini.ObjectBuilder) {
		builder.
			DefineGoFunction("registerRequestHandler", "registerRequestHandler", hkm.jsRegisterRequestHandler(bundle)).
			DefineObjectProperty("RequestMethod", hkm.jsRequestMethod).
			DefineObjectProperty("ResponseCode", hkm.jsResponseCode)
	}
}

func (hkm *httpKernelModule) jsRequestMethod(builder gomini.ObjectBuilder) {
	builder.
		DefineConstant("GET", REQUEST_METHOD_GET).
		DefineConstant("POST", REQUEST_METHOD_POST)
}

func (hkm *httpKernelModule) jsResponseCode(builder gomini.ObjectBuilder) {
	builder.
		DefineConstant("OK", OK).
		DefineConstant("NotFound", NotFound).
		DefineConstant("InternalServerError", InternalServerError)
}

func (hkm *httpKernelModule) jsRegisterRequestHandler(bundle gomini.Bundle) func(string, int, requestHandler) bool {
	return func(path string, requestMethod int, handler requestHandler) bool {
		handlerAdapter := func(context echo.Context) error {
			jsContext := hkm.jsRequestContext(bundle, context)

			method := hkm.findRequestMethodByName(context.Request().Method)
			if method.handling(requestMethod) {
				defer func() {
					if x := recover(); x != nil {
						switch t := x.(type) {
						case *goja.Exception:
							context.String(http.StatusInternalServerError, t.String())
						case error:
							context.String(http.StatusInternalServerError, errors.New(t).ErrorStack())
						default:
							context.NoContent(http.StatusInternalServerError)
						}
					}
				}()
				return handler(jsContext)
			}

			return context.NoContent(int(NotFound))
		}

		return hkm.e.Any(path, handlerAdapter) != nil
	}
}

func (hkm *httpKernelModule) jsRequestContext(bundle gomini.Bundle, context echo.Context) gomini.Object {
	jsContext := bundle.NewObjectBuilder("requestContext")

	jsContext.DefineGoFunction("pathParam", "pathParam", func(key string) string {
		return context.Param(key)
	})

	jsContext.DefineGoFunction("queryParam", "queryParam", func(key string) string {
		return context.QueryParam(key)
	})

	jsContext.DefineGoFunction("formParam", "formParam", func(key string) string {
		return context.FormValue(key)
	})

	jsContext.DefineConstant("request", hkm.jsRequest(bundle, context.Request()))
	jsContext.DefineConstant("response", hkm.jsResponse(bundle, context))
	jsContext.DefineConstant("path", context.Path())

	return jsContext.Build()
}

func (hkm *httpKernelModule) jsRequest(bundle gomini.Bundle, request *http.Request) gomini.Object {
	jsRequest := bundle.NewObjectBuilder("request")

	jsRequest.DefineGoFunction("header", "header", func(key string) string {
		return request.Header.Get(key)
	})

	jsRequest.DefineConstant("method", hkm.findRequestMethodByName(request.Method))
	jsRequest.DefineConstant("url", request.URL.String())
	jsRequest.DefineConstant("protocol", request.Proto)
	jsRequest.DefineConstant("contentLength", request.ContentLength)
	jsRequest.DefineConstant("host", request.Host)

	return jsRequest.Build()
}

func (hkm *httpKernelModule) jsResponse(bundle gomini.Bundle, context echo.Context) gomini.Object {
	jsResponse := bundle.NewObjectBuilder("response")

	jsResponse.DefineGoFunction("respondWithString", "respondWithString", func(responseCode ResponseCode, content string) error {
		return context.String(int(responseCode), content)
	})

	jsResponse.DefineGoFunction("respondWithError", "respondWithError", func(code ResponseCode) error {
		return context.NoContent(int(code))
	})

	return jsResponse.Build()
}

func (hkm *httpKernelModule) findRequestMethodByName(method string) RequestMethod {
	switch method {
	case "GET":
		return REQUEST_METHOD_GET
	case "POST":
		return REQUEST_METHOD_POST
	}
	panic("Undefined request method: " + method)
}

func (hkm *httpKernelModule) findRequestMethodById(method int) RequestMethod {
	switch method {
	case int(REQUEST_METHOD_GET):
		return REQUEST_METHOD_GET
	case int(REQUEST_METHOD_POST):
		return REQUEST_METHOD_POST
	}
	panic("Undefined request method: " + strconv.Itoa(method))
}
