package service

import "errors"

var (
	ErrServerInternal    = errors.New("internal server error")
	ErrInvalidNameLength = errors.New("name must be 1 to 100 long")

	// client service
	ErrInvalidSurnameLength = errors.New("surname must be 1 to 100 long")
	ErrNotAdult             = errors.New("client must be at least 18 years old")
	ErrInvalidGender        = errors.New("gender must be either male or female")
	ErrUserNotFound         = errors.New("user not found")

	// address service
	ErrInvalidCountryLength = errors.New("country must be 1 to 100 long")
	ErrInvalidCityLength    = errors.New("city must be 1 to 100 long")
	ErrInvalidStreetLength  = errors.New("street must be 1 to 100 long")

	// supplier service
	ErrInvalidPhoneLength = errors.New("phone number must be 1 to 20 long")

	// product service
	ErrInvalidCategoryLength = errors.New("category name must be 1 to 50 long")
	ErrInvalidPrice          = errors.New("price must not be negative")
	ErrInvalidAvailableStock = errors.New("available stock must be positive")
	ErrSupplierNotFound      = errors.New("supplier with given id not found")
	ErrProductNotFound       = errors.New("product not found")
)
