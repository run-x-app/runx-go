# Run X API - Go Client

Welcome to the **Run X API** Go client wrapper. This client provides an easy-to-use interface for interacting with the Run X API, allowing you to manage users, applications, sessions, and more.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Authentication](#authentication)
- [Usage](#usage)
  - [User Operations](#user-operations)
    - [Get User Information](#get-user-information)
    - [Reveal User Number](#reveal-users-number)
    - [Generate a New API Key](#generate-a-new-api-key)
    - [Get User Billing Information](#get-user-billing-information)
    - [Get User Sessions](#get-user-sessions)
    - [Delete a Session](#delete-a-session)
  - [Application Operations](#application-operations)
    - [Create New Applications](#create-new-applications)
    - [Get List of Applications](#get-list-of-applications)
    - [Get an Application by ID](#get-an-application-by-id)
    - [Update an Application](#update-an-application)
    - [Delete an Application](#delete-an-application)
    - [Enable or Disable an Application](#enable-or-disable-an-application)
    - [Restart an Application](#restart-an-application)
    - [Get Catalog Applications](#get-catalog-applications)
- [API Reference](#api-reference)
  - [Models](#models)
- [License](#license)

## Introduction

The Run X API allows you to programmatically manage users and applications on the Run X platform. With this Go client wrapper, you can easily integrate these capabilities into your Go applications.

## Installation

To install the Go client for the Run X API, use the following command:

```bash
go get github.com/run-x-app/runx-go
```

## Usage

First, import the client in your Go application:

```go
import "github.com/run-x-app/runx-go"
```

Initialize the client:

```go
client := runx.NewClient("https://api.run-x.cloud", "<your_api_key>")
```

### User Operations

#### Get User Information

Retrieve information about the authenticated user.

```go
ctx := context.Background()

resp, err := client.Me(ctx)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get user information
```

- **Endpoint**: `GET /me`
- **Function**:

  ```go
  Me(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Responses**:

  - **200 OK**: Returns user information and consumption data.
  - **404 Not Found**: User not found.
  - **500 Internal Server Error**: Server error.

#### Reveal User's Number

Reveal the authenticated user's phone number.

```go
ctx := context.Background()

resp, err := client.RevealNumber(ctx)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get the phone number
```

- **Endpoint**: `GET /me/number`
- **Function**:

  ```go
  RevealNumber(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Responses**:

  - **200 OK**: Returns the user's phone number.
  - **500 Internal Server Error**: Server error.

#### Generate a New API Key

Generate a new API key for the authenticated user.

```go
ctx := context.Background()

resp, err := client.GenerateApiKey(ctx)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get the new API key
```

- **Endpoint**: `POST /me/key/generate`
- **Function**:

  ```go
  GenerateApiKey(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Responses**:

  - **200 OK**: Returns the new API key.
  - **404 Not Found**: User not found.
  - **500 Internal Server Error**: Server error.

#### Get User Billing Information

Retrieve billing information for the authenticated user.

```go
ctx := context.Background()

resp, err := client.MeBilling(ctx)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get billing information
```

- **Endpoint**: `GET /me/billing`
- **Function**:

  ```go
  MeBilling(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Responses**:

  - **200 OK**: Returns an array of payment information.
  - **404 Not Found**: User not found.
  - **500 Internal Server Error**: Server error.

#### Get User Sessions

Retrieve active sessions for the authenticated user.

```go
ctx := context.Background()

resp, err := client.MeSession(ctx)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get session information
```

- **Endpoint**: `GET /me/session`
- **Function**:

  ```go
  MeSession(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Responses**:

  - **200 OK**: Returns an array of session information.
  - **500 Internal Server Error**: Server error.

#### Delete a Session

Delete a specific session by session ID.

```go
ctx := context.Background()
sessionId := "session-id"

resp, err := client.DeleteSession(ctx, sessionId)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Check response status for success
```

- **Endpoint**: `DELETE /me/session/{sessionId}`
- **Function**:

  ```go
  DeleteSession(ctx context.Context, sessionId string, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Parameters**:

  - `sessionId` (string): The ID of the session to delete.
- **Responses**:

  - **200 OK**: Session deleted successfully.
  - **404 Not Found**: Session not found.
  - **500 Internal Server Error**: Server error.

### Application Operations

#### Create New Applications

Create one or more new applications.

```go
ctx := context.Background()

appRequest := runx.CreateAppJSONRequestBody{
    Apps: []runx.AppRequest{
        {
            Name: "my-app",
            App:  "app-type",
            Env:  &[]string{"ENV_VAR=value"},
            Cmd:  runx.PtrString("start-command"),
            Gpu:  runx.PtrInt64(0),
            Cpu:  runx.PtrInt64(1),
            Ram:  runx.PtrInt64(1024),
            Disk: runx.PtrInt64(10240),
        },
    },
    Pack: runx.PtrString("optional-pack"),
}

resp, err := client.CreateApp(ctx, appRequest)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get application creation status
```

- **Endpoint**: `POST /app`
- **Function**:

  ```go
  CreateApp(ctx context.Context, body CreateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Request Body**:

  ```json
  {
    "apps": [
      {
        "name": "string",
        "app": "string",
        "env": ["string"],
        "cmd": "string",
        "gpu": 0,
        "cpu": 0,
        "ram": 0,
        "disk": 0
      }
    ],
    "pack": "string"  // Optional
  }
  ```
- **Responses**:

  - **200 OK**: Applications are being created.
  - **400 Bad Request**: Invalid request data.
  - **401 Unauthorized**: User does not own the app.
  - **404 Not Found**: User or app not found.
  - **409 Conflict**: No app created due to conflict.
  - **500 Internal Server Error**: Server error.

#### Get List of Applications

Retrieve a list of applications for the authenticated user.

```go
ctx := context.Background()

resp, err := client.GetApps(ctx)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get list of applications
```

- **Endpoint**: `GET /app`
- **Function**:

  ```go
  GetApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Responses**:

  - **200 OK**: Returns an array of applications.
  - **404 Not Found**: User not found.
  - **500 Internal Server Error**: Server error.

#### Get an Application by ID

Retrieve details of a specific application.

```go
ctx := context.Background()
appId := "app-id"

resp, err := client.GetApp(ctx, appId)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get application details
```

- **Endpoint**: `GET /app/{appId}`
- **Function**:

  ```go
  GetApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Parameters**:

  - `appId` (string): The ID of the application.
- **Responses**:

  - **200 OK**: Returns application details and logs.
  - **401 Unauthorized**: User does not own the app.
  - **404 Not Found**: User or app not found.
  - **500 Internal Server Error**: Server error.

#### Update an Application

Update an existing application.

```go
ctx := context.Background()
appId := "app-id"

updateAppRequest := runx.UpdateAppJSONRequestBody{
    Name: runx.PtrString("new-app-name"),
    Env:  &[]string{"ENV_VAR=new_value"},
    Cmd:  runx.PtrString("new-start-command"),
    Gpu:  runx.PtrInt64(1),
    Cpu:  runx.PtrInt64(2),
    Ram:  runx.PtrInt64(2048),
    Disk: runx.PtrInt64(20480),
}

resp, err := client.UpdateApp(ctx, appId, updateAppRequest)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Check response for success
```

- **Endpoint**: `PUT /app/{appId}`
- **Function**:

  ```go
  UpdateApp(ctx context.Context, appId string, body UpdateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Parameters**:

  - `appId` (string): The ID of the application.
- **Request Body**:

  ```json
  {
    "name": "string",
    "env": ["string"],
    "cmd": "string",
    "gpu": 0,
    "cpu": 0,
    "ram": 0,
    "disk": 0
  }
  ```
- **Responses**:

  - **200 OK**: Application updated successfully.
  - **401 Unauthorized**: User does not own the app.
  - **404 Not Found**: User or app not found.
  - **500 Internal Server Error**: Server error.

#### Delete an Application

Delete an application by ID.

```go
ctx := context.Background()
appId := "app-id"

resp, err := client.DeleteApp(ctx, appId)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Check response status for success
```

- **Endpoint**: `DELETE /app/{appId}`
- **Function**:

  ```go
  DeleteApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Parameters**:

  - `appId` (string): The ID of the application.
- **Responses**:

  - **200 OK**: Application deleted successfully.
  - **401 Unauthorized**: User does not own the app.
  - **404 Not Found**: User or app not found.
  - **500 Internal Server Error**: Server error.

#### Enable or Disable an Application

Enable or disable an application.

```go
ctx := context.Background()
appId := "app-id"
enabled := runx.EnableAppParamsEnabled("true") // Use "false" to disable

resp, err := client.EnableApp(ctx, appId, enabled)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Check response for success
```

- **Endpoint**: `PATCH /app/{appId}/enable/{enabled}`
- **Function**:

  ```go
  EnableApp(ctx context.Context, appId string, enabled EnableAppParamsEnabled, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Parameters**:

  - `appId` (string): The ID of the application.
  - `enabled` (EnableAppParamsEnabled): `"true"` or `"false"`.
- **Responses**:

  - **200 OK**: Application status updated.
  - **401 Unauthorized**: User does not own the app.
  - **404 Not Found**: User or app not found.
  - **500 Internal Server Error**: Server error.

#### Restart an Application

Restart an application.

```go
ctx := context.Background()
appId := "app-id"

resp, err := client.RestartApp(ctx, appId)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Check response for success
```

- **Endpoint**: `PATCH /app/{appId}/restart`
- **Function**:

  ```go
  RestartApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Parameters**:

  - `appId` (string): The ID of the application.
- **Responses**:

  - **200 OK**: Application restarted successfully.
  - **401 Unauthorized**: User does not own the app.
  - **404 Not Found**: User or app not found.
  - **500 Internal Server Error**: Server error.

#### Get Catalog Applications

Retrieve catalog applications available to the user.

```go
ctx := context.Background()

resp, err := client.GetCatalogApps(ctx)
if err != nil {
    // Handle error
}
defer resp.Body.Close()

// Process response to get catalog applications
```

- **Endpoint**: `GET /catalog`
- **Function**:

  ```go
  GetCatalogApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
  ```
- **Responses**:

  - **200 OK**: Returns catalog applications and related information.
  - **404 Not Found**: User not found.
  - **500 Internal Server Error**: Server error.

## API Reference

### Client Interface

The client implements the following interface:

```go
type ClientInterface interface {
    // User Authentication and Registration
    Auth(ctx context.Context, body AuthJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
    Register(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

    // User Information
    Me(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
    RevealNumber(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
    GenerateApiKey(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
    MeBilling(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
    MeSession(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
    DeleteSession(ctx context.Context, sessionId string, reqEditors ...RequestEditorFn) (*http.Response, error)

    // Application Management
    CreateApp(ctx context.Context, body CreateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
    GetApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
    GetApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)
    UpdateApp(ctx context.Context, appId string, body UpdateAppJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
    DeleteApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)
    EnableApp(ctx context.Context, appId string, enabled EnableAppParamsEnabled, reqEditors ...RequestEditorFn) (*http.Response, error)
    RestartApp(ctx context.Context, appId string, reqEditors ...RequestEditorFn) (*http.Response, error)
    GetCatalogApps(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}
```

### Models

#### AuthJSONRequestBody

```go
type AuthJSONRequestBody struct {
    Number string `json:"number"`
}
```

- **Fields**:
  - `Number` (string, required): The phone number for authentication.

#### CreateAppJSONRequestBody

```go
type CreateAppJSONRequestBody struct {
    Apps []AppRequest `json:"apps"`
    Pack *string      `json:"pack,omitempty"`
}
```

- **Fields**:
  - `Apps` ([]AppRequest, required): List of applications to create.
  - `Pack` (*string, optional): Optional pack information.

#### AppRequest

```go
type AppRequest struct {
    Name string    `json:"name"`
    App  string    `json:"app"`
    Env  *[]string `json:"env,omitempty"`
    Cmd  *string   `json:"cmd,omitempty"`
    Gpu  *int64    `json:"gpu,omitempty"`
    Cpu  *int64    `json:"cpu,omitempty"`
    Ram  *int64    `json:"ram,omitempty"`
    Disk *int64    `json:"disk,omitempty"`
}
```

- **Fields**:
  - `Name` (string, required)
  - `App` (string, required)
  - `Env` (*[]string, optional)
  - `Cmd` (*string, optional)
  - `Gpu` (*int64, optional)
  - `Cpu` (*int64, optional)
  - `Ram` (*int64, optional)
  - `Disk` (*int64, optional)

#### UpdateAppJSONRequestBody

```go
type UpdateAppJSONRequestBody struct {
    Name *string   `json:"name,omitempty"`
    Env  *[]string `json:"env,omitempty"`
    Cmd  *string   `json:"cmd,omitempty"`
    Gpu  *int64    `json:"gpu,omitempty"`
    Cpu  *int64    `json:"cpu,omitempty"`
    Ram  *int64    `json:"ram,omitempty"`
    Disk *int64    `json:"disk,omitempty"`
}
```

- **Fields**:
  - `Name` (*string, optional)
  - `Env` (*[]string, optional)
  - `Cmd` (*string, optional)
  - `Gpu` (*int64, optional)
  - `Cpu` (*int64, optional)
  - `Ram` (*int64, optional)
  - `Disk` (*int64, optional)

#### EnableAppParamsEnabled

```go
type EnableAppParamsEnabled string
```

- Possible values: `"true"`, `"false"`
