package template

import (
	"errors"
	"html/template"
)

var CustomFuncMap = template.FuncMap{
	"Dict": Dict,
	"Iif":  Iif,
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
