package viewmodel

import "errors"

// Address represents the view model for an Address API response.
type Address struct {
	Address string `json:"address"`
}

// IsValid performs validation on the Address view model.
func (v *Address) IsValid() error {
	if v.Address == "" {
		return errors.New("address cannot be empty")
	}
	// Add more specific address validation if needed
	return nil
}