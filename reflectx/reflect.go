package reflectx

import (
	"fmt"
	"reflect"
)

type FieldInfo struct {
	FieldName string
	FieldType string
	DbTag     string
}

func GetStructFields[T any](tags ...string) []string {
	var fields []string
	var filter string

	if len(tags) > 0 {
		filter = tags[0]
	}
	// Instantiate Struct of type T to use for reflection
	var myStruct T
	myType := reflect.TypeOf(myStruct)

	// Iterate through fields
	numFields := myType.NumField()
	for i := 0; i < numFields; i++ {

		// Get f info
		f := myType.Field(i)

		if filter != "" {
			if f.Tag.Get(filter) != "" {
				fields = append(fields, f.Name)
			}
		} else {
			fields = append(fields, f.Name)
		}

		//f.Tag.Get("db")

		// Add column name to list
		fmt.Printf("%d. The name is  %s\n", i, f.Name)
	}

	return fields
}

func GetStructFieldsX[T any](tags ...string) []FieldInfo {
	var filter string

	var fieldInfo []FieldInfo

	// DB fields have a 'db' tag
	if len(tags) > 0 {
		filter = tags[0]
	} else {
		filter = "db"
	}

	// Instantiate Struct of type T to use for reflection
	var myStruct T
	myType := reflect.TypeOf(myStruct)

	// Iterate through fields
	numFields := myType.NumField()
	for i := 0; i < numFields; i++ {

		// Get f info
		f := myType.Field(i)

		if filter != "" {
			if f.Tag.Get(filter) != "" {
				fieldInfo = append(fieldInfo, FieldInfo{
					FieldName: f.Name,
					FieldType: f.Type.Name(),
				})
			}
		} else {
			fieldInfo = append(fieldInfo, FieldInfo{
				FieldName: f.Name,
				FieldType: f.Type.Name(),
			})
		}

	}

	return fieldInfo
}

func GetDbStructFields[T any](tags ...string) []string {
	var fields []string
	var filter string

	// DB fields have a 'db' tag
	if len(tags) > 0 {
		filter = tags[0]
	} else {
		filter = "db"
	}

	// Instantiate Struct of type T to use for reflection
	var myStruct T
	myType := reflect.TypeOf(myStruct)

	// Iterate through fields
	numFields := myType.NumField()
	for i := 0; i < numFields; i++ {

		// Get f info
		f := myType.Field(i)

		if f.Tag.Get(filter) != "" {
			fields = append(fields, f.Tag.Get(filter))
		}

	}

	return fields
}
