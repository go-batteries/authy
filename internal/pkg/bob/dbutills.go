package bob

import (
	"database/sql/driver"
	"net/url"
	"strings"
)

type StringArray []string

func (v StringArray) Value() (driver.Value, error) {
	return stringArrayToString(v), nil
}

func stringArrayToString(arr StringArray) string {
	encodedArr := []string{}

	for _, v := range arr {
		encodedArr = append(encodedArr, url.QueryEscape(v))
	}

	return strings.Join(encodedArr, ",")
}

func (v *StringArray) Scan(value interface{}) error {
	values := []string{}

	if v == nil {
		*v = values
		return nil
	}

	if res, err := driver.String.ConvertValue(value); err == nil {
		arr, ok := res.(string)
		if !ok {
			*v = values
			return nil
		}

		for _, v := range strings.Split(arr, ",") {
			x, err := url.QueryUnescape(v)
			if err != nil {
				return err
			}
			values = append(values, x)
		}
	}

	*v = values
	return nil
}
