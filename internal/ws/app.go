package ws

import (
	"net/http"
	"sync"
)

var house sync.Map
var roomMutexes = make(map[string]*sync.Mutex)
var mutexForRoomMutexes = new(sync.Mutex)

func Start(roomId string, w http.ResponseWriter, r *http.Request) {
	mutexForRoomMutexes.Lock()
	roomMutex, ok := roomMutexes[roomId]
	if ok {
		roomMutex.Lock()
	} else {
		roomMutexes[roomId] = new(sync.Mutex)
		roomMutexes[roomId].Lock()
	}
	mutexForRoomMutexes.Unlock()
	room, ok := house.Load(roomId)
	var hub *Hub
	if ok {
		hub = room.(*Hub)
	} else {
		hub = NewHub(roomId)
		house.Store(roomId, hub)
		go hub.Run()
	}
	ServeWs(hub, w, r)
}
