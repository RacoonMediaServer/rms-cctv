// Code generated by go-swagger; DO NOT EDIT.

package recorder_service

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

// NewAddRecordingParams creates a new AddRecordingParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewAddRecordingParams() *AddRecordingParams {
	return &AddRecordingParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewAddRecordingParamsWithTimeout creates a new AddRecordingParams object
// with the ability to set a timeout on a request.
func NewAddRecordingParamsWithTimeout(timeout time.Duration) *AddRecordingParams {
	return &AddRecordingParams{
		timeout: timeout,
	}
}

// NewAddRecordingParamsWithContext creates a new AddRecordingParams object
// with the ability to set a context for a request.
func NewAddRecordingParamsWithContext(ctx context.Context) *AddRecordingParams {
	return &AddRecordingParams{
		Context: ctx,
	}
}

// NewAddRecordingParamsWithHTTPClient creates a new AddRecordingParams object
// with the ability to set a custom HTTPClient for a request.
func NewAddRecordingParamsWithHTTPClient(client *http.Client) *AddRecordingParams {
	return &AddRecordingParams{
		HTTPClient: client,
	}
}

/*
AddRecordingParams contains all the parameters to send to the API endpoint

	for the add recording operation.

	Typically these are written to a http.Request.
*/
type AddRecordingParams struct {

	// Recording.
	Recording AddRecordingBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the add recording params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AddRecordingParams) WithDefaults() *AddRecordingParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the add recording params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AddRecordingParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the add recording params
func (o *AddRecordingParams) WithTimeout(timeout time.Duration) *AddRecordingParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the add recording params
func (o *AddRecordingParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the add recording params
func (o *AddRecordingParams) WithContext(ctx context.Context) *AddRecordingParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the add recording params
func (o *AddRecordingParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the add recording params
func (o *AddRecordingParams) WithHTTPClient(client *http.Client) *AddRecordingParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the add recording params
func (o *AddRecordingParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRecording adds the recording to the add recording params
func (o *AddRecordingParams) WithRecording(recording AddRecordingBody) *AddRecordingParams {
	o.SetRecording(recording)
	return o
}

// SetRecording adds the recording to the add recording params
func (o *AddRecordingParams) SetRecording(recording AddRecordingBody) {
	o.Recording = recording
}

// WriteToRequest writes these params to a swagger request
func (o *AddRecordingParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Recording); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
