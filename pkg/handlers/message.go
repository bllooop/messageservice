package handlers

import (
	"encoding/json"
	"fmt"
	"messageservice"
	"messageservice/pkg/messaging"
	"net/http"
	"time"
)

func (h *Handler) createMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		clientErr(w, http.StatusMethodNotAllowed, "only post method allowed")
		return
	}
	var input messageservice.MessageItem
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Text == "" {
		clientErr(w, http.StatusBadRequest, "invalid input body")
		return
	}
	input.Time = time.Now()
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
	messaging.MessageToKafka(input)
}

func (h *Handler) consumerKafka(w http.ResponseWriter, r *http.Request) {
	messaging.ConsumeKafka()
}
