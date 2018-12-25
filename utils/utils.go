package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
)

// CreateHashFromString create a md5 hash of text
func CreateHashFromString(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// IsJSONArray checks if text is a JSON Array
func IsJSONArray(text string) string {
	if string(text[0]) == "[" {
		return "{\"data\":" + text + "}"
	}

	return text
}

// ConvertToJSON combines columns & values and returns a json string
func ConvertToJSON(columns []string, values []interface{}) (string, string) {
	var id string
	entry := make(map[string]interface{})
	for i, col := range columns {
		var v interface{}
		val := values[i]
		b, ok := val.([]byte)
		if ok {
			v = string(b)
		} else {
			v = val
		}
		entry[col] = v
		if strings.ToLower(col) == "id" {
			id, ok = v.(string)
			if !ok {
				id = ""
			}
		}
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Fatalln(err)
	}

	return id, string(jsonData)
}
