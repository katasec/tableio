package reflectx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type FieldInfo struct {
	FieldName string
	FieldType string
	DbTag     string
}

// GetDbStructFields Returns list of fields in a struct of type "T" that have
// a "db" tag
func GetDbStructFields[T any](tags ...string) []FieldInfo {
	var fieldInfo []FieldInfo
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

	// Iterate through fields in the Struct
	for i := 0; i < myType.NumField(); i++ {

		// Get current field from struct
		f := myType.Field(i)

		// Add field name and type to list of fields
		if f.Tag.Get(filter) != "" {
			fieldInfo = append(fieldInfo, FieldInfo{
				FieldName: f.Tag.Get(filter),
				FieldType: f.Type.Name(),
			})
		}

	}

	return fieldInfo
}

// GetStructValues Gets struct values of the fields in a struct
// used for the 'VALUES' clause in an SQL INSERT statement
func GetStructValues[T any](data T) string {

	var sb strings.Builder

	// Use reflection to inspect struct values
	dataValue := reflect.ValueOf(data)

	// Iterate thru fields in struct
	numFields := dataValue.NumField()
	for i := 0; i < numFields; i++ {
		// Get next field and value
		myValue := getValueFromStruct(dataValue.Field(i))

		// Comma separated
		sb.WriteString("'" + myValue + "',")
	}

	result := TrimSuffix(sb.String(), ",")

	return result
}

func getValueFromStruct(field reflect.Value) string {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.String:
		return field.String()
	default:
		bVal, err := json.Marshal(field.Interface())
		if err != nil {
			panic(err)
		}
		return string(bVal)
	}
}

func GenSqlForFields(fields []FieldInfo) string {
	var sb strings.Builder
	var sql string
	for i, field := range fields {
		switch field.FieldType {
		case "string":
			sql = fmt.Sprintf("\t%s VARCHAR(255) NULL", field.FieldName)
			sb.WriteString(sql)
		case "int32":
			sql = fmt.Sprintf("\t%s INTEGER NULL", field.FieldName)
			sb.WriteString(sql)
		default:
			sql = fmt.Sprintf("\t%s TEXT NULL", field.FieldName)
			sb.WriteString(sql)
		}

		if i < len(fields)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}

	}
	//fieldName := f.Name
	return sb.String()

}

// GetDbStructFields Returns list of fields in a struct of type "T" that have
// a "db" tag
func GetStructFields[T any]() []FieldInfo {
	var fieldInfo []FieldInfo

	// Instantiate Struct of type T to use for reflection
	var myStruct T
	myType := reflect.TypeOf(myStruct)

	// Iterate through fields in the Struct
	for i := 0; i < myType.NumField(); i++ {

		// Get current field from struct
		f := myType.Field(i)

		// Add field name and type to list of fields
		fieldInfo = append(fieldInfo, FieldInfo{
			FieldName: f.Name,
			FieldType: f.Type.Name(),
		})

	}

	return fieldInfo
}

func GetDbStuff[T any]() {
	// Instantiate Struct of type T to use for reflection
	var myStruct T
	myType := reflect.TypeOf(myStruct)

	// Iterate through fields in the Struct
	for i := 0; i < myType.NumField(); i++ {

		// Get current field from struct
		f := myType.Field(i)

		// Add field name and type to list of fields
		fmt.Println(f.Name, f.Type.Name())

	}
}
