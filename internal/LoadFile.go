package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/govel-framework/lamb/ast"
	"github.com/govel-framework/lamb/lexer"
	"github.com/govel-framework/lamb/object"
	"github.com/govel-framework/lamb/parser"
)

type evalFunc func(ast.Node, *object.Environment) interface{}

// LoadFile parse the file received and writes the result in the io.Writer.
func LoadFile(fileName string, vars map[string]interface{}, out io.Writer, evaluator evalFunc, env object.Environment) error {
	// get the base directory from the env.
	baseDir := os.Getenv("GOVEL_LAMB_BASE_DIR")

	// replace every '.' in the file path with '/' and append '.lamb.html' at the end.
	file := strings.ReplaceAll(fileName, ".", "/") + ".lamb.html"
	file = baseDir + file

	// add the vars
	for key, value := range vars {
		env.Set(key, value)
	}

	// check the cache
	var cache string

	if cacheValue, exists := vars["__cache"]; exists && fmt.Sprintf("%T", cacheValue) == "string" {
		cache = fmt.Sprintf("%s", cacheValue)
	}

	cacheDir := os.Getenv("GOVEL_LAMB_CACHE_DIR")

	cacheFile := cacheDir + "/" + fileName

	// check if the file exists
	if stat, statError := os.Stat(cacheFile); statError == nil && cache != "" {
		cacheTime, _ := time.ParseDuration(os.Getenv("GOVEL_LAMB_CACHE_TIME"))

		// check if the file is older than the cache time
		modTime := stat.ModTime()

		if time.Since(modTime) > cacheTime {
			// delete the file
			os.Remove(cacheFile)

		} else {
			// read the file
			content, err := os.ReadFile(cacheFile)

			if err != nil {
				return err
			}

			out.Write(content)

			return nil
		}

	}

	// set the file name
	env.FileName = file

	content, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	l := lexer.New(string(content))

	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {

		for _, e := range p.Errors() {
			return fmt.Errorf("%s: %s\n", file, e)
		}
	}

	evaluated := evaluator(program, &env)

	if evaluated != nil {

		if _, isError := evaluated.(error); isError {
			return errors.New(fmt.Sprintf("%s", evaluated))
		}

		out.Write([]byte(fmt.Sprintf("%s", evaluated)))

		go func() {
			// check if the cache is enabled
			if cache != "" {
				switch cache {
				case "all":
					// create the cache directory
					if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
						os.Mkdir(cacheDir, os.ModePerm)
					}

					// write the file
					err = os.WriteFile(cacheFile, []byte(fmt.Sprintf("%s", evaluated)), 0644)

					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}

	return nil
}
