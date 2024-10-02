package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Raaffs/FluxMap/internal/models"
)


func RequestAndCalculatePERTCPM[T models.Analytic](a []*T) (models.Result, error) {
	pertData := models.Result{}
	data, err := json.Marshal(a)
	if err != nil {
		return models.Result{}, err
	}

	resp, err := http.Post("http://localhost:5000/api/pert", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return models.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Result{}, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&pertData.Result)
	if err != nil {
		return models.Result{}, err
	}
	log.Println("data pert",pertData)
	return pertData, nil
}