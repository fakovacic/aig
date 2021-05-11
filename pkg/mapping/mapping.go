package mapping

import (
	"encoding/json"
	"reflect"
)

// Mapper map for elasticsearch
type Mapper struct {
	Mappings Mappings `json:"mappings"`
}

// JSON return json representation of mapper
func (m Mapper) JSON() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Mappings contain properties
type Mappings struct {
	Properties map[string]Field `json:"properties"`
}

// Field describe struct field
type Field struct {
	Type string `json:"type"`
}

// Map convert interface to mappings
func Map(in interface{}) Mapper {
	fields := make(map[string]Field)
	s := reflect.ValueOf(in)
	xType := reflect.TypeOf(
		s.Interface(),
	)
	fieldNum := xType.NumField()

	i := 0
	for {
		if i >= fieldNum {
			break
		}
		fieldMap := xType.Field(i).Tag.Get("mapping")
		if fieldMap == "-" {
			i++
			continue
		}

		fieldName := xType.Field(i).Tag.Get("json")
		if fieldName == "-" {
			i++
			continue
		}

		fieldType := ""
		switch xType.Field(i).Type.Kind() {
		case reflect.String:
			fieldType = "keyword"
		case reflect.Bool:
			fieldType = "boolean"
		case reflect.Int16, reflect.Uint16:
			fieldType = "short"
		case reflect.Int, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint32:
			fieldType = "integer"
		case reflect.Uint64:
			fieldType = "unsigned_long"
		case reflect.Float32:
			fieldType = "float"
		case reflect.Float64:
			fieldType = "double"
		}

		if fieldMap != "" {
			fieldType = fieldMap
		}

		if fieldType == "" {
			i++
			continue
		}

		fields[fieldName] = Field{
			Type: fieldType,
		}

		i++
	}

	return Mapper{
		Mappings: Mappings{
			Properties: fields,
		},
	}
}
