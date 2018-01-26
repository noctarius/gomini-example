package main

import (
	"github.com/labstack/echo"
	"github.com/relationsone/gomini"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/afero"
	"path/filepath"
	"log"
	"example"
	"os"
)

func main() {
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
		log.Fatal(err)
	}

	basePath := filepath.Join(base, "scripting", "scripts")
	typesPath := filepath.Join(base, "scripting", "@types")
	cachePath := filepath.Join(base, "target", "cache")

	os.MkdirAll(cachePath, os.ModePerm)

	kernelfs := buildKernelFilesystem(basePath, typesPath, cachePath)

	kernel, err := gomini.NewScriptKernel(afero.NewOsFs(), kernelfs, true)
	if err != nil {
		log.Fatal(err)
	}

	if err := kernel.LoadKernelModule(example.NewHttpKernelModule(e)); err != nil {
		log.Fatal(err)
	}

	if err := kernel.EntryPoint("main.ts"); err != nil {
		log.Fatal(err)
	}

	if err := kernel.Start(); err != nil {
		log.Fatal(err)
	}

	if err := e.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}

func buildKernelFilesystem(basePath, typesPath, cachePath string) afero.Fs {
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

	// Writable caching filesystem (for caching transpiled scripts)
	cachefs := afero.NewBasePathFs(osfs, cachePath)
	kernelfs.Mount(cachefs, "/kernel/cache")

	return kernelfs
}
