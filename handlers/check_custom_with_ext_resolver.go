package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tomkaith13/openfga-authz-engine/utils"
)

func CheckCustomRelationWithExt(w http.ResponseWriter, r *http.Request) {

	var customCheckReqBody CustomCheckRequest
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&customCheckReqBody)
	if err != nil {
		http.Error(w, "Invalid body in POST", http.StatusBadRequest)
		return
	}
	err = utils.CheckImpersonatorWithExternalResolver(FgaClient, ModelId, customCheckReqBody.ImpersonatorId, customCheckReqBody.UserId)
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

	err = utils.Check(FgaClient, ModelId, customCheckReqBody.ImpersonatorId, customCheckReqBody.Relation, customCheckReqBody.CapabilityId)
	if err != nil {
		fmt.Println("Relation check failed:", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Allowed: true"))
	w.WriteHeader(http.StatusOK)

}
