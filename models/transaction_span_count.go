// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// TransactionSpanCount transaction span count
// swagger:model transactionSpanCount
type TransactionSpanCount struct {

	// dropped
	Dropped *TransactionSpanCountDropped `json:"dropped,omitempty"`
}

// Validate validates this transaction span count
func (m *TransactionSpanCount) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDropped(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TransactionSpanCount) validateDropped(formats strfmt.Registry) error {

	if swag.IsZero(m.Dropped) { // not required
		return nil
	}

	if m.Dropped != nil {

		if err := m.Dropped.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("dropped")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TransactionSpanCount) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TransactionSpanCount) UnmarshalBinary(b []byte) error {
	var res TransactionSpanCount
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
