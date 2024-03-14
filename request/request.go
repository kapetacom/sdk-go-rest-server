// Copyright 2023 Kapeta Inc.
// SPDX-License-Identifier: MIT
package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/ggicci/httpin"
	"github.com/labstack/echo/v4"
)

func GetRequestParameters[T any](req *http.Request, param *T) error {
	paramHandler, err := httpin.New(param)
	if err != nil {
		return err
	}
	paramValues, err := paramHandler.Decode(req)
	if err != nil {
		return err
	}
	castParamValues := paramValues.(*T)
	*param = *castParamValues
	return nil
}

// GetBody function takes two arguments: an echo context and a pointer to the return value.
// It used a JSON decoder to convert the request body into the return value.
// If the decoding fails, the function returns an error.
func GetBody[T any](ctx echo.Context, returnValue *T) error {
	err := json.NewDecoder(ctx.Request().Body).Decode(returnValue)
	if err != nil {
		return err
	}
	return nil
}

// GetPathParams function takes three arguments: an echo context, a string key, and a pointer to the return value.
// It returns the value of the key from the path parameters.
// If the key is not found, the function returns an error.
func GetPathParams[T any](ctx echo.Context, key string, returnValue *T) error {
	vals := ctx.ParamValues()
	keys := ctx.ParamNames()
	for i, k := range keys {
		if k == key {
			val := vals[i]
			return convertToType[T](returnValue, val)
		}
	}
	return fmt.Errorf("key not found")
}

// GetQueryParam function takes three arguments: an echo context, a string key, and a pointer to the return value.
// It returns the value of the key from the query parameters like /test?key=value
// If the key is not found, the function returns an error.
func GetQueryParam[T any](ctx echo.Context, key string, returnValue *T, optional bool) error {
	vals := ctx.Request().URL.Query()[key]
	if len(vals) == 0 {
		if optional {
			return nil
		}
		return fmt.Errorf("key '%v' not found", key)
	}
	return convertToType[T](returnValue, strings.Join(vals, ","))
}

// GetHeaderParams function takes three arguments: an echo context, a string key, and a pointer to the return value.
// It returns the value of the key from the header parameters.
// If the key is not found, the function returns an error.
func GetHeaderParams[T any](ctx echo.Context, key string, returnValue *T, optional bool) error {
	vals := ctx.Request().Header.Get(key)
	if vals == "" {
		if optional {
			return nil
		}
		return fmt.Errorf("key '%v' not found", key)
	}
	return convertToType[T](returnValue, vals)
}

func convertToType[T any](target *T, val string) error {

	var bodyBytes []byte
	var err error

	rt := reflect.TypeOf(*target)
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
			return err
		}
		bodyBytes, err = json.Marshal(x)
		if err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		bodyBytes, err = json.Marshal(x)
		if err != nil {
			return err
		}
	case reflect.Int:
		x, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		bodyBytes, err = json.Marshal(x)
		if err != nil {
			return err
		}
	default:
		bodyBytes, err = json.Marshal(val)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, target)
	if err != nil {
		return err
	}
	return nil
}
