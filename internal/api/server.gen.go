// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.0 DO NOT EDIT.
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
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// DocumentPayload defines model for DocumentPayload.
type DocumentPayload struct {
	Device     *string  `json:"device,omitempty"`
	DeviceId   *string  `json:"device_id,omitempty"`
	Document   *string  `json:"document,omitempty"`
	Percentage *float32 `json:"percentage,omitempty"`
	Progress   *string  `json:"progress,omitempty"`
}

// GetSyncProgressResponse defines model for GetSyncProgressResponse.
type GetSyncProgressResponse struct {
	Device     *string  `json:"device,omitempty"`
	DeviceId   *string  `json:"device_id,omitempty"`
	Document   *string  `json:"document,omitempty"`
	Percentage *float32 `json:"percentage,omitempty"`
	Progress   *string  `json:"progress,omitempty"`
}

// Response defines model for Response.
type Response struct {
	Message  *string `json:"message,omitempty"`
	UserName *string `json:"userName,omitempty"`
}

// UserPayload defines model for UserPayload.
type UserPayload struct {
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}

// PutSyncsProgressParams defines parameters for PutSyncsProgress.
type PutSyncsProgressParams struct {
	XAuthUser string `json:"x-auth-user"`
	XAuthKey  string `json:"x-auth-key"`
}

// GetSyncsProgressDocumentHashParams defines parameters for GetSyncsProgressDocumentHash.
type GetSyncsProgressDocumentHashParams struct {
	XAuthUser string `json:"x-auth-user"`
	XAuthKey  string `json:"x-auth-key"`
}

// GetUsersAuthParams defines parameters for GetUsersAuth.
type GetUsersAuthParams struct {
	XAuthUser string `json:"x-auth-user"`
	XAuthKey  string `json:"x-auth-key"`
}

// PutSyncsProgressJSONRequestBody defines body for PutSyncsProgress for application/json ContentType.
type PutSyncsProgressJSONRequestBody = DocumentPayload

// PostUsersCreateJSONRequestBody defines body for PostUsersCreate for application/json ContentType.
type PostUsersCreateJSONRequestBody = UserPayload

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /doc)
	GetDoc(w http.ResponseWriter, r *http.Request)

	// (GET /healthcheck)
	GetHealthcheck(w http.ResponseWriter, r *http.Request)

	// (GET /openapi.yaml)
	GetOpenapiYaml(w http.ResponseWriter, r *http.Request)

	// (PUT /syncs/progress)
	PutSyncsProgress(w http.ResponseWriter, r *http.Request, params PutSyncsProgressParams)

	// (GET /syncs/progress/{documentHash})
	GetSyncsProgressDocumentHash(w http.ResponseWriter, r *http.Request, documentHash string, params GetSyncsProgressDocumentHashParams)

	// (GET /users/auth)
	GetUsersAuth(w http.ResponseWriter, r *http.Request, params GetUsersAuthParams)

	// (POST /users/create)
	PostUsersCreate(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /doc)
func (_ Unimplemented) GetDoc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /healthcheck)
func (_ Unimplemented) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /openapi.yaml)
func (_ Unimplemented) GetOpenapiYaml(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (PUT /syncs/progress)
func (_ Unimplemented) PutSyncsProgress(w http.ResponseWriter, r *http.Request, params PutSyncsProgressParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /syncs/progress/{documentHash})
func (_ Unimplemented) GetSyncsProgressDocumentHash(w http.ResponseWriter, r *http.Request, documentHash string, params GetSyncsProgressDocumentHashParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /users/auth)
func (_ Unimplemented) GetUsersAuth(w http.ResponseWriter, r *http.Request, params GetUsersAuthParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /users/create)
func (_ Unimplemented) PostUsersCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetDoc operation middleware
func (siw *ServerInterfaceWrapper) GetDoc(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDoc(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetHealthcheck operation middleware
func (siw *ServerInterfaceWrapper) GetHealthcheck(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealthcheck(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetOpenapiYaml operation middleware
func (siw *ServerInterfaceWrapper) GetOpenapiYaml(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOpenapiYaml(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PutSyncsProgress operation middleware
func (siw *ServerInterfaceWrapper) PutSyncsProgress(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PutSyncsProgressParams

	headers := r.Header

	// ------------- Required header parameter "x-auth-user" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("x-auth-user")]; found {
		var XAuthUser string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "x-auth-user", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "x-auth-user", valueList[0], &XAuthUser, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "x-auth-user", Err: err})
			return
		}

		params.XAuthUser = XAuthUser

	} else {
		err := fmt.Errorf("Header parameter x-auth-user is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "x-auth-user", Err: err})
		return
	}

	// ------------- Required header parameter "x-auth-key" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("x-auth-key")]; found {
		var XAuthKey string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "x-auth-key", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "x-auth-key", valueList[0], &XAuthKey, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "x-auth-key", Err: err})
			return
		}

		params.XAuthKey = XAuthKey

	} else {
		err := fmt.Errorf("Header parameter x-auth-key is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "x-auth-key", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PutSyncsProgress(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetSyncsProgressDocumentHash operation middleware
func (siw *ServerInterfaceWrapper) GetSyncsProgressDocumentHash(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "documentHash" -------------
	var documentHash string

	err = runtime.BindStyledParameterWithOptions("simple", "documentHash", chi.URLParam(r, "documentHash"), &documentHash, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: false})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "documentHash", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSyncsProgressDocumentHashParams

	headers := r.Header

	// ------------- Required header parameter "x-auth-user" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("x-auth-user")]; found {
		var XAuthUser string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "x-auth-user", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "x-auth-user", valueList[0], &XAuthUser, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "x-auth-user", Err: err})
			return
		}

		params.XAuthUser = XAuthUser

	} else {
		err := fmt.Errorf("Header parameter x-auth-user is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "x-auth-user", Err: err})
		return
	}

	// ------------- Required header parameter "x-auth-key" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("x-auth-key")]; found {
		var XAuthKey string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "x-auth-key", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "x-auth-key", valueList[0], &XAuthKey, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "x-auth-key", Err: err})
			return
		}

		params.XAuthKey = XAuthKey

	} else {
		err := fmt.Errorf("Header parameter x-auth-key is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "x-auth-key", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSyncsProgressDocumentHash(w, r, documentHash, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetUsersAuth operation middleware
func (siw *ServerInterfaceWrapper) GetUsersAuth(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUsersAuthParams

	headers := r.Header

	// ------------- Required header parameter "x-auth-user" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("x-auth-user")]; found {
		var XAuthUser string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "x-auth-user", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "x-auth-user", valueList[0], &XAuthUser, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "x-auth-user", Err: err})
			return
		}

		params.XAuthUser = XAuthUser

	} else {
		err := fmt.Errorf("Header parameter x-auth-user is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "x-auth-user", Err: err})
		return
	}

	// ------------- Required header parameter "x-auth-key" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("x-auth-key")]; found {
		var XAuthKey string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "x-auth-key", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "x-auth-key", valueList[0], &XAuthKey, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "x-auth-key", Err: err})
			return
		}

		params.XAuthKey = XAuthKey

	} else {
		err := fmt.Errorf("Header parameter x-auth-key is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "x-auth-key", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUsersAuth(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostUsersCreate operation middleware
func (siw *ServerInterfaceWrapper) PostUsersCreate(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostUsersCreate(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/doc", wrapper.GetDoc)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/healthcheck", wrapper.GetHealthcheck)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/openapi.yaml", wrapper.GetOpenapiYaml)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/syncs/progress", wrapper.PutSyncsProgress)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/syncs/progress/{documentHash}", wrapper.GetSyncsProgressDocumentHash)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/users/auth", wrapper.GetUsersAuth)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/users/create", wrapper.PostUsersCreate)
	})

	return r
}

type GetDocRequestObject struct {
}

type GetDocResponseObject interface {
	VisitGetDocResponse(w http.ResponseWriter) error
}

type GetDoc200TexthtmlResponse struct {
	Body          io.Reader
	ContentLength int64
}

func (response GetDoc200TexthtmlResponse) VisitGetDocResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html")
	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}
	w.WriteHeader(200)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}
	_, err := io.Copy(w, response.Body)
	return err
}

type GetDoc400JSONResponse Response

func (response GetDoc400JSONResponse) VisitGetDocResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetHealthcheckRequestObject struct {
}

type GetHealthcheckResponseObject interface {
	VisitGetHealthcheckResponse(w http.ResponseWriter) error
}

type GetHealthcheck200JSONResponse Response

func (response GetHealthcheck200JSONResponse) VisitGetHealthcheckResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetOpenapiYamlRequestObject struct {
}

type GetOpenapiYamlResponseObject interface {
	VisitGetOpenapiYamlResponse(w http.ResponseWriter) error
}

type GetOpenapiYaml200ApplicationyamlResponse struct {
	Body          io.Reader
	ContentLength int64
}

func (response GetOpenapiYaml200ApplicationyamlResponse) VisitGetOpenapiYamlResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/yaml")
	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}
	w.WriteHeader(200)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}
	_, err := io.Copy(w, response.Body)
	return err
}

type GetOpenapiYaml400JSONResponse Response

func (response GetOpenapiYaml400JSONResponse) VisitGetOpenapiYamlResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PutSyncsProgressRequestObject struct {
	Params PutSyncsProgressParams
	Body   *PutSyncsProgressJSONRequestBody
}

type PutSyncsProgressResponseObject interface {
	VisitPutSyncsProgressResponse(w http.ResponseWriter) error
}

type PutSyncsProgress200JSONResponse Response

func (response PutSyncsProgress200JSONResponse) VisitPutSyncsProgressResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PutSyncsProgress401JSONResponse Response

func (response PutSyncsProgress401JSONResponse) VisitPutSyncsProgressResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetSyncsProgressDocumentHashRequestObject struct {
	DocumentHash string `json:"documentHash,omitempty"`
	Params       GetSyncsProgressDocumentHashParams
}

type GetSyncsProgressDocumentHashResponseObject interface {
	VisitGetSyncsProgressDocumentHashResponse(w http.ResponseWriter) error
}

type GetSyncsProgressDocumentHash200JSONResponse GetSyncProgressResponse

func (response GetSyncsProgressDocumentHash200JSONResponse) VisitGetSyncsProgressDocumentHashResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetSyncsProgressDocumentHash401JSONResponse Response

func (response GetSyncsProgressDocumentHash401JSONResponse) VisitGetSyncsProgressDocumentHashResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetSyncsProgressDocumentHash502JSONResponse Response

func (response GetSyncsProgressDocumentHash502JSONResponse) VisitGetSyncsProgressDocumentHashResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(502)

	return json.NewEncoder(w).Encode(response)
}

type GetUsersAuthRequestObject struct {
	Params GetUsersAuthParams
}

type GetUsersAuthResponseObject interface {
	VisitGetUsersAuthResponse(w http.ResponseWriter) error
}

type GetUsersAuth200JSONResponse Response

func (response GetUsersAuth200JSONResponse) VisitGetUsersAuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetUsersAuth401JSONResponse Response

func (response GetUsersAuth401JSONResponse) VisitGetUsersAuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type PostUsersCreateRequestObject struct {
	Body *PostUsersCreateJSONRequestBody
}

type PostUsersCreateResponseObject interface {
	VisitPostUsersCreateResponse(w http.ResponseWriter) error
}

type PostUsersCreate201JSONResponse Response

func (response PostUsersCreate201JSONResponse) VisitPostUsersCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type PostUsersCreate402JSONResponse Response

func (response PostUsersCreate402JSONResponse) VisitPostUsersCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(402)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /doc)
	GetDoc(ctx context.Context, request GetDocRequestObject) (GetDocResponseObject, error)

	// (GET /healthcheck)
	GetHealthcheck(ctx context.Context, request GetHealthcheckRequestObject) (GetHealthcheckResponseObject, error)

	// (GET /openapi.yaml)
	GetOpenapiYaml(ctx context.Context, request GetOpenapiYamlRequestObject) (GetOpenapiYamlResponseObject, error)

	// (PUT /syncs/progress)
	PutSyncsProgress(ctx context.Context, request PutSyncsProgressRequestObject) (PutSyncsProgressResponseObject, error)

	// (GET /syncs/progress/{documentHash})
	GetSyncsProgressDocumentHash(ctx context.Context, request GetSyncsProgressDocumentHashRequestObject) (GetSyncsProgressDocumentHashResponseObject, error)

	// (GET /users/auth)
	GetUsersAuth(ctx context.Context, request GetUsersAuthRequestObject) (GetUsersAuthResponseObject, error)

	// (POST /users/create)
	PostUsersCreate(ctx context.Context, request PostUsersCreateRequestObject) (PostUsersCreateResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetDoc operation middleware
func (sh *strictHandler) GetDoc(w http.ResponseWriter, r *http.Request) {
	var request GetDocRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetDoc(ctx, request.(GetDocRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetDoc")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetDocResponseObject); ok {
		if err := validResponse.VisitGetDocResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetHealthcheck operation middleware
func (sh *strictHandler) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	var request GetHealthcheckRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetHealthcheck(ctx, request.(GetHealthcheckRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetHealthcheck")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetHealthcheckResponseObject); ok {
		if err := validResponse.VisitGetHealthcheckResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetOpenapiYaml operation middleware
func (sh *strictHandler) GetOpenapiYaml(w http.ResponseWriter, r *http.Request) {
	var request GetOpenapiYamlRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetOpenapiYaml(ctx, request.(GetOpenapiYamlRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetOpenapiYaml")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetOpenapiYamlResponseObject); ok {
		if err := validResponse.VisitGetOpenapiYamlResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PutSyncsProgress operation middleware
func (sh *strictHandler) PutSyncsProgress(w http.ResponseWriter, r *http.Request, params PutSyncsProgressParams) {
	var request PutSyncsProgressRequestObject

	request.Params = params

	var body PutSyncsProgressJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PutSyncsProgress(ctx, request.(PutSyncsProgressRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PutSyncsProgress")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PutSyncsProgressResponseObject); ok {
		if err := validResponse.VisitPutSyncsProgressResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetSyncsProgressDocumentHash operation middleware
func (sh *strictHandler) GetSyncsProgressDocumentHash(w http.ResponseWriter, r *http.Request, documentHash string, params GetSyncsProgressDocumentHashParams) {
	var request GetSyncsProgressDocumentHashRequestObject

	request.DocumentHash = documentHash
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetSyncsProgressDocumentHash(ctx, request.(GetSyncsProgressDocumentHashRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetSyncsProgressDocumentHash")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetSyncsProgressDocumentHashResponseObject); ok {
		if err := validResponse.VisitGetSyncsProgressDocumentHashResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetUsersAuth operation middleware
func (sh *strictHandler) GetUsersAuth(w http.ResponseWriter, r *http.Request, params GetUsersAuthParams) {
	var request GetUsersAuthRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetUsersAuth(ctx, request.(GetUsersAuthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetUsersAuth")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetUsersAuthResponseObject); ok {
		if err := validResponse.VisitGetUsersAuthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostUsersCreate operation middleware
func (sh *strictHandler) PostUsersCreate(w http.ResponseWriter, r *http.Request) {
	var request PostUsersCreateRequestObject

	var body PostUsersCreateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostUsersCreate(ctx, request.(PostUsersCreateRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostUsersCreate")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostUsersCreateResponseObject); ok {
		if err := validResponse.VisitPostUsersCreateResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xX32/bNhD+VwhuwDpAsRS7bjq9ZQuwFgNWY0Ufhq0YGOpssZFJjnd0IgT+3weS8s/Y",
	"zpItdYHmSQLBO353933H4y2XZmqNBk3Iy1uOsoapiL8XRvopaBqJtjGiCkvWGQuOFMQNFcyUhPAHN2Jq",
	"G+Alv1K6auDEXr/kGafWhjUkp/SEz7PO4i9VbRqdwWBQVD8MTi9fn42LwaD/anA2lAJ2uuhQPd6DBSdB",
	"k5hsQi96w+Vm7aeX4OJmZyYOEDePK3rDu57nyxVz+QkkBfOfgd63Wo46L78BWqMRnnP5iFzuT94UELch",
	"8PfgZuCYQuYt+9MXRf8Vc17r4H9HJB7B/SqmW06Eptro9t8B/IDg9orFCsRr47aqRU7ofWj0f0ETlpQe",
	"m+CAFKWMtFpeGZexX4wDUYFjYYWlRPGMz8ChMpqX/LRX9IqAw1jQwipe8kFcyrgVVMeA8srI8J1A5E+I",
	"VZAy+m3Fy8D7CyN5xl1XtGjSL4rwkUZTxzqCG8prmjar3hOXt8KLdEfplKUE8J0FfT56yzqTecZf3vEt",
	"rG2UjJjyTxjM1o/41sGYl/ybfNX/8q755Uum7Tj5XLfMa/RSAuLYN2wRIXuxYlxlrnXGxqoBVhlA/R0x",
	"uFFIGQOS3we384znNYiGalmDvDqUyDdr2+5N6NMEfUhLXTAdU3qtSMXcF827tO/3sO1B0Swcf20swVZL",
	"zNd7p/U7Ujvy8abBxVUTperEFAgc8vKPW64CsDoKn2c8dRd+cyI81Seh38Ry/O2Vg4qX5Dxka4nYDO9D",
	"156YGTOqgXXmq0617F87WtVhJFfQPgDIqOuqrBZY70Wz7L130XxMZwHSj6Zq/zdubM9OqZbH0W7o8Qv6",
	"MG8rQVCxFTWbNgnj9Gho2FioBir2QmkhSc1SBTN27YyesEXxDqkiv12MMm8E1vND/WdDJBdrVnsEE+67",
	"FUmrTYNVPh4zPt0nha9blE8klX0T+b1cdUBOweyL0U6HRzQPk0/Gh0X/i4C70BLThtjYeL0l8RAF5oF/",
	"h/QcWI/nnurnC+/o2jrEhpBBFgIGTcH7EXW0DUUZ/ZgrKPFTOhCU3qQGd81lBhNHf0obn2beWH967pw1",
	"PmNmU0Z2lbd/7Nk7Vb4Jr982zdqYMQcThZRKFsdyheKygfVqzzOOcWxPjWXzvAuYQWNsbGW4eEt71wTJ",
	"E9kyzxsjRVMbpPJ1MSz4/OP8nwAAAP//QTfb/ncTAAA=",
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
