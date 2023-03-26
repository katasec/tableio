package reflectx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/katasec/utils/errx"
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
	for _, field := range fields {
		switch field.FieldType {
		case "string":
			sql = fmt.Sprintf("\t%s VARCHAR(255) NULL \n", field.FieldName)
			sb.WriteString(sql)
		case "int32":
			sql = fmt.Sprintf("\t%s INTEGER NULL \n", field.FieldName)
			sb.WriteString(sql)
		default:
			sql = fmt.Sprintf("\t%s TEXT NULL \n", field.FieldName)
			sb.WriteString(sql)
		}
	}
	//fieldName := f.Name
	return sb.String()

}

// GenSqlForStructFields Generates SQL for creating database columns for a given struct
func GenSqlForStructFields[T any]() string {
	var sb strings.Builder

	var myStruct T

	t := reflect.TypeOf(myStruct)
	numFields := t.NumField()

	suffix := ""

	for i := 0; i < numFields; i++ {

		// Get current field
		currField := t.Field(i)
		// Append comma for all except last field
		if i < numFields-1 {
			suffix = ","
		} else {
			suffix = ""
		}

		// Gen SQL for field and add to buffer
		sb.WriteString(genSqlByType(currField) + suffix + "\n")
	}
	return sb.String()
}

func genSqlByType(f reflect.StructField) string {

	colName := getDbColumnName(f)

	switch f.Type.String() {
	case "string":
		return fmt.Sprintf("\t%s VARCHAR(255) NULL", colName)
	case "int32":
		return fmt.Sprintf("\t%s INTEGER NULL", colName)
	default:
		return fmt.Sprintf("\t%s TEXT NULL", colName)
	}

}

func getDbColumnName(f reflect.StructField) string {
	// Get DB struct tag
	colName := string(f.Tag.Get("db"))

	// Panic if DB struct missing, provide example
	if f.Tag.Get("db") == "" {
		example := fmt.Sprintf("For e.g:\n\n\t %s %s `db:\"%s\"`", f.Name, f.Type.String(), ToSnakeCase(f.Name))
		errMsg := fmt.Sprintf("%s field is missing a 'db' struct tag. %s", f.Name, example)
		errx.PanicOnError(errors.New(errMsg))
	}
	// return col name
	return colName
}
