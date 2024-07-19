package handlers

import (
	"encoding/json"
	"fmt"
	"messageservice"
	"net/http"
)

func (h *Handler) createMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		clientErr(w, http.StatusMethodNotAllowed, "only post method allowed")
		return
	}
	var input messageservice.MessageItem
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Time == "" || input.Time == "" {
		clientErr(w, http.StatusBadRequest, "invalid input body")
		return
	}
	id, err := h.services.Messages.Create(input)
	if err != nil {
		servErr(w, err, err.Error())
		return
	}
	res, err := JSONStruct(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		servErr(w, err, err.Error())
	}
	fmt.Fprintf(w, "%v", res)

}
