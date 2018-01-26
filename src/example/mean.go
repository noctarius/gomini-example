package example

import (
	"github.com/relationsone/gomini"
	"fmt"
)

type meanKernelModule struct {
}

func NewMeanKernelModule() gomini.KernelModuleDefinition {
	return &meanKernelModule{}
}

func (*meanKernelModule) ID() string {
	return "732fe3c9-4ca5-4558-8eaa-5055232308aa"
}

func (*meanKernelModule) Name() string {
	return "mean"
}

func (*meanKernelModule) ApiDefinitionFile() string {
	return "/kernel/@types/mean"
}

func (*meanKernelModule) SecurityInterceptor() gomini.SecurityInterceptor {
	return func(caller gomini.Bundle, property string) (accessGranted bool) {
		return true
	}
}

func (*meanKernelModule) ExtensionBinder() gomini.ExtensionBinder {
	return func(bundle gomini.Bundle, moduleBuilder gomini.ModuleBuilder) {
		moduleBuilder.DefineFunction("fail", func(callback func()) {
			fmt.Println("meanKernelModule: Just a quick go function call and going back into JS...")
			callback()
		}).EndModule()
	}
}
