package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
	"golang.org/x/net/context/ctxhttp"
)

type GraphQLOperation string

const (
	GraphQLQuery    GraphQLOperation = "query"
	GraphQLMutation GraphQLOperation = "mutation"
)

type ApiClient struct {
	Endpoint   string
	httpClient *http.Client
	Header     http.Header
}

func NewApiClient(endpoint string, httpClient *http.Client) *ApiClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &ApiClient{
		Endpoint:   endpoint,
		httpClient: httpClient,
		Header:     make(http.Header),
	}
}

func (api *ApiClient) SendGraphQLQuery(ctx context.Context, query string, variables map[string]interface{}) (*ApiResponse, error) {
	payload, err := api.CreateGraphQLPayload(query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "GraphQL query error!")
	}
	resp, err := api.Post(ctx, payload, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result = ApiResponse{
		response: resp,
	}

	result.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return &result, err
	}

	err = json.NewDecoder(bytes.NewReader(result.body)).Decode(&result)
	if err != nil {
		return &result, err
	}

	return &result, err
}

func (api *ApiClient) CreateGraphQLPayload(query string, variables map[string]interface{}) (*bytes.Buffer, error) {
	var input = struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(input)
	if err != nil {
		return nil, errors.Wrap(err, "Error while encoding GraphQL payload!")
	}
	return &buf, nil
}

func (api *ApiClient) parsePayload(p interface{}) (*bytes.Buffer, error) {
	var buf *bytes.Buffer
	switch p.(type) {
	case string:
		buf = bytes.NewBufferString(p.(string))
		break
	case *string:
		buf = bytes.NewBufferString(*(p.(*string)))
		break
	case []byte:
		buf = bytes.NewBuffer(p.([]byte))
		break
	case bytes.Buffer:
		b := p.(bytes.Buffer)
		buf = &b
		break
	case *bytes.Buffer:
		buf = p.(*bytes.Buffer)
		break
	default:
		return nil, errors.New("Invalid type for payload: " + reflect.TypeOf(p).String())
	}
	return buf, nil
}

func (api *ApiClient) addHeaders(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	for h := range api.Header {
		req.Header.Set(h, api.Header.Get(h))
	}
	return req
}

func (api *ApiClient) Post(ctx context.Context, payload interface{}, modify *func(r *http.Request) *http.Request) (*http.Response, error) {
	buf, err := api.parsePayload(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", api.Endpoint, buf)
	if err != nil {
		return nil, err
	}
	if modify != nil {
		req = (*modify)(req)
	}

	resp, err := ctxhttp.Do(ctx, api.httpClient, api.addHeaders(req))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return resp, fmt.Errorf("Non-2xx OK status code: status=%v resp=%v", resp.Status, resp)
	}

	return resp, err
}
