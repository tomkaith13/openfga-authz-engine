package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tomkaith13/openfga-authz-engine/utils"
)

type AddImpersonationRequest struct {
	UserId         string `json:"user_id"`
	ImpersonatorId string `json:"impersonator_id"`
}

func AddImpersonationRelationHandler(w http.ResponseWriter, r *http.Request) {
	var addImpersonationReq AddImpersonationRequest
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&addImpersonationReq)
	if err != nil {
		http.Error(w, "Invalid body in POST", http.StatusBadRequest)
		return
	}
	err = utils.CheckImpersonator(FgaClient, ModelId, addImpersonationReq.ImpersonatorId, addImpersonationReq.UserId)
	if err != nil {
		// This means the tuple does not exist or is no longer valid
		// so we delete and re-create it
		utils.DeleteImpersonator(FgaClient, ModelId, addImpersonationReq.ImpersonatorId, addImpersonationReq.UserId)
		err = utils.CreateImpersonator(FgaClient, ModelId, addImpersonationReq.ImpersonatorId, addImpersonationReq.UserId)
		if err != nil {
			http.Error(w, "Unable to add Impersonator tuple. err:"+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("added impersonation tuple"))
}
