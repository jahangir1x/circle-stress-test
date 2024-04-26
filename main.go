package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"stress_test/randomizer"
	"stress_test/serializer"
	"sync"

	"github.com/go-resty/resty/v2"
)

const accessTokenFile = "access_tokens.json"

func login(client *resty.Client, email string, password string) (string, error) {
	// Check if access token already exists in the file
	accessToken, err := getAccessToken(email)
	if err == nil && accessToken != "" {
		return accessToken, nil
	}

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

	// Save access token to file
	err = saveAccessToken(email, loginResp.AccessToken)
	if err != nil {
		fmt.Println("Error saving access token:", err)
	}

	return loginResp.AccessToken, nil
}

func getAccessToken(email string) (string, error) {
	file, err := os.Open(accessTokenFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var tokens []serializer.AccessToken
	err = json.NewDecoder(file).Decode(&tokens)
	if err != nil {
		return "", err
	}

	for _, token := range tokens {
		if token.Email == email {
			return token.AccessToken, nil
		}
	}

	return "", fmt.Errorf("Access token not found for email: %s", email)
}

func saveAccessToken(email, accessToken string) error {
	file, err := os.OpenFile(accessTokenFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var tokens []serializer.AccessToken
	err = json.NewDecoder(file).Decode(&tokens)
	if err != nil && err != io.EOF {
		return err
	}

	// Append new access token
	tokens = append(tokens, serializer.AccessToken{Email: email, AccessToken: accessToken})

	// Write back to the file
	file.Seek(0, 0)
	err = json.NewEncoder(file).Encode(tokens)
	if err != nil {
		return err
	}

	return nil
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

	client := resty.New()

	if shouldLogIn {
		fmt.Println("log in with lower: ", lower, "upper: ", upper)

		if _, err := os.Stat(accessTokenFile); err == nil {
			// File exists, so remove it
			err := os.Remove(accessTokenFile)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		for i := lower; i <= upper; i++ {
			email, password := "not-a-robot-"+strconv.Itoa(i)+"@gmail.com", "Abcd1234"
			_, err := login(client, email, password)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		return
	}
	fmt.Println("test with: upper: ", upper)
	var wg sync.WaitGroup
	for i := 1; i <= upper; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			email := "not-a-robot-" + strconv.Itoa(i) + "@gmail.com"
			if err != nil {
				fmt.Println(err)
				return
			}
			var wgJob sync.WaitGroup
			wgJob.Add(2)

			accessToken, err := getAccessToken(email)
			if err != nil {
				fmt.Println(err)
				return
			}
			go pingAPI(client, &wg, accessToken)
			go userFeedAPI(client, &wg, accessToken)

			wgJob.Wait()
		}()
	}
	wg.Wait()
}
