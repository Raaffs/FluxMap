package external

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func PERT()(map[string]any,error){
	pertData:=make(map[string]any)
	
	resp,err:=http.Get("http://localhost:5000/pert");if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}
	err=json.NewDecoder(resp.Body).Decode(&pertData);if err!=nil{
		return nil,err
	}
	return pertData,nil
}