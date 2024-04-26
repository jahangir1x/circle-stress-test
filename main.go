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

func parseCmd() (shouldLogIn bool, lowerBound int, upperBound int, err error) {
	flag.Parse()
	arguments := flag.Args()
	if arguments[0] == "login" {
		shouldLogIn = true
	} else if arguments[0] == "test" {
		shouldLogIn = false
	}
	if shouldLogIn {
		lowerBound, err = strconv.Atoi(arguments[1])
		if err != nil {
			fmt.Println(err)
			return false, 0, 0, err
		}
		upperBound, err := strconv.Atoi(arguments[2])
		if err != nil {
			fmt.Println(err)
			return false, 0, 0, err
		}
		return true, lowerBound, upperBound, nil
	} else {
		totalUsers, err := strconv.Atoi(arguments[1])
		if err != nil {
			fmt.Println(err)
			return false, 0, 0, err
		}
		return false, 0, totalUsers, nil
	}
}

func main() {
	shouldLogIn, lower, upper, err := parseCmd()
	if err != nil {
		fmt.Println(err)
		return
	}

	if shouldLogIn {
		fmt.Println("log in with lower: ", lower, "upper: ", upper)
	} else {
		fmt.Println("test with: upper: ", upper)
	}

	client := resty.New()

	var wg sync.WaitGroup
	for i := 1; i <= upper; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
