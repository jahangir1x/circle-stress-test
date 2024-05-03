package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"os"
	"strconv"
	"stress_test/randomizer"
	"stress_test/repository"
	"stress_test/requests"
	"stress_test/serializer"
	"time"
)

func parseCmd() (args serializer.CliArgs) {
	flag.StringVar(&args.Operation, "action", "", "Action to perform: (login/test)")
	flag.IntVar(&args.Start, "start", 0, "Start number of email")
	flag.IntVar(&args.End, "end", 0, "End number of email")
	flag.StringVar(&args.EmailPrefix, "email-prefix", "not-a-robot-", "Email prefix")
	flag.StringVar(&args.EmailSuffix, "email-suffix", "@gmail.com", "Email suffix")
	flag.StringVar(&args.Password, "password", "Abcd1234", "Password")
	flag.Float64Var(&args.LocationSpread, "location-spread", 0.05, "Location spread to feed into to get random latitude and longitude.")
	flag.Float64Var(&args.WaitTimeMin, "wait-time-min", 0.0, "Minimum wait time between API calls.")
	flag.Float64Var(&args.WaitTimeMax, "wait-time-max", 3.0, "Maximum wait time between API calls.")
	flag.Parse()

	if args.Operation == "" {
		fmt.Println("Action is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if args.Start == 0 || args.End == 0 {
		fmt.Println("Start and end are required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if args.EmailPrefix == "" {
		fmt.Println("Email prefix is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if args.EmailSuffix == "" {
		fmt.Println("Email suffix is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if args.Operation == "login" && args.Password == "" {
		fmt.Println("Password is required for login.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return args
}

func massLogin(client *resty.Client, start int, end int, emailPrefix string, emailSuffix string, password string) {
	for i := start; i <= end; i++ {
		email := emailPrefix + strconv.Itoa(i) + emailSuffix
		fmt.Println("logging in: ", email, " password: ", password)
		token, err := requests.Login(client, email, password)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := repository.SaveToken(email, token); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func testAPIs(client *resty.Client, start int, end int, emailPrefix string, emailSuffix string, locationSpread float64, waitTimeMin float64, waitTimeMax float64) {
	for i := start; i <= end; i++ {
		email := emailPrefix + strconv.Itoa(i) + emailSuffix
		accessToken, err := repository.GetToken(email)
		fmt.Println("user info: email: ", email, " access-token: ", accessToken)
		if err != nil {
			fmt.Println(err)
			return
		}
		go performUserOperations(client, email, accessToken, locationSpread, waitTimeMin, waitTimeMax)
	}

	select {}
}

func performUserOperations(client *resty.Client, email string, accessToken string, locationSpread float64, waitTimeMin float64, waitTimeMax float64) {
	for {
		lat := randomizer.GetRandomLatitude(locationSpread)
		long := randomizer.GetRandomLongitude(locationSpread)

		if err := requests.PingAPI(client, accessToken, lat, long); err != nil {
			fmt.Println(err)
			return
		}
		randomSeconds := waitTimeMin + rand.Float64()*(waitTimeMax-waitTimeMin)
		time.Sleep(time.Duration(randomSeconds) * time.Second)

		if err := requests.UserFeedAPI(client, accessToken, lat, long); err != nil {
			fmt.Println(err)
			return
		}
		randomSeconds = waitTimeMin + rand.Float64()*(waitTimeMax-waitTimeMin)
		time.Sleep(time.Duration(randomSeconds) * time.Second)
	}
}

func main() {
	args := parseCmd()

	client := resty.New()

	if args.Operation == "login" {
		massLogin(client, args.Start, args.End, args.EmailPrefix, args.EmailSuffix, args.Password)
	} else if args.Operation == "test" {
		testAPIs(client, args.Start, args.End, args.EmailPrefix, args.EmailSuffix, args.LocationSpread, args.WaitTimeMin, args.WaitTimeMax)
	}
}
