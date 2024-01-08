package utils

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/util/nulltype"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func GetQueryParamInt(r *http.Request, key string) (int, error) {
	strVal := GetQueryParam(r, key)

	if len(strVal) == 0 {
		return -1, nil
	}

	return strconv.Atoi(strVal)
}

func GetQueryParamInt64(r *http.Request, key string) (int64, error) {
	strVal := GetQueryParam(r, key)

	if len(strVal) == 0 {
		return -1, nil
	}

	return strconv.ParseInt(strVal, 10, 64)
}

func GetQueryParamInt64Slice(r *http.Request, key string) ([]int64, error) {
	queryParamValue := GetQueryParam(r, key)

	if !strings.Contains(queryParamValue, ",") {
		intVal, err := GetQueryParamInt64(r, key)
		if err != nil || intVal == -1 {
			return []int64{}, err
		}
		return []int64{intVal}, nil
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

func GetQueryParamFloat64(r *http.Request, key string) (float64, error) {
	return GetQueryParamFloat64Fallback(r, key, -1)
}

func GetQueryParamFloat64Fallback(r *http.Request, key string, fallback float64) (float64, error) {
	strVal := GetQueryParam(r, key)

	if len(strVal) == 0 {
		return fallback, nil
	}

	return strconv.ParseFloat(strVal, 64)
}

func GetQueryParamNullBoolean(r *http.Request, key string) nulltype.NullBoolean {
	intQueryParam, err := GetQueryParamInt(r, key)

	if err != nil {
		return nulltype.UndefinedValue
	}

	return nulltype.FromInt(intQueryParam)
}

func GetTimestampFromQueryParam(r *http.Request, key string, fallback time.Time) (time.Time, error) {
	seconds, err := GetQueryParamInt64(r, key)

	if err != nil {
		return fallback, fmt.Errorf(
			"invalid timestamp `%s` value `%s`: %w",
			key, GetQueryParam(r, key), err,
		)
	}

	if seconds == -1 {
		return fallback, nil
	}

	return time.Unix(seconds, 0), nil
}
