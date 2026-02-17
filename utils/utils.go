package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
	"gorm.io/datatypes"
)

var iso3Format = regexp.MustCompile(`^[A-Z]{3}$`)

func IsISOAlpha3Format(s string) bool {
	return iso3Format.MatchString(s)
}

func StringToUUID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, fmt.Errorf("uuid string is empty")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid: %w", err)
	}

	return uid, nil
}

func Pointer[T any](v T) *T { return &v }

func ToString[T constraints.Ordered](v T) string {
	return fmt.Sprintf("%v", v)
}

func NormalizeToJSONBody(body []byte) []byte {
	if len(body) == 0 {
		return []byte("{}")
	}

	var js json.RawMessage
	if err := json.Unmarshal(body, &js); err == nil {
		return body
	}

	envelope := map[string]string{
		"text": string(body),
	}

	normalized, err := json.Marshal(envelope)
	if err != nil {
		log.WithError(err).Warn("Failed to normalize response body")
		return []byte(`{"text":"<failed to normalize response>"}`)
	}

	return normalized
}

func SafeString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func ToJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

func SafeBytesToDatatypesJSON(b []byte) (datatypes.JSON, error) {
	if !json.Valid(b) {
		return nil, fmt.Errorf("invalid JSON")
	}
	return datatypes.JSON(b), nil
}

func ToJSONString(m map[string]any) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ToDatatypesJSON(v any) (datatypes.JSON, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(b), nil
}

func StringToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

func FloatToInt(f float64) int {
	return int(math.Round(f))
}

func StringToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64) // Base 10, 64-bit integer
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Converted int64:", num)
	}
	return num
}

func GetTimeInSeconds(t time.Time) int {
	h, m, s := t.Clock()
	return h*3600 + m*60 + s
}

func RoundTo(value float64, places int) float64 {
	scale := math.Pow(10, float64(places))
	return math.Round(value*scale) / scale
}

func GetFileExtension(fileName string) string {
	if dotIndex := strings.LastIndex(fileName, "."); dotIndex > 0 && dotIndex < len(fileName)-1 {
		return fileName[dotIndex:]
	}
	return ""
}

func IsFutureDate(dateStr string) bool {
	if date, err := time.Parse("2006-01-02", dateStr); err == nil {
		return date.After(time.Now())
	}
	return false
}

func IsDateBefore(date1, date2 string) bool {
	if d1, err1 := time.Parse("2006-01-02", date1); err1 == nil {
		if d2, err2 := time.Parse("2006-01-02", date2); err2 == nil {
			return d1.Before(d2)
		}
	}
	return false
}

func TimePtrToString(t *time.Time) string {
	if t == nil {
		return "" // or return some default like "N/A"
	}
	return t.Format("2006-01-02")
}

func DetectMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream" // fallback for unknown files
	}
}

func BytesToDataURI(file *[]byte, filename string) string {
	if file == nil {
		return ""
	}
	mimeType := DetectMimeType(filename)
	b64 := base64.StdEncoding.EncodeToString(*file)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64)
}

func IsStateAllowedForBridge(state string) bool {
	state = strings.ToUpper(strings.TrimSpace(state))
	prohibitedState := map[string]bool{
		"US-NY": true,
		"US-AK": true,
	}
	return !prohibitedState[state]
}

// IsCountryAllowedForBridge checks if a country is allowed for Bridge customer creation
// Based on Bridge's restricted country lists
func IsCountryAllowedForBridge(countryCode string) bool {
	if countryCode == "" {
		return false
	}

	// Convert to uppercase for consistent comparison
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))

	prohibitedCountries := map[string]bool{
		// Prohibited countries - completely blocked
		"PRK": true, // North Korea
		"SYR": true, // Syria
		"IRN": true, // Iran, Islamic Republic of
		"CUB": true, // Cuba

		// Controlled countries - restricted but may have some allowances
		"ERI": true, // Eritrea
		"SOM": true, // Somalia
		"ZWE": true, // Zimbabwe
		"NIC": true, // Nicaragua
		"MKD": true, // North Macedonia
		"SVN": true, // Slovenia
		"ALB": true, // Albania
		"MMR": true, // Myanmar
		"BGD": true, // Bangladesh
		"PAK": true, // Pakistan
		"UKR": true, // Ukraine
		"NPL": true, // Nepal
		"PSE": true, // Gaza Strip
		"CHN": true, // China
		"QAT": true, // Qatar
		"DZA": true, // Algeria
		"KEN": true, // Kenya
		"XKX": true, // Kosovo
		"SSD": true, // South Sudan
		"SDN": true, // Sudan
		"MAR": true, // Morocco
		"MLI": true, // Mali
		"YEM": true, // Yemen
		"NER": true, // Niger
	}
	return !prohibitedCountries[countryCode]
}

func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, char := range str {
		if char >= 'A' && char <= 'Z' {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteByte(byte(char + ('a' - 'A')))
		} else if char == ' ' || char == '-' {
			result.WriteByte('_')
		} else {
			result.WriteByte(byte(char))
		}
	}
	return result.String()
}

func ExtractErrorKeyValuesJSON(errMap map[string]interface{}) string {
	source, ok := errMap["source"].(map[string]interface{})
	if !ok {
		return "{}"
	}

	keyMap, ok := source["key"].(map[string]interface{})
	if !ok {
		return "{}"
	}

	// Ensure values are string-only
	result := make(map[string]string)
	for k, v := range keyMap {
		if msg, ok := v.(string); ok {
			result[k] = msg
		}
	}

	b, _ := json.Marshal(result)
	return string(b)
}

func IsValidASCII(s string) bool {
	if s == "" {
		return true
	}
	asciiRegex := regexp.MustCompile(`^[\x00-\x7F]*$`)
	return asciiRegex.MatchString(s)
}

func IsValidName(s string) bool {
	nameAllowedPattern := regexp.MustCompile(`^[A-Za-z][A-Za-z-\s]{0,69}$`)
	return nameAllowedPattern.MatchString(s)
}

func UniqueInt64Array(arr pq.Int64Array) pq.Int64Array {
	seen := make(map[int64]struct{})
	var result pq.Int64Array
	for _, id := range arr {
		if _, exists := seen[id]; !exists {
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}
	return result
}

func HideString(s string, numChars int) string {
	if len(s) <= numChars {
		return s
	}
	return strings.Repeat("*", len(s)-numChars) + s[len(s)-numChars:]
}

func ConvertCentsToDollars(cents int64) float64 {
	return float64(cents) / 100.00
}

func ConvertCentsToAmounts(cents int64) float64 {
	return float64(cents) / 100.00
}

func ConvertAmountsToCents(amounts float64) int64 {
	return int64(math.Round(amounts * 100))
}

func ConvertDollarsToCents(dollars float64) int64 {
	return int64(math.Round(dollars * 100))
}

func CalcDiffStringToFloat(a, b string) float64 {
	aFloat, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0
	}
	bFloat, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0
	}
	return aFloat - bFloat
}

func GenerateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

func ValidateSuffix(s, suffix string) error {
	if !strings.HasSuffix(s, suffix) {
		return fmt.Errorf("invalid InternalTxRefId: missing required suffix %q", suffix)
	}
	return nil
}

func RemoveSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

func CurrentTime() time.Time {
	return time.Now().UTC()
}

// DivMinorUnits divides numerator by denominator with high precision and returns the result
// in minor units (rounded to nearest integer). Both inputs are assumed to be in minor units (cents).
func DivMinorUnits(numerator, denominator int64) (int64, error) {
	const precisionMultiplier = 1_000_000 // higher multiplier for better precision

	if denominator == 0 {
		return 0, fmt.Errorf("division by zero")
	}

	// Scale numerator to maintain precision
	scaled := numerator * precisionMultiplier
	// Perform division
	result := (scaled + denominator/2) / denominator
	// Scale back to original minor units
	return (result + precisionMultiplier/2) / precisionMultiplier, nil
}

func UnixOrZero(t *time.Time) int64 {
	if t == nil || t.IsZero() {
		return 0
	}
	return t.Unix()
}
