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

// NewDeleteStreamParams creates a new DeleteStreamParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteStreamParams() *DeleteStreamParams {
	return &DeleteStreamParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteStreamParamsWithTimeout creates a new DeleteStreamParams object
// with the ability to set a timeout on a request.
func NewDeleteStreamParamsWithTimeout(timeout time.Duration) *DeleteStreamParams {
	return &DeleteStreamParams{
		timeout: timeout,
	}
}

// NewDeleteStreamParamsWithContext creates a new DeleteStreamParams object
// with the ability to set a context for a request.
func NewDeleteStreamParamsWithContext(ctx context.Context) *DeleteStreamParams {
	return &DeleteStreamParams{
		Context: ctx,
	}
}

// NewDeleteStreamParamsWithHTTPClient creates a new DeleteStreamParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteStreamParamsWithHTTPClient(client *http.Client) *DeleteStreamParams {
	return &DeleteStreamParams{
		HTTPClient: client,
	}
}

/*
DeleteStreamParams contains all the parameters to send to the API endpoint

	for the delete stream operation.

	Typically these are written to a http.Request.
*/
type DeleteStreamParams struct {

	/* ID.

	   Stream ID
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete stream params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteStreamParams) WithDefaults() *DeleteStreamParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete stream params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteStreamParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete stream params
func (o *DeleteStreamParams) WithTimeout(timeout time.Duration) *DeleteStreamParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete stream params
func (o *DeleteStreamParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete stream params
func (o *DeleteStreamParams) WithContext(ctx context.Context) *DeleteStreamParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete stream params
func (o *DeleteStreamParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete stream params
func (o *DeleteStreamParams) WithHTTPClient(client *http.Client) *DeleteStreamParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete stream params
func (o *DeleteStreamParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the delete stream params
func (o *DeleteStreamParams) WithID(id string) *DeleteStreamParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete stream params
func (o *DeleteStreamParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteStreamParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
