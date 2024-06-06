package template

import (
	"bytes"
	"errors"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"rellwnote/core/config"
	"strings"
	"time"
)

var CustomFuncMap = template.FuncMap{
	"Dict":            Dict,
	"Iif":             Iif,
	"RandomString":    RandomString,
	"JS":              JS,
	"JSCode":          JSCode,
	"CSS":             CSS,
	"CSSCode":         CSSCode,
	"URL":             URL,
	"DynamicTemplate": DynamicTemplate,
	"Add": func(a, b int) int {
		return a + b
	},
}

func Dict(args ...interface{}) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = args[i+1]
	}
	return dict, nil
}

func Iif(args ...interface{}) interface{} {
	if interfaceToBool(args[0]) {
		return args[1]
	}
	if len(args) > 2 {
		return args[2]
	}
	return ""
}

func RandomString(length int) []string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789                                \n"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return strings.Split(string(b), "\n")
}

func JS(path string) interface{} {
	res, err := os.ReadFile(filepath.Join(config.ProgramDir, config.TemplateDir, path))
	if err != nil {
		return err.Error()
	}
	return template.JS(res)
}

func JSCode(code string) interface{} {
	return template.JS(code)
}

func CSS(path string) interface{} {
	res, err := os.ReadFile(filepath.Join(config.ProgramDir, config.TemplateDir, path))
	if err != nil {
		return err.Error()
	}
	return template.CSS(res)
}

func CSSCode(code string) interface{} {
	return template.CSS(code)
}

func URL(url string) interface{} {
	return template.URL(url)
}

func DynamicTemplate(path string, args interface{}) interface{} {
	var buf bytes.Buffer

	err := LastLoadedTemplate.ExecuteTemplate(&buf, path, args)
	if err != nil {
		return err.Error()
	}
	return template.HTML(buf.String())
}

func interfaceToBool(i interface{}) bool {
	switch v := i.(type) {
	case bool:
		return v
	case string:
		return len(v) > 0
	case int:
	case uint:
	case float32:
	case float64:
		return v > 0
	}
	return false
}
