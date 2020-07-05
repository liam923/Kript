package encode

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Convert a struct or map of any type to a map with string keys.
// An error will result if the given object is not a struct or map, or if it contains a pointer.
// Additionally, special structs that are types recognized by firestore will not be converted.
func ToMap(o interface{}, tag string) (map[string]interface{}, error) {
	r, err := toMap(reflect.ValueOf(o), tag)
	if err != nil {
		return nil, err
	}
	if m, ok := r.(map[string]interface{}); ok {
		return m, nil
	} else {
		return nil, fmt.Errorf("atomic type given")
	}
}

// do the actual recursion
func toMap(val reflect.Value, tag string) (interface{}, error) {
	switch val.Type().Kind() {
	case reflect.Map:
		m := make(map[string]interface{}, val.Len())
		for _, key := range val.MapKeys() {
			if key.Kind() != reflect.String {
				return nil, fmt.Errorf("cannot have non-string map key: %v", key)
			}
			sub, err := toMap(val.MapIndex(key), tag)
			if err != nil {
				return nil, err
			}
			m[key.String()] = sub
		}
		return m, nil
	case reflect.Array, reflect.Slice:
		s := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			sub, err := toMap(val.Index(i), tag)
			if err != nil {
				return nil, err
			}
			s[i] = sub
		}
		return s, nil
	case reflect.Struct:
		if o, ok := isSpecialStruct(val); ok {
			return o, nil
		}
		m := map[string]interface{}{}
		structType := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.CanInterface() {
				key, omitempty := parseTag(structType.Field(i), tag)
				if !omitempty || !field.IsZero() {
					sub, err := toMap(field, tag)
					if err != nil {
						return nil, err
					}
					m[key] = sub
				}
			}
		}
		return m, nil
	case reflect.Ptr:
		return nil, fmt.Errorf("cannot unwrap pointer")
	default:
		if val.CanInterface() {
			return val.Interface(), nil
		} else {
			return nil, fmt.Errorf("unable to interface reflect value: %v", val)
		}
	}
}

// check if the given value is a time.Time
func isSpecialStruct(val reflect.Value) (interface{}, bool) {
	if val.CanInterface() {
		v := val.Interface()
		switch v.(type) {
		case time.Time:
			return v, true
		}
	}
	return nil, false
}

func parseTag(fieldType reflect.StructField, tag string) (key string, omitempty bool) {
	s := strings.Split(fieldType.Tag.Get(tag), ",")
	key = s[0]
	if key == "" {
		key = fieldType.Name
	}
	if len(s) >= 2 {
		omitempty = s[1] == "omitempty"
	}
	return
}
