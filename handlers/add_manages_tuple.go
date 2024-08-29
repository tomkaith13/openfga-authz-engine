package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tomkaith13/openfga-authz-engine/utils"
)

type AddManagesTupleRequest struct {
	ImpersonatorId string   `json:"impersonator_id"`
	UserIds        []string `json:"user_ids"`
}

func AddManagesRelationHandler(w http.ResponseWriter, r *http.Request) {
	var addImpersonationReq AddManagesTupleRequest
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&addImpersonationReq)
	if err != nil {
		http.Error(w, "Invalid body in POST", http.StatusBadRequest)
		return
	}
	err = utils.CreateUserManagementTuples(FgaClient, ModelId, addImpersonationReq.ImpersonatorId, addImpersonationReq.UserIds)
	if err != nil {
		if errors.Is(err, utils.ErrorAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Unable to add manages tuple. err:"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("added manages tuple"))
}
