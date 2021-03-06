package apps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/luizalabs/teresa-api/models"
)

// UpdateAppAutoScaleReader is a Reader for the UpdateAppAutoScale structure.
type UpdateAppAutoScaleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *UpdateAppAutoScaleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUpdateAppAutoScaleOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewUpdateAppAutoScaleDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewUpdateAppAutoScaleOK creates a UpdateAppAutoScaleOK with default headers values
func NewUpdateAppAutoScaleOK() *UpdateAppAutoScaleOK {
	return &UpdateAppAutoScaleOK{}
}

/*UpdateAppAutoScaleOK handles this case with default header values.

Updated version of the app
*/
type UpdateAppAutoScaleOK struct {
	Payload *models.App
}

func (o *UpdateAppAutoScaleOK) Error() string {
	return fmt.Sprintf("[PUT /apps/{app_name}/autoScale][%d] updateAppAutoScaleOK  %+v", 200, o.Payload)
}

func (o *UpdateAppAutoScaleOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.App)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAppAutoScaleDefault creates a UpdateAppAutoScaleDefault with default headers values
func NewUpdateAppAutoScaleDefault(code int) *UpdateAppAutoScaleDefault {
	return &UpdateAppAutoScaleDefault{
		_statusCode: code,
	}
}

/*UpdateAppAutoScaleDefault handles this case with default header values.

Unexpected error
*/
type UpdateAppAutoScaleDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the update app auto scale default response
func (o *UpdateAppAutoScaleDefault) Code() int {
	return o._statusCode
}

func (o *UpdateAppAutoScaleDefault) Error() string {
	return fmt.Sprintf("[PUT /apps/{app_name}/autoScale][%d] updateAppAutoScale default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateAppAutoScaleDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
