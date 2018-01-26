package example

import (
	"github.com/relationsone/gomini"
	"github.com/labstack/echo"
	"github.com/dop251/goja"
	"net/http"
	"strconv"
)

type requestHandler func(context *goja.Object) error

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

func NewHttpKernelModule(e *echo.Echo) gomini.KernelModuleDefinition {
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

func (hkm *httpKernelModule) ExtensionBinder() gomini.ExtensionBinder {
	return func(bundle gomini.Bundle, moduleBuilder gomini.ModuleBuilder) {
		moduleBuilder.
			DefineFunction("registerRequestHandler", hkm.jsRegisterRequestHandler(bundle)).
			DefineObject("RequestMethod", hkm.jsRequestMethod).
			DefineObject("ResponseCode", hkm.jsResponseCode).
			EndModule()
	}
}

func (hkm *httpKernelModule) jsRequestMethod(builder gomini.ObjectBuilder) {
	builder.
		DefineConstant("GET", REQUEST_METHOD_GET).
		DefineConstant("POST", REQUEST_METHOD_POST).
		EndObject()
}

func (hkm *httpKernelModule) jsResponseCode(builder gomini.ObjectBuilder) {
	builder.
		DefineConstant("OK", OK).
		DefineConstant("NotFound", NotFound).
		DefineConstant("InternalServerError", InternalServerError).
		EndObject()
}

func (hkm *httpKernelModule) jsRegisterRequestHandler(bundle gomini.Bundle) func(string, int, requestHandler) bool {
	return func(path string, requestMethod int, handler requestHandler) bool {
		handlerAdapter := func(context echo.Context) error {
			jsContext := hkm.jsRequestContext(bundle, context)

			method := hkm.findRequestMethodByName(context.Request().Method)
			if method.handling(requestMethod) {
				return handler(jsContext)
			}

			return context.NoContent(int(NotFound))
		}

		return hkm.e.Any(path, handlerAdapter) != nil
	}
}

func (hkm *httpKernelModule) jsRequestContext(bundle gomini.Bundle, context echo.Context) *goja.Object {
	jsContext := bundle.NewObject()

	jsContext.Set("pathParam", func(key string) string {
		return context.Param(key)
	})

	jsContext.Set("queryParam", func(key string) string {
		return context.QueryParam(key)
	})

	jsContext.Set("formParam", func(key string) string {
		return context.FormValue(key)
	})

	bundle.DefineConstant(jsContext, "request", hkm.jsRequest(bundle, context.Request()))
	bundle.DefineConstant(jsContext, "response", hkm.jsResponse(bundle, context))
	bundle.DefineConstant(jsContext, "path", context.Path())

	return jsContext
}

func (hkm *httpKernelModule) jsRequest(bundle gomini.Bundle, request *http.Request) *goja.Object {
	jsRequest := bundle.NewObject()

	jsRequest.Set("header", func(key string) string {
		return request.Header.Get(key)
	})

	bundle.DefineConstant(jsRequest, "method", hkm.findRequestMethodByName(request.Method))
	bundle.DefineConstant(jsRequest, "url", request.URL.String())
	bundle.DefineConstant(jsRequest, "protocol", request.Proto)
	bundle.DefineConstant(jsRequest, "contentLength", request.ContentLength)
	bundle.DefineConstant(jsRequest, "host", request.Host)

	return jsRequest
}

func (hkm *httpKernelModule) jsResponse(bundle gomini.Bundle, context echo.Context) *goja.Object {
	jsResponse := bundle.NewObject()

	jsResponse.Set("respondWithString", func(responseCode ResponseCode, content string) error {
		return context.String(int(responseCode), content)
	})

	return jsResponse
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