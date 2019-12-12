package utils

import "reflect"

// ParseTagSlice takes a list of structs representing tags, and returns a map of key/value pairs
func ParseTagSlice(slice interface{}) map[string]*string {
	typedSlice := reflect.ValueOf(slice)

	tags := make(map[string]*string, typedSlice.Len())
	for i := 0; i < typedSlice.Len(); i++ {
		tagStruct := reflect.Indirect(typedSlice.Index(i))
		key := tagStruct.FieldByName("Key")
		value := tagStruct.FieldByName("Value")

		if !key.IsValid() || !value.IsValid() {
			return nil
		}
		tags[*key.Interface().(*string)] = value.Interface().(*string)
	}

	return tags
}
