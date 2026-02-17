package utils

import (
	"fmt"
	_ "image/png"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type (
	ErrorResponse map[string]interface{}
	Validation    struct {
		StatusCode            int
		Status                bool
		Response              ErrorResponse
		Message               string
		CustomValidationError bool
		CustomError           interface{}
	}
)

func NewValidationError() Validation {
	v := Validation{}
	v.Message = "Validation failed!"
	v.StatusCode = http.StatusUnprocessableEntity
	v.Status = false
	return v
}

func NewPasswordValidationError() Validation {
	v := Validation{}
	v.StatusCode = http.StatusConflict
	v.Status = false
	return v
}

func (v ErrorResponse) Add(key string, value interface{}) {
	v[key] = value
}

func (v ErrorResponse) ReplaceResponse(newResp interface{}) {
	v["errors"] = newResp
}

func (v Validation) ReplaceResponse(newResp interface{}) {
	v.Response = ErrorResponse{}
	v.Response["errors"] = newResp
}

func (e ErrorResponse) AddNested(key, field, message string) {
	nested, ok := e[key].(map[string]string)
	if !ok {
		nested = make(map[string]string)
		e[key] = nested
	}
	nested[field] = message
}

func ErrorMessage(field string) string {
	return fmt.Sprintf("Appropriate %s required.", field)
}

func ErrorMessageUri(field string) string {
	return fmt.Sprintf("Appropriate %s required uri param.", field)
}

func ErrorMessageDuplicateURI() string {
	return "Duplicate uri found."
}

func ErrorMessageASCII() string {
	return "Only english is allowed."
}

func ErrorCurrencyCode(field string) string {
	return fmt.Sprintf("%s must be a valid uppercase 3-letter currency code.", field)
}

// Additional error utility functions for validation
func ErrorInvalidLength(field string, minLen, maxLen int) string {
	if maxLen > 0 {
		return fmt.Sprintf("%s must be between %d and %d characters", field, minLen, maxLen)
	}
	return fmt.Sprintf("%s must be at least %d characters", field, minLen)
}

func ErrorInvalidChars(field string) string {
	return fmt.Sprintf("%s contains invalid characters", field)
}

func ErrorInvalidFormat(field string) string {
	return fmt.Sprintf("invalid %s format", field)
}

func ErrorInvalidDate(field string) string {
	return fmt.Sprintf("invalid %s date format, use YYYY-MM-DD", field)
}

func ErrorAgeOutOfRange(minAge, maxAge int) string {
	return fmt.Sprintf("age must be between %d and %d", minAge, maxAge)
}

func ErrorFutureDate(field string) string {
	return fmt.Sprintf("%s cannot be in the future", field)
}

func ErrorDisposableDomain(field string) string {
	return fmt.Sprintf("%s uses a disposable domain", field)
}

func ErrorInvalidCountryCode(field string) string {
	return fmt.Sprintf("%s must be exactly 3 uppercase letters", field)
}

func ErrorRestrictedCountry(field string) string {
	return fmt.Sprintf("%s is a restricted country", field)
}

func ErrorPlusRequired() string {
	return "phone number must start with '+' (international format required)"
}

// Error code constants for validation
func InvalidLengthCode() string {
	return "invalid_length"
}

func InvalidCharsCode() string {
	return "invalid_chars"
}

func InvalidFormatCode() string {
	return "invalid_format"
}

func InvalidDateCode() string {
	return "invalid_date"
}

func AgeOutOfRangeCode() string {
	return "age_out_of_range"
}

func FutureDateCode() string {
	return "future_date"
}

func DisallowedPlusCode() string {
	return "disallowed_plus"
}

func DisposableDomainCode() string {
	return "disposable_domain"
}

func InvalidCountryCodeCode() string {
	return "invalid_country_code"
}

func RestrictedCountryCode() string {
	return "restricted_country"
}

func InvalidPhoneCode() string {
	return "invalid_phone"
}

func InvalidValueCode() string {
	return "invalid_value"
}

func NonPositiveCode() string {
	return "non_positive"
}

func OutOfRangeCode() string {
	return "out_of_range"
}

func PrecisionExceededCode() string {
	return "precision_exceeded"
}

func CountryMismatchCode() string {
	return "country_mismatch"
}

func InvalidSubdivisionCode() string {
	return "invalid_subdivision"
}

func InvalidPostalCodeCode() string {
	return "invalid_postal_code"
}

func InvalidAddressPoBoxCode() string {
	return "invalid_address_po_box"
}

func InvalidOccupationId() string {
	return "invalid_occupation_id"
}

func InvalidIndustryId() string {
	return "invalid_industry_id"
}

func InvalidSourceOfFundId() string {
	return "invalid_source_of_fund_id"
}

func InvalidPurposeId() string {
	return "invalid_purpose_id"
}

func InvalidReference() string {
	return "invalid_reference"
}

func RequiredCode() string {
	return "required"
}

func NotAllowed() string {
	return "not_allowed"
}

func ErrorRequiredWithLength(minLen, maxLen int) string {
	if maxLen > 0 {
		return fmt.Sprintf("required and must be %d-%d characters", minLen, maxLen)
	}
	return fmt.Sprintf("required and must be at least %d characters", minLen)
}

func ErrorMustBeLength(minLen, maxLen int) string {
	if maxLen > 0 {
		return fmt.Sprintf("must be %d-%d characters", minLen, maxLen)
	}
	return fmt.Sprintf("must be at least %d characters", minLen)
}

func ErrorMustBeExactly(description string) string {
	return fmt.Sprintf("must be exactly %s", description)
}

func StringFiledValidation(string string, minLen int, maxLen int) bool {
	n := strings.TrimSpace(string)
	if minLen <= 0 {
		minLen = 1 // Default value for minimum length
	}
	len := len(n)
	if len < minLen {
		return false
	}
	if maxLen != -1 && len > maxLen {
		return false
	}
	return true
}

func DateFieldValidation(dateStr, layout string) bool {
	_, err := time.Parse(layout, dateStr)
	return err == nil
}

func ValidString(string string) bool {
	return len(strings.TrimSpace(string)) > 0
}

func IsValidateDateTime(dateTimeStr string) (bool, time.Time) {
	var t time.Time
	// Try to parse the input string
	t, err := time.Parse(DateTimeLayout, dateTimeStr)
	if err != nil {
		log.Error(fmt.Errorf("invalid datetime format: %v", err))
		return false, t.UTC()
	}
	return true, t.UTC()
}

func HasFile(fileName string, e echo.Context) bool {
	if _, err := e.FormFile(fileName); err == nil {
		return true
	}
	return false
}

type ImageValidation struct {
	C         echo.Context
	FiledName string
	MinMB     int
	MaxMB     int
	MaxWidth  int
	MaxHeight int
}

func StringValidation(string string) bool {
	n := strings.TrimSpace(string)
	return len(n) > 0
}

func IsEmailValid(email string) bool {
	var p string
	p = "^[a-z0-9._%-]+@[a-z0-9.-]+\\.[a-z]{2,}$"
	re := regexp.MustCompile(p)
	return re.MatchString(email)
}

func StringMaxValidation(string string, maxLen int) bool {
	n := strings.TrimSpace(string)
	if maxLen != -1 && len(n) > maxLen {
		return false
	}
	return true
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsPositiveNumber(n int) bool {
	return n > 0
}

func IsPositiveFloat(n float64) bool {
	return n > 0
}

// ValidateRequiredStringWithRegex validates a required string field with length and regex pattern
func ValidateRequiredStringWithRegex(value string, fieldName string, minLen, maxLen int, pattern *regexp.Regexp, errors map[string][]string, fieldPath string) string {
	if value == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := strings.TrimSpace(value)
	if len(normalized) < minLen || len(normalized) > maxLen {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	}

	if pattern != nil && !pattern.MatchString(normalized) {
		errors[fieldPath] = append(errors[fieldPath], InvalidCharsCode())
	}

	return normalized
}

// ValidateRequiredString validates a required string field with length constraints
func ValidateRequiredString(value string, fieldName string, minLen, maxLen int, errors map[string][]string, fieldPath string) string {
	return ValidateRequiredStringWithRegex(value, fieldName, minLen, maxLen, nil, errors, fieldPath)
}

// ValidateOptionalString validates an optional string field with length constraints
func ValidateOptionalString(value string, minLen, maxLen int, errors map[string][]string, fieldPath string) string {
	if value == "" {
		return ""
	}

	normalized := strings.TrimSpace(value)
	if len(normalized) < minLen || len(normalized) > maxLen {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	}

	return normalized
}

// ValidateEmail validates email with common rules
func ValidateEmail(email string, fieldName string, errors map[string][]string, fieldPath string) string {
	if email == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := strings.ToLower(strings.TrimSpace(email))

	if len(normalized) < 3 || len(normalized) > 254 {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	}

	if strings.Contains(normalized, "+") {
		errors[fieldPath] = append(errors[fieldPath], DisallowedPlusCode())
	}

	if !IsEmailValid(normalized) {
		errors[fieldPath] = append(errors[fieldPath], InvalidFormatCode())
	}

	// Check for disposable domains (simplified check)
	disposableDomains := []string{"@10minutemail.com", "@temp-mail.org"}
	for _, domain := range disposableDomains {
		if strings.Contains(normalized, domain) {
			errors[fieldPath] = append(errors[fieldPath], DisposableDomainCode())
			break
		}
	}

	return normalized
}

// ValidateCountryCode validates ISO country codes
func ValidateCountryCode(countryCode string, fieldName string, errors map[string][]string, fieldPath string, countryCodeRegex *regexp.Regexp) string {
	if countryCode == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := strings.ToUpper(strings.TrimSpace(countryCode))

	if len(normalized) != 3 {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	} else if countryCodeRegex != nil && !countryCodeRegex.MatchString(normalized) {
		errors[fieldPath] = append(errors[fieldPath], InvalidCountryCodeCode())
	}

	return normalized
}

// NormalizeName normalizes name strings (trim, collapse spaces, preserve casing)
func NormalizeName(name string) string {
	// Trim spaces
	name = strings.TrimSpace(name)
	// Collapse multiple spaces
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
	// Preserve original casing - names should maintain their proper capitalization
	return name
}

// ValidateRequiredName validates a required name field with normalization
func ValidateRequiredName(name string, fieldName string, minLen, maxLen int, nameRegex *regexp.Regexp, errors map[string][]string, fieldPath string) string {
	if name == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := NormalizeName(name)

	if len(normalized) < minLen || len(normalized) > maxLen {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	}

	if nameRegex != nil && !nameRegex.MatchString(normalized) {
		errors[fieldPath] = append(errors[fieldPath], InvalidCharsCode())
	}

	return normalized
}

// ValidateRequiredStreet validates a required street address field
func ValidateRequiredStreet(street string, fieldName string, minLen, maxLen int, poBoxRegex *regexp.Regexp, errors map[string][]string, fieldPath string) string {
	if street == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := strings.TrimSpace(street)

	if len(normalized) < minLen || len(normalized) > maxLen {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	}

	if poBoxRegex != nil && poBoxRegex.MatchString(normalized) {
		errors[fieldPath] = append(errors[fieldPath], InvalidAddressPoBoxCode())
	}

	return normalized
}

// NormalizeCity normalizes city strings (trim, collapse spaces, preserve casing)
func NormalizeCity(city string) string {
	// Trim spaces
	city = strings.TrimSpace(city)
	// Collapse multiple spaces
	city = regexp.MustCompile(`\s+`).ReplaceAllString(city, " ")
	// Preserve original casing - cities should maintain their proper capitalization
	return city
}

// ValidateRequiredCity validates a required city field with normalization
func ValidateRequiredCity(city string, fieldName string, minLen, maxLen int, errors map[string][]string, fieldPath string) string {
	if city == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := NormalizeCity(city)

	if len(normalized) < minLen || len(normalized) > maxLen {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	}

	return normalized
}

// ValidateRequiredPositiveInt validates a required positive integer field
func ValidateRequiredPositiveInt(value int, fieldName string, errors map[string][]string, fieldPath string) int {
	if value == 0 {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return 0
	}

	if value <= 0 {
		errors[fieldPath] = append(errors[fieldPath], InvalidValueCode())
		return 0
	}

	return value
}

// ValidateRequiredPositiveFloat validates a required positive float field with range check
func ValidateRequiredPositiveFloat(value float64, fieldName string, maxValue float64, errors map[string][]string, fieldPath string) float64 {
	if value == 0 {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return 0
	}

	if value <= 0 {
		errors[fieldPath] = append(errors[fieldPath], NonPositiveCode())
		return 0
	}

	if maxValue > 0 && value > maxValue {
		errors[fieldPath] = append(errors[fieldPath], OutOfRangeCode())
	}

	return value
}

// ValidateMonthlyVolumeUSD validates monthly volume with range check and precision validation
func ValidateMonthlyVolumeUSD(value float64, maxValue float64, errors map[string][]string, fieldPath string) float64 {
	if value == 0 {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage("monthly_volume_usd"))
		return 0
	}

	if value <= 0 {
		errors[fieldPath] = append(errors[fieldPath], NonPositiveCode())
		return 0
	}

	if value > maxValue {
		errors[fieldPath] = append(errors[fieldPath], OutOfRangeCode())
		return 0
	}

	// Round to 2 decimal places
	rounded := RoundTo(value, 2)
	if rounded != value {
		errors[fieldPath] = append(errors[fieldPath], PrecisionExceededCode())
	}

	return rounded
}

// ValidateOptionalStringWithLength validates an optional string field with length constraints
func ValidateOptionalStringWithLength(value string, minLength, maxLength int, errors map[string][]string, fieldPath string, errorCode string) string {
	if value == "" {
		return ""
	}

	trimmed := strings.TrimSpace(value)
	if len(trimmed) < minLength || len(trimmed) > maxLength {
		errors[fieldPath] = append(errors[fieldPath], errorCode)
	}

	return trimmed
}

// ValidatePhone validates and normalizes a required phone field
func ValidatePhone(phone string, region string, errors map[string][]string, fieldPath string) string {
	if phone == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage("phone"))
		return ""
	}

	trimmed := strings.TrimSpace(phone)
	validatedPhone, err := ValidateAndNormalizePhone(trimmed, region)
	if err != nil {
		errors[fieldPath] = append(errors[fieldPath], InvalidPhoneCode())
		return trimmed
	}

	return validatedPhone
}

// ValidateOptionalDateOfBirth validates an optional DOB field with age constraints
func ValidateOptionalDateOfBirth(dob string, minAge, maxAge int, errors map[string][]string, fieldPath string) string {
	if dob == "" {
		return ""
	}

	trimmed := strings.TrimSpace(dob)

	// Check format using regex
	dobRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dobRegex.MatchString(trimmed) {
		errors[fieldPath] = append(errors[fieldPath], InvalidFormatCode())
		return trimmed
	}

	// Parse and validate date
	parsedDOB, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		errors[fieldPath] = append(errors[fieldPath], InvalidDateCode())
		return trimmed
	}

	now := time.Now()

	// Check if date is in the future
	if parsedDOB.After(now) {
		errors[fieldPath] = append(errors[fieldPath], FutureDateCode())
	}

	// Calculate age and validate range
	age := now.Year() - parsedDOB.Year()
	if now.YearDay() < parsedDOB.YearDay() {
		age--
	}

	if age < minAge || age > maxAge {
		errors[fieldPath] = append(errors[fieldPath], AgeOutOfRangeCode())
	}

	return trimmed
}

// ValidateOptionalStringWithMaxLength validates an optional string field with maximum length
func ValidateOptionalStringWithMaxLength(value string, maxLength int, errors map[string][]string, fieldPath string) string {
	if value == "" {
		return ""
	}

	trimmed := strings.TrimSpace(value)
	if len(trimmed) > maxLength {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	}

	return trimmed
}

// ValidateCountryCodeWithRestrictions validates a country code and checks against restricted list
func ValidateCountryCodeWithRestrictions(countryCode string, fieldName string, restrictedCountries []string, errors map[string][]string, fieldPath string, countryCodeRegex *regexp.Regexp) string {
	if countryCode == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := strings.ToUpper(strings.TrimSpace(countryCode))

	if len(normalized) != 3 {
		errors[fieldPath] = append(errors[fieldPath], InvalidLengthCode())
	} else if countryCodeRegex != nil && !countryCodeRegex.MatchString(normalized) {
		errors[fieldPath] = append(errors[fieldPath], InvalidCountryCodeCode())
	} else if len(restrictedCountries) > 0 {
		// Check for restricted countries
		for _, country := range restrictedCountries {
			if normalized == country {
				errors[fieldPath] = append(errors[fieldPath], RestrictedCountryCode())
				break
			}
		}
	}

	return normalized
}

// ValidateTIN validates Tax Identification Number with country-specific format checking
func ValidateTIN(tin string, countryCode string, ssnRegex *regexp.Regexp, errors map[string][]string, fieldPath string) string {
	if tin == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage("tin"))
		return ""
	}

	normalized := strings.ToUpper(strings.TrimSpace(tin))

	// Country-specific validation
	if countryCode == "USA" && ssnRegex != nil {
		if !ssnRegex.MatchString(normalized) {
			errors[fieldPath] = append(errors[fieldPath], InvalidFormatCode())
		}
	}
	// Add more country-specific validations as needed

	return normalized
}

// ValidateOptionalRegexField validates an optional field against a regex pattern
func ValidateOptionalRegexField(value string, validationRegex *regexp.Regexp, errorCode string, errors map[string][]string, fieldPath string) string {
	if value == "" {
		return ""
	}

	normalized := strings.ToUpper(strings.TrimSpace(value))
	if validationRegex != nil && !validationRegex.MatchString(normalized) {
		errors[fieldPath] = append(errors[fieldPath], errorCode)
	}

	return normalized
}

// ValidateRequiredRegexField validates a required field against a regex pattern
func ValidateRequiredRegexField(value string, fieldName string, validationRegex *regexp.Regexp, errorCode string, errors map[string][]string, fieldPath string) string {
	if value == "" {
		errors[fieldPath] = append(errors[fieldPath], ErrorMessage(fieldName))
		return ""
	}

	normalized := strings.ToUpper(strings.TrimSpace(value))
	if validationRegex != nil && !validationRegex.MatchString(normalized) {
		errors[fieldPath] = append(errors[fieldPath], errorCode)
	}

	return normalized
}

func IsValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	// Trim spaces
	urlStr = strings.TrimSpace(urlStr)

	// Must start with http:// or https://
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	// Extract the part after the scheme
	var domain string
	if strings.HasPrefix(urlStr, "https://") {
		domain = strings.TrimPrefix(urlStr, "https://")
	} else {
		domain = strings.TrimPrefix(urlStr, "http://")
	}

	// Domain must not be empty
	if domain == "" {
		return false
	}

	// Extract hostname (part before path/query)
	hostname := domain
	if idx := strings.IndexAny(domain, "/?#"); idx != -1 {
		hostname = domain[:idx]
	}

	// Hostname must not be empty
	if hostname == "" {
		return false
	}

	// Remove port if present
	if idx := strings.LastIndex(hostname, ":"); idx != -1 {
		// Verify port is numeric
		port := hostname[idx+1:]
		hostname = hostname[:idx]
		for _, ch := range port {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}

	// Hostname must not contain invalid characters
	validHostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !validHostnameRegex.MatchString(hostname) {
		return false
	}

	// Must have at least one dot (TLD required) - "https://jane" is invalid
	// Examples: "example.com", "sub.example.com" are valid
	// "localhost", "jane" are invalid for production URLs
	if !strings.Contains(hostname, ".") {
		return false
	}

	// Verify TLD is at least 2 characters
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		return false
	}

	tld := parts[len(parts)-1]
	if len(tld) < 2 {
		return false
	}

	// TLD should only contain letters
	for _, ch := range tld {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
			return false
		}
	}

	return true
}

func ValidateURL(urlStr string, fieldName string) string {
	if urlStr == "" {
		return RequiredCode()
	}

	if !IsValidURL(urlStr) {
		return InvalidFormatCode()
	}

	return ""
}

func IsValidCurrencyCode(s string) bool {
	if len(s) != 3 {
		return false
	}
	for i := 0; i < 3; i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}
