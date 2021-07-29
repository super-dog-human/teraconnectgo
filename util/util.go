package util

import (
	"fmt"
	"reflect"
)

func MergeStruct(new interface{}, orig interface{}) {
	allKeys := AllKeys(new)
	n := reflect.ValueOf(new).Elem()
	o := reflect.ValueOf(orig).Elem()

	for i := 0; i < o.NumField(); i++ {
		fieldType := o.Type().Field(i).Type
		fieldName := o.Type().Field(i).Name
		field := o.FieldByName(fieldName)
		fmt.Printf("diff: %v, %s\n", field, fieldName)

		if !Contains(&allKeys, fieldName) {
			continue // マージ元にないフィールドは無視する
		}

		if !o.CanSet() {
			continue
		}

		diffField := n.FieldByName(fieldName)
		switch fieldType.String() {
		case "interface":
			MergeStruct(diffField, field)
		case "bool":
			field.SetBool(diffField.Interface().(bool))
		case "float32":
		case "float64":
			field.SetFloat(diffField.Interface().(float64))
		case "int32":
		case "int64":
			field.SetInt(diffField.Interface().(int64))
		case "string":
			field.SetString(diffField.Interface().(string))
		case "uint64":
			field.SetUint(diffField.Interface().(uint64))
		}
	}
}

func AllKeys(target interface{}) []string {
	var keys []string
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < v.NumField(); i++ {
		keys = append(keys, v.Type().Field(i).Name)
	}

	return keys
}

func Contains(s *[]string, e string) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}
