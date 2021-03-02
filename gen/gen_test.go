package gen

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestGen_Build(t *testing.T) {
	t.Parallel()

	config := Config{
		SearchDir:   "../testdata/simple",
		MainAPIFile: "./main.go",
		OutputDir:   "../testdata/simple/docs",
	}

	assert.NoError(t, New().Build(&config))

	expectedFiles := []string{
		filepath.Join(config.OutputDir, "docs.go"),
		filepath.Join(config.OutputDir, "swagger.json"),
		filepath.Join(config.OutputDir, "swagger.yaml"),
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			assert.NoError(t, err)
		}
		assert.NoError(t, os.Remove(expectedFile))
	}
}

func TestGen_BuildSnakecase(t *testing.T) {
	t.Parallel()

	config := Config{
		SearchDir:          "../testdata/simple2",
		MainAPIFile:        "./main.go",
		OutputDir:          "../testdata/simple2/docs",
		PropNamingStrategy: "snakecase",
	}

	assert.NoError(t, New().Build(&config))

	expectedFiles := []string{
		filepath.Join(config.OutputDir, "docs.go"),
		filepath.Join(config.OutputDir, "swagger.json"),
		filepath.Join(config.OutputDir, "swagger.yaml"),
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			assert.NoError(t, err)
		}
		assert.NoError(t, os.Remove(expectedFile))
	}
}

func TestGen_BuildLowerCamelcase(t *testing.T) {
	t.Parallel()

	config := Config{
		SearchDir:   "../testdata/simple3",
		MainAPIFile: "./main.go",
		OutputDir:   "../testdata/simple3/docs",
	}

	assert.NoError(t, New().Build(&config))

	expectedFiles := []string{
		filepath.Join(config.OutputDir, "docs.go"),
		filepath.Join(config.OutputDir, "swagger.json"),
		filepath.Join(config.OutputDir, "swagger.yaml"),
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			assert.NoError(t, err)
		}
		assert.NoError(t, os.Remove(expectedFile))
	}
}

func TestGen_jsonIndent(t *testing.T) {
	t.Parallel()

	config := Config{
		SearchDir:   "../testdata/simple",
		MainAPIFile: "./main.go",
		OutputDir:   "../testdata/simple/docs",
	}

	gen := New()
	gen.jsonIndent = func(data interface{}) ([]byte, error) {
		return nil, errors.New("fail")
	}
	assert.Error(t, gen.Build(&config))
}

func TestGen_jsonToYAML(t *testing.T) {
	config := Config{
		SearchDir:   "../testdata/simple",
		MainAPIFile: "./main.go",
		OutputDir:   "../testdata/simple/docs",
	}

	gen := New()
	gen.jsonToYAML = func(data []byte) ([]byte, error) {
		return nil, errors.New("fail")
	}
	assert.Error(t, gen.Build(&config))

	expectedFiles := []string{
		filepath.Join(config.OutputDir, "docs.go"),
		filepath.Join(config.OutputDir, "swagger.json"),
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			assert.Error(t, err)
		}
		assert.NoError(t, os.Remove(expectedFile))
	}
}

func TestGen_SearchDirIsNotExist(t *testing.T) {
	t.Parallel()

	config := Config{
		SearchDir:   "../isNotExistDir",
		MainAPIFile: "./main.go",
	}

	assert.EqualError(t, New().Build(&config), "dir: ../isNotExistDir is not exist")
}

func TestGen_MainAPiNotExist(t *testing.T) {
	config := Config{
		SearchDir:   "../testdata/simple",
		MainAPIFile: "./notexists.go",
	}

	assert.Error(t, New().Build(&config))
}

func TestGen_OutputIsNotExist(t *testing.T) {
	config := Config{
		SearchDir:   "../testdata/simple",
		MainAPIFile: "./main.go",
		OutputDir:   "/dev/null",
	}

	assert.Error(t, New().Build(&config))
}

func TestGen_FailToWrite(t *testing.T) {
	outputDir := filepath.Join(os.TempDir(), "swagg", "test")

	var propNamingStrategy string
	config := Config{
		SearchDir:          "../testdata/simple",
		MainAPIFile:        "./main.go",
		OutputDir:          outputDir,
		PropNamingStrategy: propNamingStrategy,
	}

	assert.NoError(t, os.MkdirAll(outputDir, 0755))

	assert.NoError(t, os.RemoveAll(filepath.Join(outputDir, "swagger.yaml")))
	assert.NoError(t, os.Mkdir(filepath.Join(outputDir, "swagger.yaml"), 0755))
	assert.Error(t, New().Build(&config))

	assert.NoError(t, os.RemoveAll(filepath.Join(outputDir, "swagger.json")))
	assert.NoError(t, os.Mkdir(filepath.Join(outputDir, "swagger.json"), 0755))
	assert.Error(t, New().Build(&config))

	assert.NoError(t, os.RemoveAll(filepath.Join(outputDir, "docs.go")))

	assert.NoError(t, os.Mkdir(filepath.Join(outputDir, "docs.go"), 0755))
	assert.Error(t, New().Build(&config))

	assert.NoError(t, os.RemoveAll(outputDir))
}

func TestGen_configWithOutputDir(t *testing.T) {
	config := Config{
		SearchDir:          "../testdata/simple",
		MainAPIFile:        "./main.go",
		OutputDir:          "../testdata/simple/docs",
		PropNamingStrategy: "",
	}

	assert.NoError(t, New().Build(&config))

	expectedFiles := []string{
		filepath.Join(config.OutputDir, "docs.go"),
		filepath.Join(config.OutputDir, "swagger.json"),
		filepath.Join(config.OutputDir, "swagger.yaml"),
	}
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			assert.NoError(t, err)
		}
		assert.NoError(t, os.Remove(expectedFile))
	}
}

func TestGen_formatSource(t *testing.T) {
	t.Parallel()

	src := `package main

import "net

func main() {}
`
	g := New()

	res := g.formatSource([]byte(src))
	assert.Equal(t, []byte(src), res, "Should return same content due to fmt fail")

	src2 := `package main

import "fmt"

func main() {
fmt.Print("Helo world")
}
`
	res = g.formatSource([]byte(src2))
	assert.NotEqual(t, []byte(src2), res, "Should return fmt code")
}

type mockWriter struct {
	hook func([]byte)
}

func (w *mockWriter) Write(data []byte) (int, error) {
	if w.hook != nil {
		w.hook(data)
	}

	return len(data), nil
}

func TestGen_writeGoDoc(t *testing.T) {
	t.Parallel()

	gen := New()

	swapTemplate := packageTemplate

	var config Config

	packageTemplate = `{{{`
	err := gen.writeGoDoc("docs", nil, nil, &config)
	assert.Error(t, err)

	packageTemplate = `{{.Data}}`
	swagger := &spec.Swagger{
		VendorExtensible: spec.VendorExtensible{},
		SwaggerProps: spec.SwaggerProps{
			Info: &spec.Info{},
		},
	}
	err = gen.writeGoDoc("docs", &mockWriter{}, swagger, &config)
	assert.Error(t, err)

	config.GeneratedTime = true
	packageTemplate = `{{ if .GeneratedTime }}Fake Time{{ end }}`
	err = gen.writeGoDoc("docs",
		&mockWriter{
			hook: func(data []byte) {
				assert.Equal(t, "Fake Time", string(data))
			},
		}, swagger, &config)
	assert.NoError(t, err)

	config.GeneratedTime = false
	err = gen.writeGoDoc("docs",
		&mockWriter{
			hook: func(data []byte) {
				assert.Equal(t, "", string(data))
			},
		}, swagger, &config)
	assert.NoError(t, err)

	packageTemplate = swapTemplate
}

func TestGen_GeneratedDoc(t *testing.T) {
	config := Config{
		SearchDir:          "../testdata/simple",
		MainAPIFile:        "./main.go",
		OutputDir:          "../testdata/simple/docs",
		PropNamingStrategy: "",
	}

	assert.NoError(t, New().Build(&config))
	gocmd, err := exec.LookPath("go")
	assert.NoError(t, err)

	cmd := exec.Command(gocmd, "build", filepath.Join(config.OutputDir, "docs.go"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	assert.NoError(t, cmd.Run())

	expectedFiles := []string{
		filepath.Join(config.OutputDir, "docs.go"),
		filepath.Join(config.OutputDir, "swagger.json"),
		filepath.Join(config.OutputDir, "swagger.yaml"),
	}
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Fatal(err)
		}
		assert.NoError(t, os.Remove(expectedFile))
	}
}

func TestGen_cgoImports(t *testing.T) {
	t.Parallel()

	config := Config{
		SearchDir:          "../testdata/simple_cgo",
		MainAPIFile:        "./main.go",
		OutputDir:          "../testdata/simple_cgo/docs",
		PropNamingStrategy: "",
		ParseDependency:    true,
	}

	assert.NoError(t, New().Build(&config))

	expectedFiles := []string{
		filepath.Join(config.OutputDir, "docs.go"),
		filepath.Join(config.OutputDir, "swagger.json"),
		filepath.Join(config.OutputDir, "swagger.yaml"),
	}
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			assert.NoError(t, err)
		}
		assert.NoError(t, os.Remove(expectedFile))
	}
}
