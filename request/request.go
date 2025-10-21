// Copyright 2023 Kapeta Inc.
// SPDX-License-Identifier: MIT
package request

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// FillStruct fills a struct with data from path, query, body, and headers.
// Tags format: in:"<source>=<key>;required"
// Supported sources: path, query, body, header
func FillStruct[T any](ctx echo.Context, result *T) error {
	// Decode body into a map
	body := map[string]any{}
	if ctx.Request().Body != nil {
		if err := json.NewDecoder(ctx.Request().Body).Decode(&body); err != nil && err.Error() != "EOF" {
			return fmt.Errorf("failed to decode body: %w", err)
		}
	}

	v := reflect.ValueOf(result).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		inTag, ok := field.Tag.Lookup("in")
		if !ok {
			continue
		}

		// Parse tag: "source=key;required"
		parts := strings.Split(inTag, ";")
		sourceKey := strings.SplitN(parts[0], "=", 2)
		if len(sourceKey) != 2 {
			return fmt.Errorf("invalid in tag format for field %s", field.Name)
		}
		source, key := sourceKey[0], sourceKey[1]

		required := false
		for _, p := range parts[1:] {
			if strings.ToLower(p) == "required" {
				required = true
			}
		}

		var val string
		switch source {
		case "path":
			val = ctx.Param(key)
		case "query":
			vals := ctx.QueryParams()[key]
			if len(vals) == 0 {
				val = ""
			} else if len(vals) == 1 {
				val = vals[0]
			} else {
				val = strings.Join(vals, ",")
			}
		case "body":
			if bodyVal, ok := body[key]; ok {
				val = fmt.Sprintf("%v", bodyVal)
			}
		case "header":
			val = ctx.Request().Header.Get(key)
		default:
			return fmt.Errorf("unsupported source: %s", source)
		}

		if required && val == "" {
			return fmt.Errorf("field %s is required but missing", field.Name)
		}

		if val != "" {
			if err := setFieldValue(fieldVal, val); err != nil {
				return fmt.Errorf("failed to set field %s: %w", field.Name, err)
			}
		}
	}

	return nil
}

// setFieldValue converts a string to the appropriate type and sets the reflect.Value
func setFieldValue(fieldVal reflect.Value, val string) error {
	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(val)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		fieldVal.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		fieldVal.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		fieldVal.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		fieldVal.SetFloat(f)
	case reflect.Slice:
		parts := strings.Split(val, ",")
		slice := reflect.MakeSlice(fieldVal.Type(), len(parts), len(parts))
		for i, p := range parts {
			if err := setFieldValue(slice.Index(i), p); err != nil {
				return err
			}
		}
		fieldVal.Set(slice)
	case reflect.Map:
		// For maps, try to parse JSON string
		m := reflect.New(fieldVal.Type()).Interface()
		if err := json.Unmarshal([]byte(val), &m); err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(m).Elem())
	default:
		return fmt.Errorf("unsupported field type: %s", fieldVal.Type())
	}
	return nil
}
