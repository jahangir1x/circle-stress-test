package main

import (
	"flag"
	"fmt"
	"strconv"
	"stress_test/randomizer"
	"stress_test/serializer"
	"sync"

	"github.com/go-resty/resty/v2"
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
		Latitude:  randomizer.GetRandomLatitude(),
		Longitude: randomizer.GetRandomLongitude(),
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
		Latitude:  randomizer.GetRandomLatitude(),
		Longitude: randomizer.GetRandomLongitude(),
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

func parseCmd() (int, error) {
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) != 1 {
		return 0, fmt.Errorf("Usage:   go run main.go <total_users>\nexample: go run main.go 10")
	}
	num, err := strconv.Atoi(arguments[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid number of users")
	}
	return num, nil
}

func main() {
	totalUsers, err := parseCmd()
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	for i := 1; i <= totalUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := resty.New()
			email, password := "not-a-robot-"+strconv.Itoa(i)+"@gmail.com", "Abcd1234"
			fmt.Println("trying to login with email:", email, "password:", password)
			accessToken, err := login(client, email, password)
			if err != nil {
				fmt.Println(err)
				return
			}
			var wgJob sync.WaitGroup
			wgJob.Add(2)

			go pingAPI(client, &wg, accessToken)
			go userFeedAPI(client, &wg, accessToken)

			wgJob.Wait()
		}()
	}
	wg.Wait()
}
