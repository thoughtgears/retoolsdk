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

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	return &response, nil
}

// doSingleRequest is a helper function for making single resource requests.
func doSingleRequest[T any](client *Client, method, url string, body interface{}) (*T, error) {
	resp, err := client.Do(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	// Check if the response is empty
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	responseData, err := decodeResponse[T](resp)
	if err != nil {
		return nil, err
	}

	return &responseData.Data, nil
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
