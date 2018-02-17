package main

import (
	"github.com/labstack/echo"
	"github.com/relationsone/gomini"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/afero"
	"path/filepath"
	"example"
	"os"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/relationsone/gomini/kmodules"
)

func main() {
	log.SetHandler(text.Default)
	log.SetLevel(log.DebugLevel)

	// Create new http server
	e := echo.New()

	// Activate logging and exception handler
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	// Build basic filesystem
	base, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	basePath := filepath.Join(base, "scripting", "scripts")
	typesPath := filepath.Join(base, "scripting", "@types")
	cachePath := filepath.Join(base, "target", "cache")
	appsPath := filepath.Join(base, "apps")

	os.MkdirAll(cachePath, os.ModePerm)

	kernelfs := buildKernelFilesystem(basePath, typesPath, appsPath, cachePath)

	kernel, err := gomini.NewScriptKernel(afero.NewOsFs(), kernelfs, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := kernel.LoadKernelModule(example.NewHttpKernelModule(e)); err != nil {
		log.Fatal(err.Error())
	}

	if err := kernel.LoadKernelModule(example.NewMeanKernelModule()); err != nil {
		log.Fatal(err.Error())
	}

	if err := kernel.LoadKernelModule(kmodules.NewLoggerModule()); err != nil {
		log.Fatal(err.Error())
	}

	if err := kernel.EntryPoint("/main.ts"); err != nil {
		log.Fatal(err.Error())
	}

	if err := kernel.Start(); err != nil {
		log.Fatal(err.Error())
	}

	e.HideBanner = true
	e.HidePort = true
	log.Info("Main: Server started at [::]:8000")
	if err := e.Start(":8000"); err != nil {
		log.Fatal(err.Error())
	}
}

func buildKernelFilesystem(basePath, typesPath, appsPath, cachePath string) afero.Fs {
	// Base filesystem, delegating to real filesystem
	osfs := afero.NewOsFs()

	// Prevent modules from mutating the real filesystem
	rofs := afero.NewReadOnlyFs(osfs)

	// The script kernel base directory
	rootfs := afero.NewBasePathFs(rofs, basePath)

	// Virtual root filesystem -> /
	kernelfs := gomini.NewCompositeFs(rootfs)

	// Typescript definitions fs
	typesfs := afero.NewBasePathFs(rofs, typesPath)

	// Mount the types filesystem into the root fs
	kernelfs.Mount(typesfs, "/kernel/@types")

	// External apps fs
	appsfs := afero.NewBasePathFs(rofs, appsPath)
	kernelfs.Mount(appsfs, "/kernel/apps")

	// Writable caching filesystem (for caching transpiled scripts)
	cachefs := afero.NewBasePathFs(osfs, cachePath)
	kernelfs.Mount(cachefs, "/kernel/cache")

	return kernelfs
}
