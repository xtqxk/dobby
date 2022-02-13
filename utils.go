package dobby

import (
	"reflect"
	"strings"
)

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func getAllMethods(rt reflect.Type) (methods []reflect.Method) {
	var pStructs []reflect.Type
	for i := 0; i < rt.Elem().NumField(); i++ {
		p := rt.Elem().Field(i)
		if p.Type.Kind() == reflect.Struct {
			pStructs = append(pStructs, p.Type)
		}
	}

	pStructs = append(pStructs, rt)
	for _, fv := range pStructs {
		for i := 0; i < fv.NumMethod(); i++ {
			fm := fv.Method(i)
			methods = append(methods, fm)
		}
	}
	return
}
