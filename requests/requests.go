package requests

import (
	"github.com/go-resty/resty/v2"
	"stress_test/serializer"
)

func Login(client *resty.Client, email string, password string) (string, error) {
	payload := serializer.LoginReq{
		Email:     email,
		Password:  password,
		LongLived: true,
	}

	var loginResp serializer.LoginResp
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&payload).
		SetResult(&loginResp).
		Post("http://dev.circlenetwork.social:5977/api/log-in")
	if err != nil {
		return "", err
	}

	return loginResp.AccessToken, nil
}

func UserFeedAPI(client *resty.Client, accessToken string, latitude float64, longitude float64) error {
	location := serializer.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	var response interface{}
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(&location).
		SetResult(&response).
		Post("http://dev.circlenetwork.social:5977/api/user-feed")
	if err != nil {
		return err
	}

	return nil

}

func PingAPI(client *resty.Client, accessToken string, latitude float64, longitude float64) error {
	location := serializer.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	var response interface{}
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(&location).
		SetResult(&response).
		Post("http://dev.circlenetwork.social:5977/api/ping")
	if err != nil {
		return err
	}

	return nil
}
