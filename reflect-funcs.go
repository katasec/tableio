package tableio

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/katasec/tableio/reflectx"
	"github.com/katasec/utils/errx"
)

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

func GetDbColumnNames[T any]() string {

	var sb strings.Builder

	// Instantiate Struct of type T to use for reflection
	var myStruct T
	myType := reflect.TypeOf(myStruct)

	// Iterate through fields
	numFields := myType.NumField()
	for i := 0; i < numFields; i++ {

		// Get SQL field from strut tag for current field
		currField := myType.Field(i)
		colName := getDbColumnName(currField) + ","

		// Add column name to list
		sb.WriteString(colName)
	}

	// Remove trailing ","
	result := TrimSuffix(sb.String(), ",")

	// Return columns
	return result
}

// getDbColumnName Get the database column for the given struct field  from its struct tag
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

func GenSqlForFields(fields []reflectx.FieldInfo) string {
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

func genSqlByType(f reflect.StructField) string {

	//fieldName := f.Name
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

func ToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
