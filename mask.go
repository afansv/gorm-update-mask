package gorm_update_mask

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/stoewer/go-strcase"
)

func Mask(model interface{}, updatePaths []string) map[string]interface{} {
	t := reflect.TypeOf(model)

	fields := structs.Names(model)

	mask := make(map[string]interface{})

	for _, fName := range fields {
		field, _ := t.FieldByName(fName)

		columnName, ok := parseColumnName(field)
		if !ok {
			columnName = strcase.SnakeCase(fName)
		}

		if isStringInArray(columnName, updatePaths) {
			mask[columnName] = reflect.ValueOf(model).FieldByName(fName).Interface()
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