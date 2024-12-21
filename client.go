// Package runx provides primitives to interact with the openapi HTTP API.
package runx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oapi-codegen/runtime"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for EnableAppParamsEnabled.
const (
	False EnableAppParamsEnabled = "false"
	True  EnableAppParamsEnabled = "true"
)

// App defines model for App.
type App struct {
	App       *string    `json:"app,omitempty"`
	Cmd       *string    `json:"cmd,omitempty"`
	Cpu       *int       `json:"cpu,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Disk      *int       `json:"disk,omitempty"`
	Env       *[]string  `json:"env,omitempty"`
	Gpu       *int       `json:"gpu,omitempty"`
	Host      *string    `json:"host,omitempty"`
	Id        *string    `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Paths     *[]string  `json:"paths,omitempty"`
	Price     *float32   `json:"price,omitempty"`
	Ram       *int       `json:"ram,omitempty"`
	ShortId   *string    `json:"short_id,omitempty"`
	User      *string    `json:"user,omitempty"`
}

// AppExtended defines model for AppExtended.
type AppExtended struct {
	App        *string                 `json:"app,omitempty"`
	Cpu        *int                    `json:"cpu,omitempty"`
	CreatedAt  *time.Time              `json:"created_at,omitempty"`
	Disk       *int                    `json:"disk,omitempty"`
	Enabled    *bool                   `json:"enabled,omitempty"`
	Env        *[]string               `json:"env,omitempty"`
	Gpu        *int                    `json:"gpu,omitempty"`
	Host       *string                 `json:"host,omitempty"`
	Id         *string                 `json:"id,omitempty"`
	Monitoring *map[string]interface{} `json:"monitoring,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Paths      *[]string               `json:"paths,omitempty"`
	Ram        *int                    `json:"ram,omitempty"`
	ShortId    *string                 `json:"short_id,omitempty"`
	Status     *string                 `json:"status,omitempty"`
	UpdatedAt  *time.Time              `json:"updated_at,omitempty"`
}

// AppRequest defines model for AppRequest.
type AppRequest struct {
	App  string    `json:"app"`
	Cmd  *string   `json:"cmd,omitempty"`
	Cpu  *int      `json:"cpu,omitempty"`
	Disk *int      `json:"disk,omitempty"`
	Env  *[]string `json:"env,omitempty"`
	Gpu  *int      `json:"gpu,omitempty"`
	Name string    `json:"name"`
	Ram  *int      `json:"ram,omitempty"`
}

// AuthRequest defines model for AuthRequest.
type AuthRequest struct {
	// Number The phone number for authentication
	Number string `json:"number"`
}

// CatalogApp defines model for CatalogApp.
type CatalogApp struct {
	Enabled *bool     `json:"enabled,omitempty"`
	Gpu     *int      `json:"gpu,omitempty"`
	Id      *string   `json:"id,omitempty"`
	Level   *int      `json:"level,omitempty"`
	Name    *string   `json:"name,omitempty"`
	Paths   *[]string `json:"paths,omitempty"`
	Price   *float32  `json:"price,omitempty"`
}

// Consumption defines model for Consumption.
type Consumption struct {
	Date  *time.Time `json:"date,omitempty"`
	Limit *float32   `json:"limit,omitempty"`
	Value *float32   `json:"value,omitempty"`
}

// CreateAppRequest defines model for CreateAppRequest.
type CreateAppRequest struct {
	Apps []AppRequest `json:"apps"`

	// Pack Optional pack information
	Pack *string `json:"pack,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Error Error message
	Error *string `json:"error,omitempty"`
}

// FilteredUser defines model for FilteredUser.
type FilteredUser struct {
	ApiKey          *string    `json:"api_key,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	Credit          *float32   `json:"credit,omitempty"`
	Email           *string    `json:"email,omitempty"`
	FirstConnection *bool      `json:"first_connection,omitempty"`
	Id              *string    `json:"id,omitempty"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
	Level           *int       `json:"level,omitempty"`
	Limit           *float32   `json:"limit,omitempty"`
	ServicesHealth  *string    `json:"services_health,omitempty"`
	TotalApps       *int       `json:"total_apps,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// Pack defines model for Pack.
type Pack struct {
	Apps    *[]string `json:"apps,omitempty"`
	Enabled *bool     `json:"enabled,omitempty"`
	Id      *string   `json:"id,omitempty"`
	Level   *int      `json:"level,omitempty"`
	Name    *string   `json:"name,omitempty"`
}

// Payment defines model for Payment.
type Payment struct {
	Amount    *float32   `json:"amount,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Id        *string    `json:"id,omitempty"`
	User      *string    `json:"user,omitempty"`
}

// SessionInfo defines model for SessionInfo.
type SessionInfo struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Id        *string    `json:"id,omitempty"`
	Ip        *string    `json:"ip,omitempty"`
	User      *string    `json:"user,omitempty"`
	UserAgent *string    `json:"user_agent,omitempty"`
}

// UpdateAppRequest defines model for UpdateAppRequest.
type UpdateAppRequest struct {
	Cmd  *string   `json:"cmd,omitempty"`
	Cpu  *int      `json:"cpu,omitempty"`
	Disk *int      `json:"disk,omitempty"`
	Env  *[]string `json:"env,omitempty"`
	Gpu  *int      `json:"gpu,omitempty"`
	Name *string   `json:"name,omitempty"`
	Ram  *int      `json:"ram,omitempty"`
}

// EnableAppParamsEnabled defines parameters for EnableApp.
type EnableAppParamsEnabled string

// CreateAppJSONRequestBody defines body for CreateApp for application/json ContentType.
type CreateAppJSONRequestBody = CreateAppRequest

// UpdateAppJSONRequestBody defines body for UpdateApp for application/json ContentType.
type UpdateAppJSONRequestBody = UpdateAppRequest

// AuthJSONRequestBody defines body for Auth for application/json ContentType.
type AuthJSONRequestBody = AuthRequest

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
func NewClient(server string, bearer string, opts ...ClientOption) (*Client, error) {

	bearerAuth, err := securityprovider.NewSecurityProviderBearerToken(bearer)

	if err != nil {
		return nil, err
	}

	opts = append(opts, WithRequestEditorFn(bearerAuth.Intercept))

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
	// GetApps request
	GetApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateAppWithBody request with any body
	CreateAppWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateApp(ctx context.Context, body CreateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteApp request
	DeleteApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetApp request
	GetApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateAppWithBody request with any body
	UpdateAppWithBody(ctx context.Context, appId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateApp(ctx context.Context, appId string, body UpdateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// EnableApp request
	EnableApp(ctx context.Context, appId string, enabled EnableAppParamsEnabled, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RestartApp request
	RestartApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AuthWithBody request with any body
	AuthWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Auth(ctx context.Context, body AuthJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCatalogApps request
	GetCatalogApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Me request
	Me(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// MeBilling request
	MeBilling(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GenerateApiKey request
	GenerateApiKey(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RevealNumber request
	RevealNumber(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// MeSession request
	MeSession(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteSession request
	DeleteSession(ctx context.Context, sessionId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Register request
	Register(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAppsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateAppWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateAppRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateApp(ctx context.Context, body CreateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateAppRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteAppRequest(c.Server, appId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAppRequest(c.Server, appId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateAppWithBody(ctx context.Context, appId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateAppRequestWithBody(c.Server, appId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateApp(ctx context.Context, appId string, body UpdateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateAppRequest(c.Server, appId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) EnableApp(ctx context.Context, appId string, enabled EnableAppParamsEnabled, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewEnableAppRequest(c.Server, appId, enabled)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RestartApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRestartAppRequest(c.Server, appId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AuthWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAuthRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Auth(ctx context.Context, body AuthJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAuthRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetCatalogApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCatalogAppsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Me(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMeRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) MeBilling(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMeBillingRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GenerateApiKey(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGenerateApiKeyRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RevealNumber(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRevealNumberRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) MeSession(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMeSessionRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteSession(ctx context.Context, sessionId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteSessionRequest(c.Server, sessionId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Register(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRegisterRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAppsRequest generates requests for GetApps
func NewGetAppsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app")
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

// NewCreateAppRequest calls the generic CreateApp builder with application/json body
func NewCreateAppRequest(server string, body CreateAppJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateAppRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateAppRequestWithBody generates requests for CreateApp with any type of body
func NewCreateAppRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteAppRequest generates requests for DeleteApp
func NewDeleteAppRequest(server string, appId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "appId", runtime.ParamLocationPath, appId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetAppRequest generates requests for GetApp
func NewGetAppRequest(server string, appId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "appId", runtime.ParamLocationPath, appId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/%s", pathParam0)
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

// NewUpdateAppRequest calls the generic UpdateApp builder with application/json body
func NewUpdateAppRequest(server string, appId string, body UpdateAppJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateAppRequestWithBody(server, appId, "application/json", bodyReader)
}

// NewUpdateAppRequestWithBody generates requests for UpdateApp with any type of body
func NewUpdateAppRequestWithBody(server string, appId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "appId", runtime.ParamLocationPath, appId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewEnableAppRequest generates requests for EnableApp
func NewEnableAppRequest(server string, appId string, enabled EnableAppParamsEnabled) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "appId", runtime.ParamLocationPath, appId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "enabled", runtime.ParamLocationPath, enabled)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/%s/enable/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewRestartAppRequest generates requests for RestartApp
func NewRestartAppRequest(server string, appId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "appId", runtime.ParamLocationPath, appId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/%s/restart", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAuthRequest calls the generic Auth builder with application/json body
func NewAuthRequest(server string, body AuthJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAuthRequestWithBody(server, "application/json", bodyReader)
}

// NewAuthRequestWithBody generates requests for Auth with any type of body
func NewAuthRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/auth")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetCatalogAppsRequest generates requests for GetCatalogApps
func NewGetCatalogAppsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/catalog")
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

// NewMeRequest generates requests for Me
func NewMeRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/me")
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

// NewMeBillingRequest generates requests for MeBilling
func NewMeBillingRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/me/billing")
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

// NewGenerateApiKeyRequest generates requests for GenerateApiKey
func NewGenerateApiKeyRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/me/key/generate")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewRevealNumberRequest generates requests for RevealNumber
func NewRevealNumberRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/me/number")
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

// NewMeSessionRequest generates requests for MeSession
func NewMeSessionRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/me/session")
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

// NewDeleteSessionRequest generates requests for DeleteSession
func NewDeleteSessionRequest(server string, sessionId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "sessionId", runtime.ParamLocationPath, sessionId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/me/session/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewRegisterRequest generates requests for Register
func NewRegisterRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/register")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
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
func NewClientWithResponses(server string, token string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, token, opts...)
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
	// GetAppsWithResponse request
	GetAppsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAppsResponse, error)

	// CreateAppWithBodyWithResponse request with any body
	CreateAppWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateAppResponse, error)

	CreateAppWithResponse(ctx context.Context, body CreateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateAppResponse, error)

	// DeleteAppWithResponse request
	DeleteAppWithResponse(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*DeleteAppResponse, error)

	// GetAppWithResponse request
	GetAppWithResponse(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*GetAppResponse, error)

	// UpdateAppWithBodyWithResponse request with any body
	UpdateAppWithBodyWithResponse(ctx context.Context, appId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateAppResponse, error)

	UpdateAppWithResponse(ctx context.Context, appId string, body UpdateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateAppResponse, error)

	// EnableAppWithResponse request
	EnableAppWithResponse(ctx context.Context, appId string, enabled EnableAppParamsEnabled, reqEditors ...RequestEditorFn) (*EnableAppResponse, error)

	// RestartAppWithResponse request
	RestartAppWithResponse(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*RestartAppResponse, error)

	// AuthWithBodyWithResponse request with any body
	AuthWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AuthResponse, error)

	AuthWithResponse(ctx context.Context, body AuthJSONRequestBody, reqEditors ...RequestEditorFn) (*AuthResponse, error)

	// GetCatalogAppsWithResponse request
	GetCatalogAppsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCatalogAppsResponse, error)

	// MeWithResponse request
	MeWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MeResponse, error)

	// MeBillingWithResponse request
	MeBillingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MeBillingResponse, error)

	// GenerateApiKeyWithResponse request
	GenerateApiKeyWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GenerateApiKeyResponse, error)

	// RevealNumberWithResponse request
	RevealNumberWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*RevealNumberResponse, error)

	// MeSessionWithResponse request
	MeSessionWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MeSessionResponse, error)

	// DeleteSessionWithResponse request
	DeleteSessionWithResponse(ctx context.Context, sessionId string, reqEditors ...RequestEditorFn) (*DeleteSessionResponse, error)

	// RegisterWithResponse request
	RegisterWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*RegisterResponse, error)
}

type GetAppsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Apps *[]AppExtended `json:"apps,omitempty"`
	}
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GetAppsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAppsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateAppResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Apps *[]App `json:"apps,omitempty"`

		// Message Status message
		Message *string `json:"message,omitempty"`
	}
	JSON400 *ErrorResponse
	JSON401 *ErrorResponse
	JSON404 *ErrorResponse
	JSON409 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r CreateAppResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateAppResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteAppResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Message *string `json:"message,omitempty"`
	}
	JSON401 *ErrorResponse
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r DeleteAppResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteAppResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAppResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		App *App `json:"app,omitempty"`

		// Log Application logs
		Log *string `json:"log,omitempty"`
	}
	JSON401 *ErrorResponse
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GetAppResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAppResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateAppResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Message *string `json:"message,omitempty"`
	}
	JSON401 *ErrorResponse
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r UpdateAppResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateAppResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type EnableAppResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Message *string `json:"message,omitempty"`
	}
	JSON401 *ErrorResponse
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r EnableAppResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r EnableAppResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RestartAppResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Message *string `json:"message,omitempty"`
	}
	JSON401 *ErrorResponse
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r RestartAppResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RestartAppResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AuthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// Session Session token
		Session *string `json:"session,omitempty"`
	}
	JSON404 *ErrorResponse
	JSON406 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r AuthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AuthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetCatalogAppsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		AvailableGpus *int            `json:"availableGpus,omitempty"`
		Catalog       *[]CatalogApp   `json:"catalog,omitempty"`
		GpuAuthorized *bool           `json:"gpuAuthorized,omitempty"`
		Limit         *map[string]int `json:"limit,omitempty"`
		Packs         *[]Pack         `json:"packs,omitempty"`
		Threshold     *map[string]int `json:"threshold,omitempty"`
	}
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GetCatalogAppsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCatalogAppsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type MeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Consumptions *[]Consumption `json:"consumptions,omitempty"`
		User         *FilteredUser  `json:"user,omitempty"`
	}
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r MeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r MeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type MeBillingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Payments *[]Payment `json:"payments,omitempty"`
	}
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r MeBillingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r MeBillingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GenerateApiKeyResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// ApiKey The new API key
		ApiKey *string `json:"api_key,omitempty"`
	}
	JSON404 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GenerateApiKeyResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GenerateApiKeyResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RevealNumberResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// Number The user's phone number
		Number *string `json:"number,omitempty"`
	}
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r RevealNumberResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RevealNumberResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type MeSessionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Sessions *[]SessionInfo `json:"sessions,omitempty"`
	}
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r MeSessionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r MeSessionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteSessionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *ErrorResponse
	JSON500      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r DeleteSessionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteSessionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RegisterResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// Session Session token
		Session *string `json:"session,omitempty"`
	}
	JSON409 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r RegisterResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RegisterResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAppsWithResponse request returning *GetAppsResponse
func (c *ClientWithResponses) GetAppsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAppsResponse, error) {
	rsp, err := c.GetApps(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAppsResponse(rsp)
}

// CreateAppWithBodyWithResponse request with arbitrary body returning *CreateAppResponse
func (c *ClientWithResponses) CreateAppWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateAppResponse, error) {
	rsp, err := c.CreateAppWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateAppResponse(rsp)
}

func (c *ClientWithResponses) CreateAppWithResponse(ctx context.Context, body CreateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateAppResponse, error) {
	rsp, err := c.CreateApp(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateAppResponse(rsp)
}

// DeleteAppWithResponse request returning *DeleteAppResponse
func (c *ClientWithResponses) DeleteAppWithResponse(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*DeleteAppResponse, error) {
	rsp, err := c.DeleteApp(ctx, appId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteAppResponse(rsp)
}

// GetAppWithResponse request returning *GetAppResponse
func (c *ClientWithResponses) GetAppWithResponse(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*GetAppResponse, error) {
	rsp, err := c.GetApp(ctx, appId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAppResponse(rsp)
}

// UpdateAppWithBodyWithResponse request with arbitrary body returning *UpdateAppResponse
func (c *ClientWithResponses) UpdateAppWithBodyWithResponse(ctx context.Context, appId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateAppResponse, error) {
	rsp, err := c.UpdateAppWithBody(ctx, appId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateAppResponse(rsp)
}

func (c *ClientWithResponses) UpdateAppWithResponse(ctx context.Context, appId string, body UpdateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateAppResponse, error) {
	rsp, err := c.UpdateApp(ctx, appId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateAppResponse(rsp)
}

// EnableAppWithResponse request returning *EnableAppResponse
func (c *ClientWithResponses) EnableAppWithResponse(ctx context.Context, appId string, enabled EnableAppParamsEnabled, reqEditors ...RequestEditorFn) (*EnableAppResponse, error) {
	rsp, err := c.EnableApp(ctx, appId, enabled, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseEnableAppResponse(rsp)
}

// RestartAppWithResponse request returning *RestartAppResponse
func (c *ClientWithResponses) RestartAppWithResponse(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*RestartAppResponse, error) {
	rsp, err := c.RestartApp(ctx, appId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRestartAppResponse(rsp)
}

// AuthWithBodyWithResponse request with arbitrary body returning *AuthResponse
func (c *ClientWithResponses) AuthWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AuthResponse, error) {
	rsp, err := c.AuthWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAuthResponse(rsp)
}

func (c *ClientWithResponses) AuthWithResponse(ctx context.Context, body AuthJSONRequestBody, reqEditors ...RequestEditorFn) (*AuthResponse, error) {
	rsp, err := c.Auth(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAuthResponse(rsp)
}

// GetCatalogAppsWithResponse request returning *GetCatalogAppsResponse
func (c *ClientWithResponses) GetCatalogAppsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCatalogAppsResponse, error) {
	rsp, err := c.GetCatalogApps(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCatalogAppsResponse(rsp)
}

// MeWithResponse request returning *MeResponse
func (c *ClientWithResponses) MeWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MeResponse, error) {
	rsp, err := c.Me(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMeResponse(rsp)
}

// MeBillingWithResponse request returning *MeBillingResponse
func (c *ClientWithResponses) MeBillingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MeBillingResponse, error) {
	rsp, err := c.MeBilling(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMeBillingResponse(rsp)
}

// GenerateApiKeyWithResponse request returning *GenerateApiKeyResponse
func (c *ClientWithResponses) GenerateApiKeyWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GenerateApiKeyResponse, error) {
	rsp, err := c.GenerateApiKey(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGenerateApiKeyResponse(rsp)
}

// RevealNumberWithResponse request returning *RevealNumberResponse
func (c *ClientWithResponses) RevealNumberWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*RevealNumberResponse, error) {
	rsp, err := c.RevealNumber(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRevealNumberResponse(rsp)
}

// MeSessionWithResponse request returning *MeSessionResponse
func (c *ClientWithResponses) MeSessionWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MeSessionResponse, error) {
	rsp, err := c.MeSession(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMeSessionResponse(rsp)
}

// DeleteSessionWithResponse request returning *DeleteSessionResponse
func (c *ClientWithResponses) DeleteSessionWithResponse(ctx context.Context, sessionId string, reqEditors ...RequestEditorFn) (*DeleteSessionResponse, error) {
	rsp, err := c.DeleteSession(ctx, sessionId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteSessionResponse(rsp)
}

// RegisterWithResponse request returning *RegisterResponse
func (c *ClientWithResponses) RegisterWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*RegisterResponse, error) {
	rsp, err := c.Register(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRegisterResponse(rsp)
}

// ParseGetAppsResponse parses an HTTP response from a GetAppsWithResponse call
func ParseGetAppsResponse(rsp *http.Response) (*GetAppsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAppsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Apps *[]AppExtended `json:"apps,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseCreateAppResponse parses an HTTP response from a CreateAppWithResponse call
func ParseCreateAppResponse(rsp *http.Response) (*CreateAppResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateAppResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Apps *[]App `json:"apps,omitempty"`

			// Message Status message
			Message *string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseDeleteAppResponse parses an HTTP response from a DeleteAppWithResponse call
func ParseDeleteAppResponse(rsp *http.Response) (*DeleteAppResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteAppResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Message *string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetAppResponse parses an HTTP response from a GetAppWithResponse call
func ParseGetAppResponse(rsp *http.Response) (*GetAppResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAppResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			App *App `json:"app,omitempty"`

			// Log Application logs
			Log *string `json:"log,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseUpdateAppResponse parses an HTTP response from a UpdateAppWithResponse call
func ParseUpdateAppResponse(rsp *http.Response) (*UpdateAppResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateAppResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Message *string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseEnableAppResponse parses an HTTP response from a EnableAppWithResponse call
func ParseEnableAppResponse(rsp *http.Response) (*EnableAppResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &EnableAppResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Message *string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseRestartAppResponse parses an HTTP response from a RestartAppWithResponse call
func ParseRestartAppResponse(rsp *http.Response) (*RestartAppResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RestartAppResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Message *string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseAuthResponse parses an HTTP response from a AuthWithResponse call
func ParseAuthResponse(rsp *http.Response) (*AuthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AuthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// Session Session token
			Session *string `json:"session,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 406:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON406 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetCatalogAppsResponse parses an HTTP response from a GetCatalogAppsWithResponse call
func ParseGetCatalogAppsResponse(rsp *http.Response) (*GetCatalogAppsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCatalogAppsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			AvailableGpus *int            `json:"availableGpus,omitempty"`
			Catalog       *[]CatalogApp   `json:"catalog,omitempty"`
			GpuAuthorized *bool           `json:"gpuAuthorized,omitempty"`
			Limit         *map[string]int `json:"limit,omitempty"`
			Packs         *[]Pack         `json:"packs,omitempty"`
			Threshold     *map[string]int `json:"threshold,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseMeResponse parses an HTTP response from a MeWithResponse call
func ParseMeResponse(rsp *http.Response) (*MeResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &MeResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Consumptions *[]Consumption `json:"consumptions,omitempty"`
			User         *FilteredUser  `json:"user,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseMeBillingResponse parses an HTTP response from a MeBillingWithResponse call
func ParseMeBillingResponse(rsp *http.Response) (*MeBillingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &MeBillingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Payments *[]Payment `json:"payments,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGenerateApiKeyResponse parses an HTTP response from a GenerateApiKeyWithResponse call
func ParseGenerateApiKeyResponse(rsp *http.Response) (*GenerateApiKeyResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GenerateApiKeyResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// ApiKey The new API key
			ApiKey *string `json:"api_key,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseRevealNumberResponse parses an HTTP response from a RevealNumberWithResponse call
func ParseRevealNumberResponse(rsp *http.Response) (*RevealNumberResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RevealNumberResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// Number The user's phone number
			Number *string `json:"number,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseMeSessionResponse parses an HTTP response from a MeSessionWithResponse call
func ParseMeSessionResponse(rsp *http.Response) (*MeSessionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &MeSessionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Sessions *[]SessionInfo `json:"sessions,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseDeleteSessionResponse parses an HTTP response from a DeleteSessionWithResponse call
func ParseDeleteSessionResponse(rsp *http.Response) (*DeleteSessionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteSessionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseRegisterResponse parses an HTTP response from a RegisterWithResponse call
func ParseRegisterResponse(rsp *http.Response) (*RegisterResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RegisterResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// Session Session token
			Session *string `json:"session,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}
