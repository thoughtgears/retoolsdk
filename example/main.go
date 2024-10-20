package main

import (
	"fmt"
	"os"

	retool "github.com/thoughtgears/retoolsdk"
)

func main() {
	apiKey := os.Getenv("RETOOL_API_KEY")
	endpoint := os.Getenv("RETOOL_ENDPOINT")

	client, err := retool.NewClient(apiKey, endpoint)
	if err != nil {
		panic(err)
	}

	users, err := client.ListUsers(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(users)
}
