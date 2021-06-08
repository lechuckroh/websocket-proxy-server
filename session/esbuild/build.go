package esbuild

import (
	"github.com/evanw/esbuild/pkg/api"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	CompileOutDir      = ".temp"
	CompileOutFilename = "out.js"
)

type Compiler interface {
	Compile(filename string) api.BuildResult
}

type CompilerImpl struct{}

func (b *CompilerImpl) Transform(filename string) api.TransformResult {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	buildResult := api.Transform(string(bytes), api.TransformOptions{
		Loader: api.LoaderTS,
	})
	return buildResult
}

func (b *CompilerImpl) Compile(filename string) (string, api.BuildResult) {
	if err := os.MkdirAll(CompileOutDir, 0700); err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = os.RemoveAll(CompileOutDir)
	}()

	outputFile := path.Join(CompileOutDir, CompileOutFilename)
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{filename},
		Outfile:     outputFile,
		Bundle:      true,
		Write:       true,
		LogLevel:    api.LogLevelInfo,
	})

	if bytes, err := ioutil.ReadFile(outputFile); err != nil {
		return "", result
	} else {
		return string(bytes), result
	}
}
