package main

import (
	"fmt"
	"stress_test/serializer"
	"sync"

	"github.com/go-resty/resty/v2"
)

func login(client *resty.Client) (string, error) {
	payload := serializer.LoginReq{
		Email:     "not-a-robot-test-thirty-two@gmail.com",
		Password:  "Abcd1234",
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

func pingAPI(client *resty.Client, wg *sync.WaitGroup, accessToken string) {
	defer wg.Done()

	location := serializer.Location{
		Latitude:  85.069,
		Longitude: 180,
	}

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

func userFeedAPI(client *resty.Client, wg *sync.WaitGroup, accessToken string) {
	defer wg.Done()

	location := serializer.Location{
		Latitude:  85.05112878,
		Longitude: 180,
	}

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

func main() {

	client := resty.New()
	accessToken, err := login(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Logged in successfully")
	var wg sync.WaitGroup
	wg.Add(2)

	go pingAPI(client, &wg, accessToken)
	go userFeedAPI(client, &wg, accessToken)

	wg.Wait()
	fmt.Println("Both goroutines completed execution")

}
