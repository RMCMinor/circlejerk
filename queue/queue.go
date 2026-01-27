package queue

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/bigspaceships/circlejerk/auth"
	dq_websocket "github.com/bigspaceships/circlejerk/websocket"
)

type QueueEntry struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

type QueueRequestData struct {
	Type string `json:"type"`
}

var queue []QueueEntry

func SetupQueue() {
	queue = make([]QueueEntry, 0)
}

func LeaveQueue(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	requestData := QueueRequestData{}
	json.NewDecoder(r.Body).Decode(&requestData)

	indexOfEntry := slices.IndexFunc(queue, func(slice QueueEntry) bool {
		return slice.Username == userInfo.Username && slice.Type == requestData.Type
	})

	queue = slices.Concat(queue[:indexOfEntry], queue[(indexOfEntry+1):])
}

func JoinQueue(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	requestData := QueueRequestData{}
	json.NewDecoder(r.Body).Decode(&requestData)

	newEntry := QueueEntry{
		Name:     userInfo.Name,
		Username: userInfo.Username,
		Type:     requestData.Type,
	}

	queue = append(queue, newEntry)

	w.WriteHeader(http.StatusOK)

	dq_websocket.SendWSMessage(struct {
		Type        string     `json:"type"`
		Data QueueEntry `json:"data"`
	}{
		Type:        "new-point",
		Data: newEntry,
	})
}

func GetQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(queue)
}
