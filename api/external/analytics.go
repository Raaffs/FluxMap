package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Raaffs/FluxMap/internal/models"
)


func RequestAndCalculatePERTCPM[T models.Analytic](a []*T) (models.Result, error) {
	resultData := models.Result{}
	if len(a) == 0 || a[0] == nil {
		return models.Result{}, fmt.Errorf("input slice is empty")
	}
	data, err := json.Marshal(a)
	if err != nil {
		return models.Result{}, err
	}
	url,err:=GenerateAnalyticUrl(a);if err!=nil{
		return models.Result{},err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return models.Result{}, err		
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.Result{}, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&resultData.Result);err != nil {
		return models.Result{}, err
	}
	return resultData, nil
}

func GenerateAnalyticUrl[T models.Analytic](a []*T)(string,error){
	url:="http://localhost:5000/api/%s"
		switch any(*a[0]).(type) {
    case models.Cpm:
		url=fmt.Sprintf(url,"cpm")
    case models.Pert:
		url=fmt.Sprintf(url,"pert")
    default:
		return "",fmt.Errorf("unknown type")
    }
	return url,nil
}