package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tomkaith13/openfga-authz-engine/utils"
)

type RefreshImpersonationRequest struct {
	UserId         string `json:"user_id"`
	ImpersonatorId string `json:"impersonator_id"`
}

func RefreshImpersonationHandler(w http.ResponseWriter, r *http.Request) {
	var refreshRequest RefreshImpersonationRequest
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		http.Error(w, "Invalid body in POST", http.StatusBadRequest)
		return
	}
	// This means the tuple does not exist or is no longer valid
	// so we delete and re-create it
	err = utils.DeleteAndAddImpersonator(FgaClient, ModelId, refreshRequest.ImpersonatorId, refreshRequest.UserId)
	if err != nil {
		http.Error(w, "Unable to refresh Impersonator tuple. err:"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("added impersonation tuple"))
}
