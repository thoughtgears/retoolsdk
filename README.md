# Retool API Client

Standardized API client for the Retool API.

## Usage

### Installation

Install the package using `go get`:

```shell
go get -u github.com/thoughtgears/retoolsdk
```

### Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/thoughtgears/retoolsdk"
)
	
func main() {
    // Create a new Retool client
    client, err := retoolsdk.NewClient(os.Getenv("RETOOL_API_KEY"), os.Getenv("RETOOL_ENDPOINT"))
    if err != nil {
        log.Fatalf("Failed to create Retool client: %v", err)
    }

    users, err := client.ListUsers(nil)
    if err != nil {
        log.Fatalf("Failed to list users: %v", err)
    }
}
````

## API Documentation
The Retool API is documented using the OpenAPI 3.0 format and available at 
[https://api.retool.com/api/v2/spec](https://api.retool.com/api/v2/spec). All API documentation can be found on the 
[Retool API Docs](https://docs.retool.com/reference/api/v2/) page.

### Pagination

By default, responses include up to 100 items. When there are more items, the has_more field in the response is set to 
true and the next field has a pagination token.

```json
{
  "data": [
    {...},
  ],
  "has_more": true,
  "next": "opaque-pagination-token"
}
```

To request the next page of results, include the token in the URL parameters of the request:

```shell
curl -H "Authorization: Bearer {token}" https://api.retool.com/api/v2/resources?next=opaque-pagination-token
```

## Contributing

We welcome contributions! Please follow these steps to contribute:
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Make your changes and commit them with clear and concise messages.
4. Push your changes to your fork.
5. Open a pull request to the main repository.

## Issues

If you encounter any issues, please report them using the GitHub issue tracker. Provide as much detail as possible, 
including steps to reproduce the issue and any relevant logs or screenshots.

