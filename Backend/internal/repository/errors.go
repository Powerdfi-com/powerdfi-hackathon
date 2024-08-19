package repository

import "errors"

var ErrRecordNotFound = errors.New("no matching record found")
var ErrDuplicateDetails = errors.New("duplicate details found")
var ErrActiveListing = errors.New("asset is already listed by user")
var ErrDisabledNotification = errors.New("notification not allowed by user")
