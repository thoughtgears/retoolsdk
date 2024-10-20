package main

import (
	"fmt"
	"os"

	retool "github.com/thoughtgears/retool-sdk"
)

func main() {
	apiKey := os.Getenv("RETOOL_API_KEY")
	endpoint := os.Getenv("RETOOL_ENDPOINT")
	userID := os.Getenv("RETOOL_USER_ID")

	client, err := retool.NewClient(apiKey, endpoint)
	if err != nil {
		panic(err)
	}

	user, err := client.GetUser(userID)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

	users, err := client.ListUsers(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(users)
}
