package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	// openfgaClient "github.com/openfga/go-sdk/client"
	"github.com/tomkaith13/openfga-authz-engine/utils"
)

type CheckManagesRequest struct {
	ImpersonatorId string `json:"impersonator_id"`
	UserId         string `json:"user_id"`
}

func CheckManages(w http.ResponseWriter, r *http.Request) {

	var checkManagesRequest CheckManagesRequest
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&checkManagesRequest)
	if err != nil {
		http.Error(w, "Invalid body in POST", http.StatusBadRequest)
		return
	}
	err = utils.CheckManages(FgaClient, ModelId, checkManagesRequest.ImpersonatorId, checkManagesRequest.UserId)
	if err != nil {
		http.Error(w, "Impersonator has no access to Impersonated", http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		fmt.Println("impersonator check failed:", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Write([]byte("Allowed: true"))
	w.WriteHeader(http.StatusOK)

}
