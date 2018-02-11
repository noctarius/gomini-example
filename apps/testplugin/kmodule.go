package main

import (
	"os"
	"time"
	"io"
	"reflect"
)

// // No C code required.
//import "C"

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt

	Name() string
	Readdir(count int) ([]os.FileInfo, error)
	Readdirnames(n int) ([]string, error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	WriteString(s string) (ret int, err error)
}

type Fs interface {
	Create(name string) (File, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(name string) (os.FileInfo, error)
	Name() string
	Chmod(name string, mode os.FileMode) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
}

type Getter func() (value interface{})
type Setter func(value interface{})

type BundleStatus int

const (
	BundleStatusStopped     BundleStatus = iota
	BundleStatusStarted
	BundleStatusStarting
	BundleStatusStopping
	BundleStatusDownloading
	BundleStatusUpdating
	BundleStatusFailed
)

type Flag int

const (
	FLAG_NOT_SET Flag = iota
	FLAG_FALSE
	FLAG_TRUE
)

type JsValue interface {
	ToInteger() int64
	String() string
	ToFloat() float64
	ToNumber() JsValue
	ToBoolean() bool
	SameAs(JsValue) bool
	Equals(JsValue) bool
	StrictEquals(JsValue) bool
	Export() interface{}
	ExportType() reflect.Type
}

type JsObject interface {
	DefineAccessorProperty(name string, getter, setter JsValue, configurable, enumerable Flag) error
	DefineDataProperty(name string, value JsValue, writable, configurable, enumerable Flag) error
	Equals(other JsValue) bool
	Export() interface{}
	ExportType() reflect.Type
	Get(name string) JsValue
	Keys() []string
	MarshalJSON() ([]byte, error)
	SameAs(other JsValue) bool
	Set(name string, value interface{}) error
	StrictEquals(other JsValue) bool
	String() string
	ToBoolean() bool
	ToFloat() float64
	ToInteger() int64
	ToNumber() JsValue
}

type Bundle interface {
	ID() string
	Name() string
	Privileged() bool
	Privileges() []string
	SecurityInterceptor() SecurityInterceptor
	Export(value JsValue, target interface{}) error
	Status() BundleStatus
	Filesystem() Fs

	NewObject() JsObject
	NewException(err error) JsObject
	ToValue(value interface{}) JsValue
	Define(property string, value interface{})
	DefineProperty(object JsObject, property string, value interface{}, getter Getter, setter Setter)
	DefineConstant(object JsObject, constant string, value interface{})
	PropertyDescriptor(object JsObject, property string) (value interface{}, writable bool, getter Getter, setter Setter)
	FreezeObject(object JsObject)
}

type SecurityInterceptor func(caller Bundle, property string) (accessGranted bool)

type ModuleBuilder interface {
	DefineObject(objectName string, objectBinder ObjectBinder) ModuleBuilder
	DefineFunction(functionName string, function interface{}) ModuleBuilder
	DefineProperty(
		propertyName string,
		value interface{},
		getter func() interface{},
		setter func(value interface{})) ModuleBuilder
	DefineConstant(constantName string, value interface{}) ModuleBuilder
	EndModule()
}

type ObjectBuilder interface {
	DefineObject(objectName string, objectBinder ObjectBinder) ObjectBuilder
	DefineFunction(functionName string, function interface{}) ObjectBuilder
	DefineProperty(
		propertyName string,
		value interface{},
		getter func() interface{},
		setter func(value interface{})) ObjectBuilder
	DefineConstant(constantName string, value interface{}) ObjectBuilder
	EndObject()
}

type ExtensionBinder func(bundle Bundle, moduleBuilder ModuleBuilder)

type ObjectBinder func(objectBuilder ObjectBuilder)

type KernelModuleDefinition interface {
	ID() string
	Name() string
	ApiDefinitionFile() string
	SecurityInterceptor() SecurityInterceptor
	ExtensionBinder() ExtensionBinder
}

type moduleDefinition struct {
}

func (moduleDefinition) ID() string {
	return "23b40c8e-5625-4572-b6a2-0b6b47f5a2a1"
}

func (moduleDefinition) Name() string {
	return "test-kmodule"
}

func (moduleDefinition) ApiDefinitionFile() string {
	return "/kernel/@types/test"
}

func (moduleDefinition) SecurityInterceptor() SecurityInterceptor {
	return func(caller Bundle, property string) bool {
		return true
	}
}

func (moduleDefinition) ExtensionBinder() ExtensionBinder {
	return func(bundle Bundle, moduleBuilder ModuleBuilder) {
		moduleBuilder.
			DefineFunction("native__hello_world", func(name string) string {
			return "Hello " + name
		}).EndModule()
	}
}

//export
func KLoad() (KernelModuleDefinition, error) {
	return &moduleDefinition{}, nil
}
