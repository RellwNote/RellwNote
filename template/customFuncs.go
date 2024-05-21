package template

import (
	"errors"
	"html/template"
	"math/rand"
	"strings"
	"time"
)

var CustomFuncMap = template.FuncMap{
	"Dict":         Dict,
	"Iif":          Iif,
	"RandomString": RandomString,
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
