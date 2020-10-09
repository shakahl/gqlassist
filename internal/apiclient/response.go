package apiclient

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ApiResponseParsed struct {
	Data   map[string]interface{} `json:"data"`
	Errors GQLApiErrors           `json:"errors"`
}

type ApiResponse struct {
	Data     *json.RawMessage `json:"data"`
	Errors   GQLApiErrors     `json:"errors"`
	response *http.Response   `json:"-"`
	body     []byte           `json:"-"`
}

func (r *ApiResponse) GetResponse() *http.Response {
	return r.response
}

func (r *ApiResponse) GetBody() []byte {
	return r.body
}

func (r *ApiResponse) GetStringBody() string {
	return string(r.GetBody())
}

func (r *ApiResponse) GetStringReader() *strings.Reader {
	return strings.NewReader(r.GetStringBody())
}

func (r *ApiResponse) DecodeIntoInterface() (target interface{}, err error) {
	err = json.NewDecoder(r.GetStringReader()).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}
