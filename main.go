package main

import (
	"fmt"
	"github.com/divan/num2words"
	"github.com/go-resty/resty/v2"
	"stress_test/serializer"
	"strings"
	"sync"
)

func login(client *resty.Client, email string, password string) (string, error) {
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

func pingAPI(client *resty.Client, wg *sync.WaitGroup, accessToken string) {
	defer wg.Done()

	location := serializer.Location{
		Latitude:  getRandomLatitude(),
		Longitude: getRandomLatitude(),
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

func userFeedAPI(client *resty.Client, wg *sync.WaitGroup, accessToken string) {
	defer wg.Done()

	location := serializer.Location{
		Latitude:  getRandomLatitude(),
		Longitude: getRandomLatitude(),
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

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := resty.New()
			userId := num2words.Convert(i)
			email, password := "not-a-robot-test-"+strings.ReplaceAll(userId, " ", "-")+"@gmail.com", "Abcd1234"
			accessToken, err := login(client, email, password)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Logged in successfully")
			var wgJob sync.WaitGroup
			wgJob.Add(2)

			go pingAPI(client, &wg, accessToken)
			go userFeedAPI(client, &wg, accessToken)

			wgJob.Wait()
			fmt.Println("Both goroutines completed execution")
		}()
	}
	wg.Wait()
}
