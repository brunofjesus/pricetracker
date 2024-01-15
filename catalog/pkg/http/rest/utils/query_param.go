package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func GetQueryParamInt(r *http.Request, key string) (*int, error) {
	strVal := GetQueryParam(r, key)

	if len(strVal) == 0 {
		return nil, nil
	}

	val, err := strconv.Atoi(strVal)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func GetQueryParamInt64(r *http.Request, key string) (*int64, error) {
	strVal := GetQueryParam(r, key)

	if len(strVal) == 0 {
		return nil, nil
	}

	val, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return nil, err
	}

	return &val, err
}

func GetQueryParamInt64Slice(r *http.Request, key string) ([]int64, error) {
	queryParamValue := GetQueryParam(r, key)

	if !strings.Contains(queryParamValue, ",") {
		intPtr, err := GetQueryParamInt64(r, key)
		if err != nil || intPtr == nil {
			return []int64{}, err
		}
		return []int64{*intPtr}, nil
	}

	queryParamValueSlice := strings.Split(queryParamValue, ",")
	result := make([]int64, 0, len(queryParamValueSlice))

	for _, str := range queryParamValueSlice {
		intVal, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return []int64{}, err
		}
		result = append(result, intVal)
	}

	return result, nil
}

func GetQueryParamFloat64(r *http.Request, key string) (*float64, error) {
	return GetQueryParamFloat64Fallback(r, key, nil)
}

func GetQueryParamFloat64Fallback(r *http.Request, key string, fallback *float64) (*float64, error) {
	strVal := GetQueryParam(r, key)

	if len(strVal) == 0 {
		return fallback, nil
	}

	val, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return nil, err
	}

	return &val, err
}

func GetQueryParamNullBoolean(r *http.Request, key string) *bool {
	intQueryParam, err := GetQueryParamInt(r, key)

	if err != nil || intQueryParam == nil {
		return nil
	}

	val := *intQueryParam > 0
	return &val
}

func GetTimestampFromQueryParam(r *http.Request, key string, fallback *time.Time) (*time.Time, error) {
	seconds, err := GetQueryParamInt64(r, key)

	if err != nil {
		return fallback, fmt.Errorf(
			"invalid timestamp `%s` value `%s`: %w",
			key, GetQueryParam(r, key), err,
		)
	}

	if seconds == nil {
		return fallback, nil
	}
	val := time.Unix(*seconds, 0)
	return &val, nil
}
