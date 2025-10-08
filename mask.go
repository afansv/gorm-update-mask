package gorm_update_mask

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/stoewer/go-strcase"
)

func Mask(model interface{}, updatePaths []string) map[string]interface{} {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return make(map[string]interface{})
		}
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return make(map[string]interface{})
	}

	fields := structs.Names(model)

	mask := make(map[string]interface{})

	for _, fName := range fields {
		field, ok := t.FieldByName(fName)
		if !ok {
			continue
		}

		columnName, ok := parseColumnName(field)
		if !ok {
			columnName = strcase.SnakeCase(fName)
		}

		if isStringInArray(columnName, updatePaths) {
			fieldValue := v.FieldByName(fName)
			if fieldValue.IsValid() {
				mask[columnName] = fieldValue.Interface()
			}
		}
	}

	return mask
}

func isStringInArray(s string, ar []string) bool {
	for _, a := range ar {
		if a == s {
			return true
		}
	}

	return false
}

func parseColumnName(field reflect.StructField) (columnName string, ok bool) {
	tag, ok := field.Tag.Lookup("gorm")
	if !ok {
		tag, ok = field.Tag.Lookup("sql")
		if !ok {
			return "", false
		}
	}

	structValues := strings.Split(tag, ",")
	for _, v := range structValues {
		kv := strings.Split(v, ":")
		if len(kv) != 2 {
			continue
		}

		if kv[0] == "column" {
			return kv[1], true
		}
	}

	return "", false
}
