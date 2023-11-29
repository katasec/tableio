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

// GetStructValues Gets struct values of the fields in a struct
// used for the 'VALUES' clause in an SQL INSERT statement
func GetStructValuesForInsert[T any](data T) string {

	var sb strings.Builder

	// Use reflection to inspect struct values
	dataValue := reflect.ValueOf(data)
	dataTypes := reflect.TypeOf(data)

	// Iterate thru fields in struct
	numFields := dataValue.NumField()
	for i := 0; i < numFields; i++ {
		// Skip ID field for inserts
		if dataTypes.Field(i).Name == "ID" {
			continue
		}
		// Get current field value
		myValue := getValueFromStruct(dataValue.Field(i))

		// Comma separated
		sb.WriteString("'" + myValue + "',")
	}

	result := TrimSuffix(sb.String(), ",")

	return result
}

// getValueFromStruct Gets the value of the passed in field in a struct.
// If the field is a string, return as is. If the field is an int, convert
// to string. If the field is a struct, convert to JSON string. Otherwise,
// panic. Return as string for use in the VALUE clause for INSERT statements
// generate by TableIO
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

/*
GenSqlForFields Given a list of fields, generate SQL for the fields section in a CREATE TABLE statement. For e.g.

	`
	message VARCHAR(255) NULL,
	age INTEGER NULL,
	address TEXT NULL
	`
*/
func GenSqlForFields(fields []FieldInfo, driverName string) string {
	var sb strings.Builder
	var sql string

	// Loop through fields to generate SQL
	for i, field := range fields {

		// Generate SQL for current field
		if (field.FieldName) == "ID" {
			// Generate SQL for field based on field name for ID field
			// Note ID fields definition varies by database
			switch driverName {
			case "sqlite3":
				sql = "\tid INTEGER PRIMARY KEY AUTOINCREMENT"
			case "mysql":
				sql = "\tid INT PRIMARY KEY AUTO_INCREMENT"
			case "postgres":
				sql = "\tid SERIAL PRIMARY KEY"
			}
			sb.WriteString(sql)
		} else if (field.FieldName) == "Name" {
			// Generate SQL for field based on field name for Name field
			sql = "\tname VARCHAR(255) NOT NULL UNIQUE"
			sb.WriteString(sql)
		} else {
			// Generate SQL for field based on field type
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
		}

		// Add comma if not last field
		if i < len(fields)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}

	// Return SQL
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
