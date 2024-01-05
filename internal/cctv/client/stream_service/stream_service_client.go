// Code generated by go-swagger; DO NOT EDIT.

package stream_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new stream service API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for stream service API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetStreamURI(params *GetStreamURIParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetStreamURIOK, error)

	AddStream(params *AddStreamParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*AddStreamOK, error)

	DeleteStream(params *DeleteStreamParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteStreamOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
GetStreamURI gets stream URI

Get Stream Live URI
*/
func (a *Client) GetStreamURI(params *GetStreamURIParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetStreamURIOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetStreamURIParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetStreamURI",
		Method:             "GET",
		PathPattern:        "/stream/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetStreamURIReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetStreamURIOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetStreamURI: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
AddStream adds stream

Register stream source
*/
func (a *Client) AddStream(params *AddStreamParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*AddStreamOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewAddStreamParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "addStream",
		Method:             "POST",
		PathPattern:        "/stream",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &AddStreamReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*AddStreamOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for addStream: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
DeleteStream deletes stream

Unregister media source
*/
func (a *Client) DeleteStream(params *DeleteStreamParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteStreamOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteStreamParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "deleteStream",
		Method:             "DELETE",
		PathPattern:        "/stream/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteStreamReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteStreamOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for deleteStream: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}