package http

import (
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// Result represents a json value that is returned from Json().
// https://pkg.go.dev/github.com/tidwall/gjson#readme-result-type
type Result = gjson.Result

var (
	client *resty.Client
	null   Result
)

func init() {
	client = resty.New()
}

// JsonString is requests RESTful API then returns the response parsed value (string).
// https://pkg.go.dev/github.com/tidwall/gjson#readme-path-syntax
func JsonString(url, pattern string) (string, error) {
	ret, err := Json(url, pattern)
	if err != nil {
		return "", err
	}

	return ret.String(), nil
}

// Json is requests RESTful API then returns the response parsed value.
// https://pkg.go.dev/github.com/tidwall/gjson#readme-path-syntax
func Json(url, pattern string) (Result, error) {
	res, err := client.R().
		SetHeader("accept", "application/json").
		Get(url)
	if err != nil {
		return null, err
	}
	if res.StatusCode() != http.StatusOK {
		return null, errors.New(res.Status())
	}

	return gjson.GetBytes(res.Body(), pattern), nil
}
