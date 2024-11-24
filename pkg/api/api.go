// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for Unit.
const (
	UnitImperial Unit = "imperial"
	UnitMetric   Unit = "metric"
)

// Defines values for GetMultipleSensorsParamsUnit.
const (
	GetMultipleSensorsParamsUnitImperial GetMultipleSensorsParamsUnit = "imperial"
	GetMultipleSensorsParamsUnitMetric   GetMultipleSensorsParamsUnit = "metric"
)

// AccountResponse defines model for AccountResponse.
type AccountResponse struct {
	Id *string `json:"id,omitempty"`
}

// AccountsResponse defines model for AccountsResponse.
type AccountsResponse struct {
	Accounts *[]AccountResponse `json:"accounts,omitempty"`
}

// DeviceResponse defines model for DeviceResponse.
type DeviceResponse struct {
	Home         *string   `json:"home"`
	Name         *string   `json:"name,omitempty"`
	Sensors      *[]string `json:"sensors,omitempty"`
	SerialNumber *string   `json:"serialNumber,omitempty"`
	Type         *string   `json:"type,omitempty"`
}

// DevicesResponse defines model for DevicesResponse.
type DevicesResponse struct {
	Devices *[]DeviceResponse `json:"devices,omitempty"`
}

// Error defines model for Error.
type Error struct {
	// Message A message detailing the error encountered
	Message *string `json:"message,omitempty"`
}

// SensorResponse defines model for SensorResponse.
type SensorResponse struct {
	SensorType *string  `json:"sensorType,omitempty"`
	Unit       *string  `json:"unit,omitempty"`
	Value      *float64 `json:"value,omitempty"`
}

// SensorsResponse defines model for SensorsResponse.
type SensorsResponse struct {
	BatteryPercentage *int              `json:"batteryPercentage"`
	Recorded          *string           `json:"recorded"`
	Sensors           *[]SensorResponse `json:"sensors,omitempty"`
	SerialNumber      *string           `json:"serialNumber,omitempty"`
}

// AccountId defines model for accountId.
type AccountId = string

// DeviceSerialNumbers defines model for deviceSerialNumbers.
type DeviceSerialNumbers = []string

// PageNumber defines model for pageNumber.
type PageNumber = int

// Unit defines model for unit.
type Unit string

// GetMultipleSensorsParams defines parameters for GetMultipleSensors.
type GetMultipleSensorsParams struct {
	// Sn The serial numbers of the devices
	Sn *DeviceSerialNumbers `form:"sn,omitempty" json:"sn,omitempty"`

	// PageNumber The number of a page (of 50 records) to fetch
	PageNumber *PageNumber `form:"pageNumber,omitempty" json:"pageNumber,omitempty"`

	// Unit The units type sensors values will be returned in
	Unit *GetMultipleSensorsParamsUnit `form:"unit,omitempty" json:"unit,omitempty"`
}

// GetMultipleSensorsParamsUnit defines parameters for GetMultipleSensors.
type GetMultipleSensorsParamsUnit string

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetAccountsIds request
	GetAccountsIds(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetDevices request
	GetDevices(ctx context.Context, accountId AccountId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetMultipleSensors request
	GetMultipleSensors(ctx context.Context, accountId AccountId, params *GetMultipleSensorsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetHealth request
	GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAccountsIds(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAccountsIdsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetDevices(ctx context.Context, accountId AccountId, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetDevicesRequest(c.Server, accountId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetMultipleSensors(ctx context.Context, accountId AccountId, params *GetMultipleSensorsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetMultipleSensorsRequest(c.Server, accountId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAccountsIdsRequest generates requests for GetAccountsIds
func NewGetAccountsIdsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/accounts")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetDevicesRequest generates requests for GetDevices
func NewGetDevicesRequest(server string, accountId AccountId) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "accountId", runtime.ParamLocationPath, accountId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/accounts/%s/devices", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetMultipleSensorsRequest generates requests for GetMultipleSensors
func NewGetMultipleSensorsRequest(server string, accountId AccountId, params *GetMultipleSensorsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "accountId", runtime.ParamLocationPath, accountId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/accounts/%s/sensors", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Sn != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "sn", runtime.ParamLocationQuery, *params.Sn); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageNumber != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "pageNumber", runtime.ParamLocationQuery, *params.PageNumber); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Unit != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "unit", runtime.ParamLocationQuery, *params.Unit); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetHealthRequest generates requests for GetHealth
func NewGetHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/health")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetAccountsIdsWithResponse request
	GetAccountsIdsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAccountsIdsResponse, error)

	// GetDevicesWithResponse request
	GetDevicesWithResponse(ctx context.Context, accountId AccountId, reqEditors ...RequestEditorFn) (*GetDevicesResponse, error)

	// GetMultipleSensorsWithResponse request
	GetMultipleSensorsWithResponse(ctx context.Context, accountId AccountId, params *GetMultipleSensorsParams, reqEditors ...RequestEditorFn) (*GetMultipleSensorsResponse, error)

	// GetHealthWithResponse request
	GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error)
}

type GetAccountsIdsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *AccountsResponse
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetAccountsIdsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAccountsIdsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetDevicesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *DevicesResponse
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetDevicesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDevicesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetMultipleSensorsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// HasNext True if next pages can be fetched, false otherwise.
		HasNext    *bool              `json:"hasNext,omitempty"`
		Results    *[]SensorsResponse `json:"results,omitempty"`
		TotalPages *int               `json:"totalPages,omitempty"`
	}
	JSON429     *Error
	JSONDefault *Error
}

// Status returns HTTPResponse.Status
func (r GetMultipleSensorsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetMultipleSensorsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetHealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetHealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAccountsIdsWithResponse request returning *GetAccountsIdsResponse
func (c *ClientWithResponses) GetAccountsIdsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAccountsIdsResponse, error) {
	rsp, err := c.GetAccountsIds(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAccountsIdsResponse(rsp)
}

// GetDevicesWithResponse request returning *GetDevicesResponse
func (c *ClientWithResponses) GetDevicesWithResponse(ctx context.Context, accountId AccountId, reqEditors ...RequestEditorFn) (*GetDevicesResponse, error) {
	rsp, err := c.GetDevices(ctx, accountId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDevicesResponse(rsp)
}

// GetMultipleSensorsWithResponse request returning *GetMultipleSensorsResponse
func (c *ClientWithResponses) GetMultipleSensorsWithResponse(ctx context.Context, accountId AccountId, params *GetMultipleSensorsParams, reqEditors ...RequestEditorFn) (*GetMultipleSensorsResponse, error) {
	rsp, err := c.GetMultipleSensors(ctx, accountId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetMultipleSensorsResponse(rsp)
}

// GetHealthWithResponse request returning *GetHealthResponse
func (c *ClientWithResponses) GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error) {
	rsp, err := c.GetHealth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHealthResponse(rsp)
}

// ParseGetAccountsIdsResponse parses an HTTP response from a GetAccountsIdsWithResponse call
func ParseGetAccountsIdsResponse(rsp *http.Response) (*GetAccountsIdsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAccountsIdsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest AccountsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetDevicesResponse parses an HTTP response from a GetDevicesWithResponse call
func ParseGetDevicesResponse(rsp *http.Response) (*GetDevicesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetDevicesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest DevicesResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetMultipleSensorsResponse parses an HTTP response from a GetMultipleSensorsWithResponse call
func ParseGetMultipleSensorsResponse(rsp *http.Response) (*GetMultipleSensorsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetMultipleSensorsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// HasNext True if next pages can be fetched, false otherwise.
			HasNext    *bool              `json:"hasNext,omitempty"`
			Results    *[]SensorsResponse `json:"results,omitempty"`
			TotalPages *int               `json:"totalPages,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetHealthResponse parses an HTTP response from a GetHealthWithResponse call
func ParseGetHealthResponse(rsp *http.Response) (*GetHealthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all accounts the current user is member of
	// (GET /v1/accounts)
	GetAccountsIds(ctx echo.Context) error
	// Get all devices connected to a user
	// (GET /v1/accounts/{accountId}/devices)
	GetDevices(ctx echo.Context, accountId AccountId) error
	// Get sensors for a set of devices
	// (GET /v1/accounts/{accountId}/sensors)
	GetMultipleSensors(ctx echo.Context, accountId AccountId, params GetMultipleSensorsParams) error
	// Get the health of the API
	// (GET /v1/health)
	GetHealth(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAccountsIds converts echo context to params.
func (w *ServerInterfaceWrapper) GetAccountsIds(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAccountsIds(ctx)
	return err
}

// GetDevices converts echo context to params.
func (w *ServerInterfaceWrapper) GetDevices(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "accountId" -------------
	var accountId AccountId

	err = runtime.BindStyledParameterWithOptions("simple", "accountId", ctx.Param("accountId"), &accountId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter accountId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDevices(ctx, accountId)
	return err
}

// GetMultipleSensors converts echo context to params.
func (w *ServerInterfaceWrapper) GetMultipleSensors(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "accountId" -------------
	var accountId AccountId

	err = runtime.BindStyledParameterWithOptions("simple", "accountId", ctx.Param("accountId"), &accountId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter accountId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetMultipleSensorsParams
	// ------------- Optional query parameter "sn" -------------

	err = runtime.BindQueryParameter("form", true, false, "sn", ctx.QueryParams(), &params.Sn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sn: %s", err))
	}

	// ------------- Optional query parameter "pageNumber" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageNumber", ctx.QueryParams(), &params.PageNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pageNumber: %s", err))
	}

	// ------------- Optional query parameter "unit" -------------

	err = runtime.BindQueryParameter("form", true, false, "unit", ctx.QueryParams(), &params.Unit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter unit: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMultipleSensors(ctx, accountId, params)
	return err
}

// GetHealth converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetHealth(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v1/accounts", wrapper.GetAccountsIds)
	router.GET(baseURL+"/v1/accounts/:accountId/devices", wrapper.GetDevices)
	router.GET(baseURL+"/v1/accounts/:accountId/sensors", wrapper.GetMultipleSensors)
	router.GET(baseURL+"/v1/health", wrapper.GetHealth)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"E7x+/eRdwexg5L/d/ScAAP//kyiDPREYAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
