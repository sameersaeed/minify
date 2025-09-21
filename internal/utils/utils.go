package utils

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"minify/internal/config"
)

// JSONResponse sends a JSON response with the given data and status code
func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// JSONError sends a JSON error response with the given message and status code
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// IsValidURL validates if an input string is a valid URL
func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// GetBaseURL extracts the base URL from request or config
func GetBaseURL(r *http.Request) string {
	cfg := config.Load()
	
	// if request came through a proxy, use the original host
	if host := r.Header.Get("X-Forwarded-Host"); host != "" {
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			return "https://" + host
		}
		return "http://" + host
	}
	
	// otherwise use the base url from config
	return cfg.BaseURL
}

// GetClientIP extracts client IP from request
func GetClientIP(r *http.Request) string {
	// first try x-forwarded-for header (set by proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// header can have multiple ips, take first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// fallback to x-real-ip header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// last fallback to remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// ValidateStruct checks a struct's fields based on validate tags (see models.go)
func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			
			switch {
			case rule == "required":
				if isEmptyValue(field) {
					return ValidationError{Field: fieldType.Name, Message: "is required"}
				}
			case rule == "email":
				if field.Kind() == reflect.String {
					if !isValidEmail(field.String()) {
						return ValidationError{Field: fieldType.Name, Message: "must be a valid email"}
					}
				}
			case rule == "url":
				if field.Kind() == reflect.String {
					if !IsValidURL(field.String()) {
						return ValidationError{Field: fieldType.Name, Message: "must be a valid URL"}
					}
				}
			case strings.HasPrefix(rule, "min="):
				minLen := parseIntFromRule(rule, "min=")
				if field.Kind() == reflect.String && len(field.String()) < minLen {
					return ValidationError{Field: fieldType.Name, Message: "is too short"}
				}
			case strings.HasPrefix(rule, "max="):
				maxLen := parseIntFromRule(rule, "max=")
				if field.Kind() == reflect.String && len(field.String()) > maxLen {
					return ValidationError{Field: fieldType.Name, Message: "is too long"}
				}
			}
		}
	}

	return nil
}

// ValidationError represents an error that occurs when validating an API input
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + " " + e.Message
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		return false
	}
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func parseIntFromRule(rule, prefix string) int {
	numStr := strings.TrimPrefix(rule, prefix)
	if numStr == "3" {
		return 3
	}
	if numStr == "6" {
		return 6
	}
	if numStr == "50" {
		return 50
	}
	return 0
}
