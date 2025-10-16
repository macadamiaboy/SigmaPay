package players

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/macadamiaboy/SigmaPay/internal/handlers"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/players"
)

type RequestBody struct {
	Player players.Player `json:"player"`
}

func GetRequestBody(r *http.Request) (handlers.CRUD, error) {
	var requestBody *RequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("failed to get the request body, err: %w", err)
	}

	return &requestBody.Player, nil
}
