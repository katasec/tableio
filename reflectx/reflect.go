package reflectx

import (
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

	// Iterate through fields
	numFields := myType.NumField()
	for i := 0; i < numFields; i++ {

		// Get f info
		f := myType.Field(i)

		if f.Tag.Get(filter) != "" {
			fieldInfo = append(fieldInfo, FieldInfo{
				FieldName: f.Tag.Get(filter),
				FieldType: f.Type.Name(),
			})
		}

	}

	return fieldInfo
}

func GetTableName[T any]() string {
	var data []T

	// Tablename is dervied from the name of the objects's Type
	tableName := reflect.TypeOf(data).String()

	// Remove package names from the resulting string
	if strings.Contains(tableName, ".") {
		tableName = strings.Split(tableName, ".")[1]
	}

	// Add an S
	return tableName + "s"
}

// GetStructValues Gets struct values
func GetStructValues[T any](data T) string {

	var sb strings.Builder

	// Use reflection to inspect struct values
	dataValue := reflect.ValueOf(data)

	// Iterate thru fields in struct
	numFields := dataValue.NumField()
	for i := 0; i < numFields; i++ {
		// Get next field and value
		myValue := getValue(dataValue.Field(i))

		// Comma separated
		sb.WriteString("'" + myValue + "',")
	}

	result := TrimSuffix(sb.String(), ",")

	return result
}

func getValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.String:
		return field.String()
	}

	return ""
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
