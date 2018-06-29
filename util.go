package storage

import (
	"encoding/base64"
	"encoding/json"
)

// ParseBase64GoogleCredential :
func ParseBase64GoogleCredential(data string) (*GoogleClient, error) {
	byteData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	googleClient := new(GoogleClient)
	if err := json.Unmarshal(byteData, googleClient); err != nil {
		return nil, err
	}

	return googleClient, nil
}
