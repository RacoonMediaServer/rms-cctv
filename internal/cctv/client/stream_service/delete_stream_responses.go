// Code generated by go-swagger; DO NOT EDIT.

package stream_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// DeleteStreamReader is a Reader for the DeleteStream structure.
type DeleteStreamReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteStreamReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteStreamOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDeleteStreamNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteStreamInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[DELETE /stream/{id}] deleteStream", response, response.Code())
	}
}

// NewDeleteStreamOK creates a DeleteStreamOK with default headers values
func NewDeleteStreamOK() *DeleteStreamOK {
	return &DeleteStreamOK{}
}

/*
DeleteStreamOK describes a response with status code 200, with default header values.

OK
*/
type DeleteStreamOK struct {
}

// IsSuccess returns true when this delete stream o k response has a 2xx status code
func (o *DeleteStreamOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete stream o k response has a 3xx status code
func (o *DeleteStreamOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete stream o k response has a 4xx status code
func (o *DeleteStreamOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete stream o k response has a 5xx status code
func (o *DeleteStreamOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete stream o k response a status code equal to that given
func (o *DeleteStreamOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete stream o k response
func (o *DeleteStreamOK) Code() int {
	return 200
}

func (o *DeleteStreamOK) Error() string {
	return fmt.Sprintf("[DELETE /stream/{id}][%d] deleteStreamOK ", 200)
}

func (o *DeleteStreamOK) String() string {
	return fmt.Sprintf("[DELETE /stream/{id}][%d] deleteStreamOK ", 200)
}

func (o *DeleteStreamOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteStreamNotFound creates a DeleteStreamNotFound with default headers values
func NewDeleteStreamNotFound() *DeleteStreamNotFound {
	return &DeleteStreamNotFound{}
}

/*
DeleteStreamNotFound describes a response with status code 404, with default header values.

Not Found
*/
type DeleteStreamNotFound struct {
}

// IsSuccess returns true when this delete stream not found response has a 2xx status code
func (o *DeleteStreamNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete stream not found response has a 3xx status code
func (o *DeleteStreamNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete stream not found response has a 4xx status code
func (o *DeleteStreamNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete stream not found response has a 5xx status code
func (o *DeleteStreamNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete stream not found response a status code equal to that given
func (o *DeleteStreamNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete stream not found response
func (o *DeleteStreamNotFound) Code() int {
	return 404
}

func (o *DeleteStreamNotFound) Error() string {
	return fmt.Sprintf("[DELETE /stream/{id}][%d] deleteStreamNotFound ", 404)
}

func (o *DeleteStreamNotFound) String() string {
	return fmt.Sprintf("[DELETE /stream/{id}][%d] deleteStreamNotFound ", 404)
}

func (o *DeleteStreamNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteStreamInternalServerError creates a DeleteStreamInternalServerError with default headers values
func NewDeleteStreamInternalServerError() *DeleteStreamInternalServerError {
	return &DeleteStreamInternalServerError{}
}

/*
DeleteStreamInternalServerError describes a response with status code 500, with default header values.

Server Error
*/
type DeleteStreamInternalServerError struct {
}

// IsSuccess returns true when this delete stream internal server error response has a 2xx status code
func (o *DeleteStreamInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete stream internal server error response has a 3xx status code
func (o *DeleteStreamInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete stream internal server error response has a 4xx status code
func (o *DeleteStreamInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete stream internal server error response has a 5xx status code
func (o *DeleteStreamInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this delete stream internal server error response a status code equal to that given
func (o *DeleteStreamInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the delete stream internal server error response
func (o *DeleteStreamInternalServerError) Code() int {
	return 500
}

func (o *DeleteStreamInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /stream/{id}][%d] deleteStreamInternalServerError ", 500)
}

func (o *DeleteStreamInternalServerError) String() string {
	return fmt.Sprintf("[DELETE /stream/{id}][%d] deleteStreamInternalServerError ", 500)
}

func (o *DeleteStreamInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
