package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

// DomainError is foundation for all of our specific domain errors.
type DomainError struct {
	Code    ErrorCode
	Message string
}

func NewDomainError(code ErrorCode, message string) *DomainError {
	if code == "" {
		code = CodeErrDefaultCode
	}

	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// If DomainErrors are compared using errors.Is, they will be equal if they
// have the same code.
func (e DomainError) Is(err error) bool {
	var de *DomainError
	if errors.As(err, &de) {
		return e.Code == de.Code
	}

	return false
}

// Useful for when you want to quickly override a preset error's message by
// doing something like:
//
// err := ErrInternalServerError.Clone().SetMessage("new message")
func (e DomainError) Clone() *DomainError {
	return &DomainError{
		Code:    e.Code,
		Message: e.Message,
	}
}

func (e *DomainError) SetCode(code ErrorCode) *DomainError {
	e.Code = code
	return e
}

func (e *DomainError) SetMessage(message string) *DomainError {
	e.Message = message
	return e
}

func (e *DomainError) SetMessagef(format string, args ...interface{}) *DomainError {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

func (e *DomainError) Error() string {
	code := e.Code
	if e.Code == "" {
		code = CodeErrDefaultCode
	}
	return fmt.Sprintf("%s::%s", code, e.Message)
}

// HTTPStatusCode maps the DomainError's code to an HTTP status code.
// Defaults to the given `defaultCode` if there is no override set.
func (e DomainError) HTTPStatusCode(defaultCode int) int {
	switch e.Code {
	case CodeRecordNotFound, CodeApplicationNotFound, CodeUserNotFound, CodeCountryNotFound:
		return http.StatusNotFound
	case CodeUnauthorized, CodeSessionExpired:
		return http.StatusUnauthorized
	default:
		// This is parameterized because of differing use cases - AbortWithError
		// wants to default to 400, others want to default to 500.
		return defaultCode
	}
}

const (
	CodeErrDefaultCode ErrorCode = "Error"
	CodeUnauthorized   ErrorCode = "Unauthorized"
	CodeInternalError  ErrorCode = "InternalError"
	CodeDBError        ErrorCode = "DBError"

	// Generic errors (not specific to any entity).
	CodeCreditLimitIsLow          ErrorCode = "lowCreditLimit"
	CodeEntityAlreadyExists       ErrorCode = "entityAlreadyExists"
	CodeRecordNotFound            ErrorCode = "RecordNotFound"
	CodeTitleVinDoesNotMatch      ErrorCode = "titleVinDoesNotMatch"
	CodeTitleFullnameDoesNotMatch ErrorCode = "titleFullnameDoesNotMatch"
	CodeSessionExpired            ErrorCode = "ExpiredSession"

	// Address.
	CodeAddressNotFound      ErrorCode = "addressNotFound"
	CodeInvalidAddressBody   ErrorCode = "invalidAddressBody"
	CodeInvalidState         ErrorCode = "invalidState"
	CodeInvalidAddressSmarty ErrorCode = "invalidAddress"

	// Application.
	CodeInvalidApplicationTransition ErrorCode = "invalidApplicationTransition"
	CodeInvalidOfferDecision         ErrorCode = "invalidOfferDecision"
	CodeApplicationDenied            ErrorCode = "applicationDenied"

	// Appraisal.
	CodeAppraisalNotFound     ErrorCode = "appraisalNotFound"
	CodeInvalidAppraisalType  ErrorCode = "invalidAppraisalType"
	CodeAppraisalCreateFailed ErrorCode = "appraisalCreateFailed"
	CodeInvalidAppraisalBody  ErrorCode = "invalidAppraisalBody"

	// Device.
	CodeDeviceNotFound ErrorCode = "deviceNotFound"

	// Phone.
	CodeInvalidPhoneNumber   ErrorCode = "invalidPhoneNumber"
	CodeInvalidPhoneLandline ErrorCode = "invalidPhoneLandline"
	CodeInvalidOTP           ErrorCode = "invalidOTP"
	CodeExpiredOTP           ErrorCode = "expiredOTP"
	CodeMaxAttemptsOTP       ErrorCode = "maxAttemptsOTP"

	// Email
	CodeInvalidEmailFormat ErrorCode = "invalidEmailFormat"
	CodeInvalidEmailCode   ErrorCode = "invalidEmailCode"

	// User.
	CodeInvalidBodyPhone      ErrorCode = "invalidBodyPhone"
	CodeUserNotFound          ErrorCode = "userNotFound"
	CodeApplicationNotFound   ErrorCode = "applicationNotFound"
	CodeInvalidID             ErrorCode = "invalidID"
	CodeAgrementNotFound      ErrorCode = "agrementNotFound"
	CodeInvalidAgreementType  ErrorCode = "invalidAgreementType"
	CodeInvalidDeliveryMethod ErrorCode = "invalidDeliveryMethod"
	CodeInvalidSSN            ErrorCode = "invalidSSN"

	// Vehicle.
	CodeInvalidMake     ErrorCode = "invalidMake"
	CodeInvalidModel    ErrorCode = "invalidModel"
	CodeInvalidYear     ErrorCode = "invalidYear"
	CodeInvalidTrim     ErrorCode = "invalidTrim"
	CodeInvalidVin      ErrorCode = "invalidVin"
	CodeInvalidVinNonUs ErrorCode = "invalidVinNonUs"
	CodeInvalidUvc      ErrorCode = "invalidUvc"
	CodeInvalidUvNonUs  ErrorCode = "invalidUvc"
	CodeUvcDoesNotMatch ErrorCode = "uvcDoesNotMatch"

	// Cardholder.
	CodeCardholderServiceError    ErrorCode = "cardholderServiceError"
	CodePrimaryCardNotFound       ErrorCode = "primaryCardNotFound"
	CodeCardholderSSNAlreadyExist ErrorCode = "cardholderSSNAlreadyExist"

	// Personal info/ID
	CodePersonalInfoVerificationFailed  ErrorCode = "personalInfoVerificationFailed"
	CodeMustDoPersonalInfoVerification  ErrorCode = "mustDoPersonalInfoVerification"
	CodeIDTypeNotAllowedForVerification ErrorCode = "idTypeNotAllowedForVerification"
	CodeIDVerificationFailed            ErrorCode = "idVerificationFailed"
	CodeInvalidFrontImage               ErrorCode = "invalidFrontImage"
	CodeInvalidBackImage                ErrorCode = "invalidBackImage"
	CodeInvalidCountryCode              ErrorCode = "invalidCountryCode"
	CodeCountryNotFound                 ErrorCode = "countryNotFound"
	CodeImageDetectionFailed            ErrorCode = "imageDetectionFailed"

	// Webhook
	CodeInvalidPOAPDF ErrorCode = "invalidPOAPDF"
)

var (
	ErrDefault = NewDomainError(CodeErrDefaultCode, "unknown")
	// ErrInternalServerError          = NewDomainError(CodeErrDefaultCode, "internal server error")
	ErrCreditLimitIsLow = NewDomainError(
		CodeCreditLimitIsLow,
		"credit limit is too low to have any profit",
	)
	ErrEntityAlreadyExist = NewDomainError(
		CodeEntityAlreadyExists,
		"entity already exists",
	)
	ErrRecordNotFound       = NewDomainError(CodeRecordNotFound, "record not found")
	ErrAddressNotFound      = NewDomainError(CodeAddressNotFound, "Address not found")
	ErrAppraisalNotFound    = NewDomainError(CodeAppraisalNotFound, "appraisal not found")
	ErrTitleVinDoesNotMatch = NewDomainError(
		CodeTitleVinDoesNotMatch,
		"Title vin does not match pre-approved vin",
	)
	ErrTitleFullnameDoesNotMatch = NewDomainError(
		CodeTitleFullnameDoesNotMatch,
		"Title fullname does not match pre-approved name",
	)
	ErrInvalidApplicationTransition = NewDomainError(
		CodeInvalidApplicationTransition,
		"Invalid application state transition",
	)
	ErrSessionExpired = NewDomainError(CodeSessionExpired, "session expired")

	// Address Errors
	ErrInvalidAddressBody   = NewDomainError(CodeInvalidAddressBody, "Invalid address body")
	ErrInvalidState         = NewDomainError(CodeInvalidState, "Invalid state")
	ErrInvalidAddressSmarty = NewDomainError(
		CodeInvalidAddressSmarty,
		"Smarty could not validate the given address",
	)

	// Phone Errors
	ErrInvalidPhoneNumber   = NewDomainError(CodeInvalidPhoneNumber, "Invalid phone number")
	ErrInvalidPhoneLandline = NewDomainError(CodeInvalidPhoneLandline, "Invalid phone landline")
	ErrInvalidOTP           = NewDomainError(CodeInvalidOTP, "Invalid OTP")
	ErrExpiredOTP           = NewDomainError(CodeExpiredOTP, "OTP has expired")
	ErrMaxAttemptsOTP       = NewDomainError(CodeMaxAttemptsOTP, "Max attempts OTP reached")

	// Email Errors
	ErrInvalidEmailFormat = NewDomainError(CodeInvalidEmailFormat, "Invalid email format")
	ErrInvalidEmailCode   = NewDomainError(CodeInvalidEmailCode, "Invalid email verification code")

	// Vehicle Errors
	ErrInvalidMake     = NewDomainError(CodeInvalidMake, "Invalid make")
	ErrInvalidModel    = NewDomainError(CodeInvalidModel, "Invalid model")
	ErrInvalidYear     = NewDomainError(CodeInvalidYear, "Invalid year")
	ErrInvalidTrim     = NewDomainError(CodeInvalidTrim, "Invalid trim")
	ErrInvalidVin      = NewDomainError(CodeInvalidVin, "Invalid vin")
	ErrInvalidVinNonUs = NewDomainError(CodeInvalidVinNonUs, "Vin number is not supported")
	ErrInvalidUvc      = NewDomainError(CodeInvalidUvc, "Invalid uvc")
	ErrInvalidUvNonUs  = NewDomainError(CodeInvalidUvNonUs, "Uvc number is not supported")
	ErrUvcDoesNotMatch = NewDomainError(CodeUvcDoesNotMatch, " UVC does not match pre-approved uvc")
	// User Errors
	ErrInvalidBodyPhone      = NewDomainError(CodeInvalidBodyPhone, "Invalid phone number")
	ErrUserNotFound          = NewDomainError(CodeUserNotFound, "User not found")
	ErrApplicationNotFound   = NewDomainError(CodeApplicationNotFound, "Application not found")
	ErrInvalidID             = NewDomainError(CodeInvalidID, "Invalid ID")
	ErrAgrementNotFound      = NewDomainError(CodeAgrementNotFound, "Agreement not found")
	ErrInvalidAgreementType  = NewDomainError(CodeInvalidAgreementType, "Invalid agreement type")
	ErrInvalidDeliveryMethod = NewDomainError(CodeInvalidDeliveryMethod, "Invalid delivery method")
	ErrInvalidSSN            = NewDomainError(CodeInvalidSSN, "Invalid SSN")
	// Appraisal Errors
	ErrInvalidAppraisalType  = NewDomainError(CodeInvalidAppraisalType, "Invalid appraisal type")
	ErrAppraisalCreateFailed = NewDomainError(CodeAppraisalCreateFailed, "Appraisal create failed")
	ErrInvalidAppraisalBody  = NewDomainError(
		CodeInvalidAppraisalBody,
		"Bad Request - the changeset contains attributes which are invalid",
	)

	// Application errors
	ErrApplicationDenied    = NewDomainError(CodeApplicationDenied, "Application has been denied")
	ErrInvalidOfferDecision = NewDomainError(CodeInvalidOfferDecision, "Invalid offer decision")

	// Device errors
	ErrDeviceNotFound = NewDomainError(CodeDeviceNotFound, "Device not found")

	// Cardholder service errors
	ErrCardholderServiceError = NewDomainError(
		CodeCardholderServiceError,
		"Cardholder service error",
	)
	ErrPrimaryCardNotFound       = NewDomainError(CodePrimaryCardNotFound, "primary card not found")
	ErrCardholderSSNAlreadyExist = NewDomainError(
		CodeCardholderSSNAlreadyExist,
		"SSN already exist with a user",
	)

	// Personal Info/ID verification errors
	ErrPersonalInfoVerificationFailed = NewDomainError(
		CodePersonalInfoVerificationFailed,
		"Personal info verification failed",
	)
	ErrMustDoPersonalInfoVerification = NewDomainError(
		CodeMustDoPersonalInfoVerification,
		"The user must do personal info verification before ID verification",
	)

	ErrIDVerificationFailed = NewDomainError(CodeIDVerificationFailed, "ID verification failed")
	ErrInvalidFrontImage    = NewDomainError(CodeInvalidFrontImage, "Invalid front image")
	ErrInvalidBackImage     = NewDomainError(CodeInvalidBackImage, "Invalid back image")
	ErrInvalidCountryCode   = NewDomainError(CodeInvalidCountryCode, "Invalid country code")
	ErrCountryNotFound      = NewDomainError(CodeCountryNotFound, "Country not found")
	ErrImageDetectionFailed = NewDomainError(CodeImageDetectionFailed, "Image detection failed")

	// Webhook errors
	ErrWrongHandshakeKey = ErrUnauthorized("incorrect handshake key")
	ErrInvalidPOAPDF     = NewDomainError(CodeInvalidPOAPDF, "Invalid POA PDF")
)

func ErrInternalServerError(message string) *DomainError {
	return NewDomainError(CodeInternalError, message)
}

func ErrInternalServerErrorf(format string, args ...interface{}) *DomainError {
	return NewDomainError(CodeInternalError, fmt.Sprintf(format, args...))
}

func ErrUnauthorized(message string) *DomainError {
	return NewDomainError(CodeUnauthorized, message)
}

func ErrDBError(message string) *DomainError {
	return NewDomainError(CodeDBError, message)
}

func ErrIDTypeNotAllowedForVerification(idType string) *DomainError {
	return NewDomainError(
		CodeIDTypeNotAllowedForVerification,
		fmt.Sprintf("ID type '%s' is not allowed for verification", idType),
	)
}
