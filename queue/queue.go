package queue

import (
	"encoding/json"
	"net/http"

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

type Queue struct {
	NewPoints  []QueueEntry `json:"points"`
	Clarifiers []QueueEntry `json:"clarifiers"`
	wsServer   *dq_websocket.WsServer
}

func SetupQueue(wsServer *dq_websocket.WsServer) *Queue {
	return &Queue{
		NewPoints:  make([]QueueEntry, 0),
		Clarifiers: make([]QueueEntry, 0),
		wsServer:   wsServer,
	}
}

// func (queue *Queue) LeaveQueue(w http.ResponseWriter, r *http.Request) {
// 	userInfo := auth.GetUserClaims(r)
//
// 	requestData := QueueRequestData{}
// 	json.NewDecoder(r.Body).Decode(&requestData)
//
// 	indexOfEntry := slices.IndexFunc(queue.Entries, func(slice QueueEntry) bool {
// 		return slice.Username == userInfo.Username && slice.Type == requestData.Type
// 	})
//
// 	queue.Entries = slices.Concat(queue.Entries[:indexOfEntry], queue.Entries[(indexOfEntry+1):])
// }

func (queue *Queue) NewClarifier(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	requestData := QueueRequestData{}
	json.NewDecoder(r.Body).Decode(&requestData)

	newEntry := QueueEntry{
		Name:     userInfo.Name,
		Username: userInfo.Username,
		Type:     "clarifier",
	}

	queue.Clarifiers = append(queue.Clarifiers, newEntry)

	w.WriteHeader(http.StatusOK)

	queue.wsServer.SendWSMessage(struct {
		Type string     `json:"type"`
		Data QueueEntry `json:"data"`
	}{
		Type: "clarifier",
		Data: newEntry,
	})
}

func (queue *Queue) NewPoint(w http.ResponseWriter, r *http.Request) {
	userInfo := auth.GetUserClaims(r)

	requestData := QueueRequestData{}
	json.NewDecoder(r.Body).Decode(&requestData)

	newEntry := QueueEntry{
		Name:     userInfo.Name,
		Username: userInfo.Username,
		Type:     "point",
	}

	queue.NewPoints = append(queue.NewPoints, newEntry)

	w.WriteHeader(http.StatusOK)

	queue.wsServer.SendWSMessage(struct {
		Type string     `json:"type"`
		Data QueueEntry `json:"data"`
	}{
		Type: "point",
		Data: newEntry,
	})
}

func (queue *Queue) GetQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(queue)
}
