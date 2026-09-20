// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Schema is inline JSON Schema accepted by fm respond.
type Schema json.RawMessage

// ParseSchema validates inline JSON Schema.
func ParseSchema(data []byte) (Schema, error) {
	if !json.Valid(data) {
		return nil, fmt.Errorf("fmgo: invalid schema JSON")
	}
	return Schema(append([]byte(nil), data...)), nil
}

// SchemaType identifies a property declaration accepted by fm schema object.
type SchemaType string

const (
	SchemaString  SchemaType = "string"
	SchemaInteger SchemaType = "integer"
	SchemaDouble  SchemaType = "double"
	SchemaBoolean SchemaType = "boolean"
	SchemaObject  SchemaType = "object"
	SchemaAnyOf   SchemaType = "anyOf"
)

// SchemaProperty describes a native fm schema property.
type SchemaProperty struct {
	Name        string
	Type        SchemaType
	Description string
	Optional    bool
	Array       bool
	Schema      Schema
	Choices     []Schema
}

// ObjectSchema describes an object passed to fm schema object.
type ObjectSchema struct {
	Name       string
	Properties []SchemaProperty
}

// GenerateSchema creates a schema through fm schema object.
func (c *Client) GenerateSchema(ctx context.Context, object ObjectSchema) (Schema, error) {
	if object.Name == "" {
		return nil, fmt.Errorf("fmgo: schema object name is required")
	}
	args := []string{"schema", "object", "--name", object.Name}
	for _, property := range object.Properties {
		if property.Name == "" {
			return nil, fmt.Errorf("fmgo: schema property name is required")
		}
		switch property.Type {
		case SchemaString, SchemaInteger, SchemaDouble, SchemaBoolean:
			args = append(args, "--"+string(property.Type), property.Name)
		case SchemaObject:
			if len(property.Schema) == 0 {
				return nil, fmt.Errorf("fmgo: object property %q requires a schema", property.Name)
			}
			args = append(args, "--object", property.Name, "--schema", string(property.Schema))
		case SchemaAnyOf:
			if len(property.Choices) == 0 {
				return nil, fmt.Errorf("fmgo: anyOf property %q requires choices", property.Name)
			}
			args = append(args, "--anyOf")
			for _, choice := range property.Choices {
				args = append(args, "--schema", string(choice))
			}
		default:
			return nil, fmt.Errorf("fmgo: unsupported schema property type %q", property.Type)
		}
		if property.Array {
			args = append(args, "--array")
		}
		if property.Description != "" {
			args = append(args, "--description", property.Description)
		}
		if property.Optional {
			args = append(args, "--optional")
		}
	}
	stdout, _, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return ParseSchema(stdout)
}

// SchemaFor generates JSON Schema for T without invoking fm.
func SchemaFor[T any]() (Schema, error) {
	var value T
	return SchemaFrom(value)
}

// SchemaFrom generates JSON Schema for a Go value's type.
func SchemaFrom(value any) (Schema, error) {
	typeOf := reflect.TypeOf(value)
	if typeOf == nil {
		return nil, fmt.Errorf("fmgo: schema requires a concrete type")
	}
	for typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	visiting := map[reflect.Type]bool{}
	schema, err := schemaForType(typeOf, visiting)
	if err != nil {
		return nil, err
	}
	if typeOf.Kind() == reflect.Struct && typeOf.Name() != "" {
		schema["title"] = typeOf.Name()
	}
	encoded, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("fmgo: encode schema: %w", err)
	}
	return Schema(encoded), nil
}

func schemaForType(typeOf reflect.Type, visiting map[reflect.Type]bool) (map[string]any, error) {
	for typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	switch typeOf.Kind() {
	case reflect.String:
		return map[string]any{"type": "string"}, nil
	case reflect.Bool:
		return map[string]any{"type": "boolean"}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return map[string]any{"type": "integer"}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return map[string]any{"type": "integer", "minimum": 0}, nil
	case reflect.Uint64:
		return nil, fmt.Errorf("fmgo: unsupported schema type %s", typeOf)
	case reflect.Float32, reflect.Float64:
		return map[string]any{"type": "number"}, nil
	case reflect.Slice, reflect.Array:
		items, err := schemaForType(typeOf.Elem(), visiting)
		if err != nil {
			return nil, err
		}
		return map[string]any{"type": "array", "items": items}, nil
	case reflect.Struct:
		return schemaForStruct(typeOf, visiting)
	default:
		return nil, fmt.Errorf("fmgo: unsupported schema type %s", typeOf)
	}
}

func schemaForStruct(typeOf reflect.Type, visiting map[reflect.Type]bool) (map[string]any, error) {
	if visiting[typeOf] {
		return nil, fmt.Errorf("fmgo: recursive schema type %s", typeOf)
	}
	visiting[typeOf] = true
	defer delete(visiting, typeOf)

	properties := map[string]any{}
	required := []string{}
	order := []string{}
	for index := range typeOf.NumField() {
		field := typeOf.Field(index)
		if !field.IsExported() || field.Anonymous {
			continue
		}
		name, optional := jsonField(field)
		if name == "" {
			continue
		}
		property, err := schemaForType(field.Type, visiting)
		if err != nil {
			return nil, err
		}
		properties[name] = property
		order = append(order, name)
		if !optional {
			required = append(required, name)
		}
	}
	schema := map[string]any{
		"type":                 "object",
		"x-order":              order,
		"properties":           properties,
		"additionalProperties": false,
	}
	if len(required) != 0 {
		schema["required"] = required
	}
	return schema, nil
}

func jsonField(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	name, options, _ := strings.Cut(tag, ",")
	if name == "-" {
		return "", false
	}
	if name == "" {
		name = field.Name
	}
	optional := false
	for _, option := range strings.Split(options, ",") {
		if option == "omitempty" {
			optional = true
		}
	}
	return name, optional
}
