package helper

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Dump is function to dump variable to string
func Dump(v any) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(marshal)
}

// TimeToStringFormat is function to convert time to string
func TimeToStringFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ExpectNumber is function to convert any to number
func ExpectNumber[T ~int | ~uint | ~float32 | ~float64](v any) T {
	var result T
	var valueString string

	// convert to string
	switch v.(type) {
	case string:
		valueString = v.(string)
	default:
		marshal, _ := json.Marshal(&v)
		if err := json.Unmarshal(marshal, &valueString); err != nil {
			logrus.Error(err)
			return T(0)
		}
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint:
		// if number float to int
		if strings.Contains(valueString, ".") {
			valueString = strings.Split(valueString, ".")[0]
		}

		number, err := strconv.ParseInt(valueString, 10, 64)
		if err != nil {
			logrus.Error(err)
			return T(0)
		}

		result = T(number)
	case reflect.Float64, reflect.Float32:
		float, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			logrus.Error(err)
			return T(0)
		}

		result = T(float)
	}

	return result
}
