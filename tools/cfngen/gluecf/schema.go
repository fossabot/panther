package gluecf

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Infers Glue table column types from Go types, recursively descends types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Functions to infer schema by reflection

type CustomMapping struct {
	From reflect.Type // type to map (result of reflect.TypeOf() )
	To   string       // glue type to emit
}

// Walk object, create columns using JSON Serde expected types, allow optional custom mappings
func InferJSONColumns(obj interface{}, customMappings ...CustomMapping) (cols []Column) {
	customMappingsTable := make(map[string]string)
	for _, customMapping := range customMappings {
		customMappingsTable[customMapping.From.String()] = customMapping.To
	}

	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	// dereference pointers
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	for i := 0; i < objType.NumField(); i++ {
		fieldName, jsonType, skip := inferStructFieldType(objType.Field(i), customMappingsTable)
		if skip {
			continue
		}
		cols = append(cols, Column{Name: fieldName, Type: jsonType})
	}

	return cols
}

func inferStructFieldType(sf reflect.StructField, customMappingsTable map[string]string) (fieldName, jsonType string, skip bool) {
	t := sf.Type

	// deference pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	isUnexported := sf.PkgPath != ""
	if sf.Anonymous {
		if isUnexported && t.Kind() != reflect.Struct { // I can't seem to find a way to exercise this block in my tests
			// Ignore embedded fields of unexported non-struct types.
			skip = true
			return
		}
		// Do not ignore embedded fields of unexported struct types
		// since they may have exported fields.
	} else if isUnexported {
		// Ignore unexported non-embedded fields.
		skip = true
		return
	}

	// use json tag name if present
	tag := sf.Tag.Get("json")
	if tag == "-" {
		skip = true
		return
	}
	fieldName, _ = parseTag(tag)
	if fieldName == "" {
		fieldName = sf.Name
	}

	switch t.Kind() { // NOTE: not all possible nestings have been implemented
	case reflect.Slice:

		sliceOfType := t.Elem()
		switch sliceOfType.Kind() {
		case reflect.Struct:
			jsonType = fmt.Sprintf("array<struct<%s>>", inferStruct(sliceOfType, customMappingsTable))
			return
		case reflect.Map:
			jsonType = fmt.Sprintf("array<%s>", inferMap(sliceOfType, customMappingsTable))
			return
		default:
			jsonType = fmt.Sprintf("array<%s>", toJSONType(sliceOfType))
			return
		}

	case reflect.Map:
		return fieldName, inferMap(t, customMappingsTable), skip

	case reflect.Struct:

		if to, found := customMappingsTable[t.String()]; found {
			jsonType = to
			return
		}

		jsonType = fmt.Sprintf("struct<%s>", inferStruct(t, customMappingsTable))
		return

	default:
		if mappedType, found := customMappingsTable[t.String()]; found {
			jsonType = mappedType
			return
		}

		// simple types
		jsonType = toJSONType(t)
		return
	}
}

// Recursively expand a struct
func inferStruct(structType reflect.Type, customMappingsTable map[string]string) string { // return comma delimited
	// recurse over components to get types
	numFields := structType.NumField()
	var keyPairs []string
	for i := 0; i < numFields; i++ {
		subFieldName, subFieldJSONType, subFieldSkip := inferStructFieldType(structType.Field(i), customMappingsTable)
		if subFieldSkip {
			continue
		}
		keyPairs = append(keyPairs, subFieldName+":"+subFieldJSONType)
	}
	return strings.Join(keyPairs, ",")
}

// Recursively expand a map
func inferMap(t reflect.Type, customMappingsTable map[string]string) (jsonType string) {
	mapOfType := t.Elem()
	if mapOfType.Kind() == reflect.Struct {
		jsonType = fmt.Sprintf("map<%s,struct<%s>>", t.Key(), inferStruct(mapOfType, customMappingsTable))
		return
	}
	jsonType = fmt.Sprintf("map<%s,%s>", t.Key(), toJSONType(mapOfType))
	return
}

// Primitive mappings
func toJSONType(t reflect.Type) (jsonType string) {
	switch t.String() {
	case "bool":
		jsonType = "boolean"
	case "string":
		jsonType = "string"
	case "int8":
		jsonType = "tinyint"
	case "int16":
		jsonType = "smallint"
	case "int":
		// int is problematic due to definition (at least 32bits ...)
		switch strconv.IntSize {
		case 32:
			jsonType = "int"
		case 64:
			jsonType = "bigint"
		default:
			panic(fmt.Sprintf("Size of native int unexpected: %d", strconv.IntSize))
		}
	case "int32":
		jsonType = "int"
	case "int64":
		jsonType = "bigint"
	case "float32":
		jsonType = "float"
	case "float64":
		jsonType = "double"
	case "interface {}":
		jsonType = "string" // best we can do in this case
	default:
		panic("Cannot map " + t.String())
	}

	return jsonType
}

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tag[idx+1:]
	}
	return tag, ""
}
