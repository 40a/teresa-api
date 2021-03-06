package apps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new apps API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for apps API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetApps gets a list of apps

Get a list of apps
*/
func (a *Client) GetApps(params *GetAppsParams, authInfo runtime.ClientAuthInfoWriter) (*GetAppsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetAppsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getApps",
		Method:             "GET",
		PathPattern:        "/apps",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetAppsReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetAppsOK), nil
}

/*
PartialUpdateApp partials update app

Update some app properties, for now, only accept envvars
*/
func (a *Client) PartialUpdateApp(params *PartialUpdateAppParams, authInfo runtime.ClientAuthInfoWriter) (*PartialUpdateAppOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPartialUpdateAppParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "partialUpdateApp",
		Method:             "PATCH",
		PathPattern:        "/apps",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PartialUpdateAppReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PartialUpdateAppOK), nil
}

/*
UpdateApp updates an app

Update app properties, such as number of replicas and other things.
*/
func (a *Client) UpdateApp(params *UpdateAppParams, authInfo runtime.ClientAuthInfoWriter) (*UpdateAppOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateAppParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "updateApp",
		Method:             "PUT",
		PathPattern:        "/apps",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateAppReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UpdateAppOK), nil
}

/*
UpdateAppAutoScale updates app auto scale

Update app auto scale
*/
func (a *Client) UpdateAppAutoScale(params *UpdateAppAutoScaleParams, authInfo runtime.ClientAuthInfoWriter) (*UpdateAppAutoScaleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateAppAutoScaleParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "updateAppAutoScale",
		Method:             "PUT",
		PathPattern:        "/apps/{app_name}/autoScale",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateAppAutoScaleReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UpdateAppAutoScaleOK), nil
}

/*
UpdateAppScale updates app scale

Update app scale size (replicas)
*/
func (a *Client) UpdateAppScale(params *UpdateAppScaleParams, authInfo runtime.ClientAuthInfoWriter) (*UpdateAppScaleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateAppScaleParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "updateAppScale",
		Method:             "PUT",
		PathPattern:        "/apps/{app_name}/scale",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateAppScaleReader{formats: a.formats},
		AuthInfo:           authInfo,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UpdateAppScaleOK), nil
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
