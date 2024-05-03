package requests

import (
	"fmt"
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
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&payload).
		SetResult(&loginResp).
		Post("http://dev.circlenetwork.social:5977/api/log-in")
	if err != nil {
		return "", err
	}

	fmt.Println("Status:", response.Status())

	return loginResp.AccessToken, nil
}

func UserFeedAPI(client *resty.Client, accessToken string, latitude float64, longitude float64) {
	location := serializer.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	fmt.Println("user feed location: ", location)

	var response interface{}
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(&location).
		SetResult(&response).
		Post("http://dev.circlenetwork.social:5977/api/user-feed")
	if err != nil {
		fmt.Println("Error fetching user feed:", err)
		return
	}

	fmt.Println("User feed API response:", response)
}

func PingAPI(client *resty.Client, accessToken string, latitude float64, longitude float64) {
	location := serializer.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	fmt.Println("ping location: ", location)

	var response interface{}
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		SetBody(&location).
		SetResult(&response).
		Post("http://dev.circlenetwork.social:5977/api/ping")
	if err != nil {
		fmt.Println("Error pinging API:", err)
		return
	}

	fmt.Println("Ping API response:", response)
}
