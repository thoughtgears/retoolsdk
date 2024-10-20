package retoolsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// decodeResponse is a helper function to decode JSON responses
func decodeResponse[T any](resp *http.Response) (*Response[T], error) {
	var response Response[T]

	// Check if the response is empty
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	return &response, nil
}

// doPaginatedRequest is a helper function to make paginated requests to the API.
func doPaginatedRequest[T any](client *Client, baseURL string, query url.Values) ([]T, error) {
	var allItems []T
	var nextToken string
	hasMore := true

	for hasMore {
		if nextToken != "" {
			query.Set("next", nextToken)
		}

		urlWithQuery := fmt.Sprintf("%s?%s", baseURL, query.Encode())
		resp, err := client.Do("GET", urlWithQuery, nil)
		if err != nil {
			return nil, fmt.Errorf("making request: %w", err)
		}

		responseData, err := decodeResponse[[]T](resp)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, responseData.Data...)
		nextToken = responseData.NextToken
		hasMore = responseData.HasMore
	}

	return allItems, nil
}
