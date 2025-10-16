package positions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/positions"
)

type RequestBody struct {
	Position positions.Position `json:"position"`
}

func GetRequestBody(r *http.Request) (handlers.CRUD, error) {
	var requestBody *RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("failed to get the request body, err: %w", err)
	}

	return &requestBody.Position, nil
}
