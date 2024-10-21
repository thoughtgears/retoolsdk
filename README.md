# Retool API Client

Standardized API client for the Retool API.

## Implemented features

#### User

- **Get User by ID**: Get a user information by ID.
- **List Users**: List all users in the organization.
- **Create User**: Create a new user with optional parameters.
- **Update User**: Update a user's information by ID.
- **Delete User**: Delete a user by ID.

#### User Attributes

- **Update User Attributes**: Update a user's attributes by ID.
- **Delete User Attributes**: Delete a user's attributes by ID.
- **List Organization Attributes**: List all organization attributes.

#### Group

- **Get Group by ID**: Get a group information by ID.
- **List Groups**: List all groups in the organization.
- **Create Group**: Create a new group with optional parameters.
- **Update Group**: Update a group's information by ID.
- **Delete Group**: Delete a group by ID.

#### Permissions
- 
- **List Permissions**: List all permissions in the organization.
- **Grant Permission**: Grant a permission to a user or group.
- **Revoke Permission**: Revoke a permission from a user or group.

#### Folders

- **Get Folder by ID**: Get a folder information by ID.
- **List Folders**: List all folders in the organization.
- **Create Folder**: Create a new folder with optional parameters.
- **Update Folder**: Update a folder's information by ID.
- **Delete Folder**: Delete a folder by ID.

#### Spaces

- **Get Space by ID**: Get a space information by ID.
- **List Spaces**: List all spaces in the organization.
- **Create Space**: Create a new space with optional parameters.
- **Update Space**: Update a space's information by ID.
- **Delete Space**: Delete a space by ID.

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

## Development

## Contributing

## Issues

