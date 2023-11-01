// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MessageErr MessageErr represents a error message.
//
// swagger:model MessageErr
type MessageErr struct {

	// error
	Error string `json:"Error,omitempty"`

	// message
	Message string `json:"Message,omitempty"`

	// status
	Status int64 `json:"Status,omitempty"`
}

// Validate validates this message err
func (m *MessageErr) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this message err based on context it is used
func (m *MessageErr) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MessageErr) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MessageErr) UnmarshalBinary(b []byte) error {
	var res MessageErr
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
