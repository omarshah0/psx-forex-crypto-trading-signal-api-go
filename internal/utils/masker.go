package utils

import (
	"encoding/json"
	"reflect"
)

// MaskSensitiveData masks sensitive data in a map based on sensitive and ignored keys
func MaskSensitiveData(data interface{}, sensitiveKeys, ignoredKeys []string) interface{} {
	if data == nil {
		return nil
	}

	// Convert to map for processing
	var dataMap map[string]interface{}
	
	// Handle different input types
	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case string:
		// Try to unmarshal JSON string
		if err := json.Unmarshal([]byte(v), &dataMap); err != nil {
			return data
		}
	case []byte:
		// Try to unmarshal JSON bytes
		if err := json.Unmarshal(v, &dataMap); err != nil {
			return data
		}
	default:
		// Try to convert struct to map using JSON marshaling
		jsonData, err := json.Marshal(data)
		if err != nil {
			return data
		}
		if err := json.Unmarshal(jsonData, &dataMap); err != nil {
			return data
		}
	}

	// Create sets for faster lookup
	sensitiveSet := make(map[string]bool)
	ignoredSet := make(map[string]bool)
	
	for _, key := range sensitiveKeys {
		sensitiveSet[key] = true
	}
	for _, key := range ignoredKeys {
		ignoredSet[key] = true
	}

	// Mask the data
	masked := maskMap(dataMap, sensitiveSet, ignoredSet)
	return masked
}

// maskMap recursively masks sensitive data in a map
func maskMap(data map[string]interface{}, sensitiveSet, ignoredSet map[string]bool) map[string]interface{} {
	result := make(map[string]interface{})
	
	for key, value := range data {
		// Skip ignored keys
		if ignoredSet[key] {
			continue
		}
		
		// Mask sensitive keys
		if sensitiveSet[key] {
			result[key] = "****"
			continue
		}
		
		// Recursively process nested structures
		switch v := value.(type) {
		case map[string]interface{}:
			result[key] = maskMap(v, sensitiveSet, ignoredSet)
		case []interface{}:
			result[key] = maskSlice(v, sensitiveSet, ignoredSet)
		default:
			result[key] = value
		}
	}
	
	return result
}

// maskSlice recursively masks sensitive data in a slice
func maskSlice(data []interface{}, sensitiveSet, ignoredSet map[string]bool) []interface{} {
	result := make([]interface{}, len(data))
	
	for i, item := range data {
		switch v := item.(type) {
		case map[string]interface{}:
			result[i] = maskMap(v, sensitiveSet, ignoredSet)
		case []interface{}:
			result[i] = maskSlice(v, sensitiveSet, ignoredSet)
		default:
			result[i] = item
		}
	}
	
	return result
}

// IsZeroValue checks if a value is the zero value for its type
func IsZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() == 0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return val.IsNil()
	}
	
	return false
}

