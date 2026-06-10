package service

import "errors"

var (
	ErrServerInternal    = errors.New("internal server error")
	ErrInvalidNameLength = errors.New("name must be 1 to 100 long")
	ErrDependentEntity   = errors.New("entity is dependent")

	// client service
	ErrInvalidSurnameLength = errors.New("surname must be 1 to 100 long")
	ErrNotAdult             = errors.New("client must be at least 18 years old")
	ErrInvalidGender        = errors.New("gender must be either male or female")
	ErrClientNotFound       = errors.New("client not found")

	// address service
	ErrInvalidCountryLength = errors.New("country must be 1 to 100 long")
	ErrInvalidCityLength    = errors.New("city must be 1 to 100 long")
	ErrInvalidStreetLength  = errors.New("street must be 1 to 100 long")
	ErrAddressNotFound      = errors.New("address not found")

	// supplier service
	ErrInvalidPhoneLength = errors.New("phone number must be 1 to 20 long")
	ErrSupplierNotFound   = errors.New("supplier not found")

	// product service
	ErrInvalidCategoryLength = errors.New("category name must be 1 to 50 long")
	ErrInvalidPrice          = errors.New("price must not be negative")
	ErrInvalidAvailableStock = errors.New("available stock must be positive")
	ErrInvalidAmount         = errors.New("amount must be positive")
	ErrProductNotFound       = errors.New("product not found")

	// image service
	ErrImageNotFound = errors.New("image not found")
)
