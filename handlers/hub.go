package handlers

type Hub struct {
	Rooms map[uint64]*Room
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[uint64]*Room),
	}
}

func (h *Hub) GetRoom(roomName uint64) *Room {
	if room, ok := h.Rooms[roomName]; ok {
		return room
	}
	room := &Room{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	h.Rooms[roomName] = room
	go room.run()
	return room
}

func (h *Hub) RemoveRoom(roomName uint64) {
	delete(h.Rooms, roomName)
}
