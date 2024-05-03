package repository

import (
	"encoding/json"
	"os"
)

const accessTokenFile = "access_tokens.json"

func SaveToken(email string, token string) error {
	file, err := os.OpenFile(accessTokenFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var tokenInfos []AccessTokenInfo
	err = json.NewDecoder(file).Decode(&tokenInfos)
	if err != nil && err.Error() != "EOF" {
		return err
	}

	// check if id is already present
	idExists := false
	for _, tokenData := range tokenInfos {
		if tokenData.Email == email {
			idExists = true
			tokenData.AccessToken = token
		}
	}

	if !idExists {
		tokenInfos = append(tokenInfos, AccessTokenInfo{Email: email, AccessToken: token})
	}

	file.Seek(0, 0)
	err = json.NewEncoder(file).Encode(tokenInfos)
	if err != nil {
		return err
	}

	return nil
}

func GetToken(email string) (token string, err error) {
	file, err := os.Open(accessTokenFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var tokenInfos []AccessTokenInfo
	err = json.NewDecoder(file).Decode(&tokenInfos)
	if err != nil {
		return "", err
	}

	for _, tokenData := range tokenInfos {
		if tokenData.Email == email {
			return tokenData.AccessToken, nil
		}
	}

	return "", &TokenNotFoundError{}
}
