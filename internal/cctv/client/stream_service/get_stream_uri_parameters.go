// Code generated by go-swagger; DO NOT EDIT.

package stream_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetStreamURIParams creates a new GetStreamURIParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetStreamURIParams() *GetStreamURIParams {
	return &GetStreamURIParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetStreamURIParamsWithTimeout creates a new GetStreamURIParams object
// with the ability to set a timeout on a request.
func NewGetStreamURIParamsWithTimeout(timeout time.Duration) *GetStreamURIParams {
	return &GetStreamURIParams{
		timeout: timeout,
	}
}

// NewGetStreamURIParamsWithContext creates a new GetStreamURIParams object
// with the ability to set a context for a request.
func NewGetStreamURIParamsWithContext(ctx context.Context) *GetStreamURIParams {
	return &GetStreamURIParams{
		Context: ctx,
	}
}

// NewGetStreamURIParamsWithHTTPClient creates a new GetStreamURIParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetStreamURIParamsWithHTTPClient(client *http.Client) *GetStreamURIParams {
	return &GetStreamURIParams{
		HTTPClient: client,
	}
}

/*
GetStreamURIParams contains all the parameters to send to the API endpoint

	for the get stream URI operation.

	Typically these are written to a http.Request.
*/
type GetStreamURIParams struct {

	/* ID.

	   Stream ID
	*/
	ID string

	/* Transport.

	   Video transport

	   Default: "RTSP"
	*/
	Transport *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get stream URI params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetStreamURIParams) WithDefaults() *GetStreamURIParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get stream URI params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetStreamURIParams) SetDefaults() {
	var (
		transportDefault = string("RTSP")
	)

	val := GetStreamURIParams{
		Transport: &transportDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get stream URI params
func (o *GetStreamURIParams) WithTimeout(timeout time.Duration) *GetStreamURIParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get stream URI params
func (o *GetStreamURIParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get stream URI params
func (o *GetStreamURIParams) WithContext(ctx context.Context) *GetStreamURIParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get stream URI params
func (o *GetStreamURIParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get stream URI params
func (o *GetStreamURIParams) WithHTTPClient(client *http.Client) *GetStreamURIParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get stream URI params
func (o *GetStreamURIParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get stream URI params
func (o *GetStreamURIParams) WithID(id string) *GetStreamURIParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get stream URI params
func (o *GetStreamURIParams) SetID(id string) {
	o.ID = id
}

// WithTransport adds the transport to the get stream URI params
func (o *GetStreamURIParams) WithTransport(transport *string) *GetStreamURIParams {
	o.SetTransport(transport)
	return o
}

// SetTransport adds the transport to the get stream URI params
func (o *GetStreamURIParams) SetTransport(transport *string) {
	o.Transport = transport
}

// WriteToRequest writes these params to a swagger request
func (o *GetStreamURIParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if o.Transport != nil {

		// query param transport
		var qrTransport string

		if o.Transport != nil {
			qrTransport = *o.Transport
		}
		qTransport := qrTransport
		if qTransport != "" {

			if err := r.SetQueryParam("transport", qTransport); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
