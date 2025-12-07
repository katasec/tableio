package reflectx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// parseTableIOTag parses the tableio struct tag and returns a map of attributes
// Supports: pk, auto, unique, required
func parseTableIOTag(tag string) map[string]bool {
	attrs := make(map[string]bool)
	if tag == "" {
		return attrs
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch part {
		case "pk", "primarykey":
			attrs["pk"] = true
		case "auto", "autoincrement":
			attrs["auto"] = true
		case "unique":
			attrs["unique"] = true
		case "required", "notnull":
			attrs["required"] = true
		}
	}
	return attrs
}

type FieldInfo struct {
	FieldName     string
	FieldType     string
	DbTag         string
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Required      bool
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

// GetStructValuesForInsert Gets struct values of the fields in a struct
// used for the 'VALUES' clause in an SQL INSERT statement
// Skips auto-increment fields
func GetStructValuesForInsert[T any](data T) string {

	var values []string

	// Use reflection to inspect struct values
	dataValue := reflect.ValueOf(data)

	// Get field info to check for auto-increment
	fieldInfos := GetStructFields[T]()

	// Iterate thru fields in struct
	numFields := dataValue.NumField()
	for i := 0; i < numFields; i++ {
		// Skip auto-increment fields (database generates these)
		if fieldInfos[i].AutoIncrement {
			continue
		}
		// Get current field value
		myValue := getValueFromStruct(dataValue.Field(i))

		// Add to values list
		values = append(values, "'"+myValue+"'")
	}

	return strings.Join(values, ",")
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

	// Loop through fields to generate SQL
	for i, field := range fields {

		// Build column definition based on field attributes
		var colDef strings.Builder
		colDef.WriteString("\t" + field.FieldName + " ")

		// Determine column type
		if field.PrimaryKey && field.AutoIncrement {
			// Auto-increment primary key - database-specific syntax
			switch driverName {
			case "sqlite3":
				colDef.WriteString("INTEGER PRIMARY KEY AUTOINCREMENT")
			case "mysql":
				colDef.WriteString("INT PRIMARY KEY AUTO_INCREMENT")
			case "postgres":
				colDef.WriteString("SERIAL PRIMARY KEY")
			case "mssql", "sqlserver":
				colDef.WriteString("INT PRIMARY KEY IDENTITY(1,1)")
			}
		} else {
			// Regular field - determine type
			switch field.FieldType {
			case "string":
				colDef.WriteString("VARCHAR(255)")
			case "int32", "int64", "int":
				colDef.WriteString("INTEGER")
			default:
				// Nested structs and other types - use driver-specific column type
				switch driverName {
				case "sqlite3":
					colDef.WriteString("TEXT")
				case "mssql", "sqlserver":
					colDef.WriteString("NVARCHAR(MAX)")
				default:
					colDef.WriteString("JSON")
				}
			}

			// Add constraints
			if field.PrimaryKey {
				colDef.WriteString(" PRIMARY KEY")
			}
			if field.Required {
				colDef.WriteString(" NOT NULL")
			}
			if field.Unique {
				colDef.WriteString(" UNIQUE")
			}
			if !field.Required && !field.PrimaryKey {
				colDef.WriteString(" NULL")
			}
		}

		sb.WriteString(colDef.String())

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

// GetStructFields Returns list of fields in a struct of type "T" with parsed tableio tags
func GetStructFields[T any]() []FieldInfo {
	var fieldInfo []FieldInfo

	// Instantiate Struct of type T to use for reflection
	var myStruct T
	myType := reflect.TypeOf(myStruct)

	// Iterate through fields in the Struct
	for i := 0; i < myType.NumField(); i++ {

		// Get current field from struct
		f := myType.Field(i)

		// Parse tableio tag
		tableioTag := f.Tag.Get("tableio")
		attrs := parseTableIOTag(tableioTag)

		// Add field name and type to list of fields
		fieldInfo = append(fieldInfo, FieldInfo{
			FieldName:     f.Name,
			FieldType:     f.Type.Name(),
			PrimaryKey:    attrs["pk"],
			AutoIncrement: attrs["auto"],
			Unique:        attrs["unique"],
			Required:      attrs["required"],
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
