package requester

import (
	"context"
	"dummy-endpoints-http/structs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func MakeWG(url string) (*structs.Response, error) {
	var req *http.Request
	var err error
	req, err = http.NewRequest("GET", url, http.NoBody)

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(context.Background(), structs.RequestTimeOut*time.Millisecond)
	defer cancel()

	// Make the request
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	// Read the response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for non-200 status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", res.StatusCode)
	}

	// Parse the JSON response
	var responses structs.Response
	if err := json.Unmarshal(resBody, &responses); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &responses, nil
}
