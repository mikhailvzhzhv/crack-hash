package service

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/util"
)

type ResultSenderRest struct {
	client *http.Client
}

func NewResultSenderRest() *ResultSenderRest {
	return &ResultSenderRest{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (t *ResultSenderRest) SendResult(result *models.Result) {
	resultJSON := shared.StructToJSON(result)
	if resultJSON == nil {
		log.Printf("Failed to marshal result %+v", result)
		return
	}

	req, err := http.NewRequest(
		http.MethodPatch,
		"http://manager:8080/internal/api/manager/hash/crack/request",
		bytes.NewBuffer(resultJSON),
	)

	if err != nil {
		log.Printf("Failed to send result %+v: %v", result, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Received non-success status %d for result %+v", resp.StatusCode, result)
	}
}
