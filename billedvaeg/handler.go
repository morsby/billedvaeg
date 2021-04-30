package function

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"github.com/morsby/billedvaeg/api/web"
	handler "github.com/openfaas/templates-sdk/go-http"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var err error
	var body []byte
	status := http.StatusOK

	if req.Method == http.MethodGet {
		var b bytes.Buffer
		web.Compile(&b)
		body = b.Bytes()
	} else {
		errMsg := "Method not implemented"
		err = errors.New(strings.ToLower(errMsg))
		body = []byte(errMsg)
	}

	return handler.Response{
		Body:       body,
		StatusCode: status,
	}, err
}
