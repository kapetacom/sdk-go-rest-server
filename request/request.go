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

// GetBody function takes two arguments: an echo context and a pointer to the return value.
// It used a JSON decoder to convert the request body into the return value.
// If the decoding fails, the function returns an error.
func GetBody[T any](ctx echo.Context, returnValue *T) (*T, error) {
	err := json.NewDecoder(ctx.Request().Body).Decode(returnValue)
	if err != nil {
		return nil, err
	}
	return returnValue, nil
}

// GetPathParams function takes three arguments: an echo context, a string key, and a pointer to the return value.
// It returns the value of the key from the path parameters.
// If the key is not found, the function returns an error.
func GetPathParams[T any](ctx echo.Context, key string, returnValue *T) (T, error) {
	vals := ctx.ParamValues()
	keys := ctx.ParamNames()
	for i, k := range keys {
		if k == key {
			val := vals[i]
			return convertToType[T](returnValue, val)
		}
	}
	return *returnValue, fmt.Errorf("key not found")
}

// GetQueryParam function takes three arguments: an echo context, a string key, and a pointer to the return value.
// It returns the value of the key from the query parameters like /test?key=value
// If the key is not found, the function returns an error.
func GetQueryParam[T any](ctx echo.Context, key string, returnValue *T) (T, error) {
	vals := ctx.Request().URL.Query()[key]
	return convertToType[T](returnValue, strings.Join(vals, ","))
}

func convertToType[T any](returnValue *T, val string) (T, error) {

	var bodyBytes []byte
	var err error

	rt := reflect.TypeOf(*returnValue)
	kind := rt.Kind()
	switch kind {
	case reflect.Map:
		bodyBytes, err = json.Marshal(val)

	case reflect.Array | reflect.Slice:
		x := strings.Split(val, ",")
		bodyBytes, err = json.Marshal(x)
	case reflect.String:
		bodyBytes, err = json.Marshal(val)
	case reflect.Bool:
		x, err := strconv.ParseBool(val)
		if err != nil {
			return *returnValue, err
		}
		bodyBytes, err = json.Marshal(x)
		if err != nil {
			return *returnValue, err
		}
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return *returnValue, err
		}
		bodyBytes, err = json.Marshal(x)
		if err != nil {
			return *returnValue, err
		}
	case reflect.Int:
		x, err := strconv.Atoi(val)
		if err != nil {
			return *returnValue, err
		}
		bodyBytes, err = json.Marshal(x)
		if err != nil {
			return *returnValue, err
		}
	default:
		bodyBytes, err = json.Marshal(val)
	}
	if err != nil {
		return *returnValue, err
	}
	x := json.Unmarshal(bodyBytes, &returnValue)
	return *returnValue, x
}
