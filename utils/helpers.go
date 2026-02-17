package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

const salt = "orchestration"
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

var DateTimeLayout = "2006-01-02 15:04:05-07"
var TimeLayout = "03:04 PM"

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)

	// Use crypto/rand for cryptographic-level randomness
	if _, err := rand.Read(bytes); err != nil {
		log.Error("RandomString Error", err)
		return ""
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

func HashMake(str string) (string, error) {
	str = salt + str
	return hashCreate(str)
}

func IsSameHash(str string, hashed string) bool {
	str = salt + str
	return isSameHashCheck(str, hashed)
}

func InArray(needle interface{}, haystack interface{}) bool {
	s := reflect.ValueOf(haystack)

	// Ensure haystack is an array or slice
	if s.Kind() != reflect.Array && s.Kind() != reflect.Slice {
		return false
	}

	// Iterate through haystack and compare elements
	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(needle, s.Index(i).Interface()) {
			return true
		}
	}
	return false
}
func InArrayString(target string, array []string) bool {
	for _, v := range array {
		if v == target {
			return true
		}
	}
	return false
}

func InterfaceToMap(i interface{}) (map[string]interface{}, error) {
	// If it's already a map[string]interface{}, just return it
	if m, ok := i.(map[string]interface{}); ok {
		return m, nil
	}

	// Marshal interface{} to JSON, then unmarshal into map[string]interface{}
	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal interface: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into map: %w", err)
	}

	return result, nil
}

func MapToStruct(input map[string]interface{}, output interface{}) error {
	// Marshal the map to JSON
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	// Unmarshal the JSON to the struct
	return json.Unmarshal(jsonData, output)
}

func hashCreate(str string) (string, error) {
	// Pre-hash with SHA-256
	hash := sha256.Sum256([]byte(str))

	// Create bcrypt hash
	bcryptHash, err := bcrypt.GenerateFromPassword(hash[:], bcrypt.DefaultCost)
	return string(bcryptHash), err
}

// generateIdempotencyKey generates a unique idempotency key for the request
func GenerateIdempotencyKey() string {
	// Generate a proper UUID for idempotency as recommended by Bridge API
	return uuid.New().String()
}

func isSameHashCheck(str string, hashed string) bool {
	// Pre-hash the input string with SHA-256
	hash := sha256.Sum256([]byte(str))

	// Compare the bcrypt hash with the hash of the entered password
	err := bcrypt.CompareHashAndPassword([]byte(hashed), hash[:])
	return err == nil // Returns true if the hash matches, false otherwise
}

func ImageUpload(file *multipart.FileHeader, cdnPath string, upload string) (fullPath string, err error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	path := cdnPath + upload
	fileExt := filepath.Ext(file.Filename) // e.g., .jpg, .png
	fileBase := strings.TrimSuffix(file.Filename, fileExt)

	fileName := EncodeStringBase64(GenerateRandomString(20)) + "-" + ToString(time.Now().Unix()) + "-" + EncodeParam(fileBase) + fileExt

	dstPath := fmt.Sprintf("%s/%s", path, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", upload, fileName), nil
}
func EncodeParam(s string) string {
	return url.QueryEscape(s)
}

func EncodeStringBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
func MoveFileManually(srcPath, destPath string) error {
	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy the file contents
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Close both files explicitly before deletion
	srcFile.Close()
	destFile.Close()

	// Remove the original file
	err = os.Remove(srcPath)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

func CopyFileManually(srcPath, destPath string) error {
	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy the file contents
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Close both files explicitly before deletion
	srcFile.Close()
	destFile.Close()
	return nil
}
func DeleteImage(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func FindExtraData(slice1, slice2 []string) []string {
	extra := []string{}
	existingMap := make(map[string]bool)

	// Add all elements from slice2 to a map for quick lookup
	for _, item := range slice2 {
		existingMap[item] = true
	}

	// Check elements in slice1 that are not in slice2
	for _, item := range slice1 {
		if !existingMap[item] {
			extra = append(extra, item)
		}
	}
	return extra
}

func EnsureDir(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		fmt.Printf("Directory created: %s\n", dir)
	} else if err != nil {
		// Error occurred while checking
		return fmt.Errorf("error checking directory: %w", err)
	} else if !info.IsDir() {
		// Path exists but is not a directory
		return fmt.Errorf("path exists but is not a directory: %s", dir)
	}
	// Directory exists
	return nil
}

func ParsePaginationParams(c echo.Context) (limit int, page int, offset int) {
	limitStr := c.QueryParam("per_page")
	if limitStr == "" {
		limitStr = c.QueryParam("limit")
	}
	offsetStr := c.QueryParam("offset")
	pageStr := c.QueryParam("current_page")
	if pageStr == "" {
		pageStr = c.QueryParam("page")
	}

	limit = 10 // default
	page = 1   // default
	offset = 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
			offset = (page - 1) * limit
		}
	} else if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
			page = (offset / limit) + 1
		}
	}

	return limit, page, offset
}

func CalculateTotalPages(total int, limit int) int {
	if total == 0 || limit == 0 {
		return 1
	}
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	return totalPages
}

func FloatToDecimal2PointString(val float64) string {
	return fmt.Sprintf("%.2f", math.Trunc(val*100)/100)
}
func GetStringFromMetadata(m datatypes.JSONMap, key string) (string, bool) {
	if raw, ok := m[key]; ok {
		return fmt.Sprintf("%v", raw), true
	}
	return "", false
}

func InterfaceToStruct(data interface{}, result interface{}) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := json.Unmarshal(bytesData, result); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}
