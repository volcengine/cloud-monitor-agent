package error

import "errors"

// loggedErrors define some common errors.
type loggedErrors struct {
	NoConfigFileFound         error
	ConfigFileValidationError error
	NoMatchedFileFound        error
	AkSkStrInvalid            error
	MetaServiceError          error
	OpenConfigFileError       error
	CastTypeError             error
	DirPathError              error
}

// Errors is the set of errors that can occur
var Errors = loggedErrors{
	NoConfigFileFound:         errors.New("no config file found, please check if missing"),
	ConfigFileValidationError: errors.New("config file validation failed, please re-check the content on it"),
	NoMatchedFileFound:        errors.New("no matched file found with a pattern"),
	AkSkStrInvalid:            errors.New("aksk data is invalid"),
	MetaServiceError:          errors.New("meta service resp error"),
	OpenConfigFileError:       errors.New("open config file error"),
	CastTypeError:             errors.New("cast type error"),
	DirPathError:              errors.New("dir is not a path"),
}
