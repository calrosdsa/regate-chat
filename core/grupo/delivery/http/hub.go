package ws

import (
	"context"
	"encoding/json"
	"log"
	"message/domain/repository"
	"strconv"
	// "soporte-go/core/model/ws"
)

type message struct {
	data []byte
	room string
}

type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

// func NewHub()*hub {
// 	return &hub{
// 		broadcast:  make(chan message),
// 		register:   make(chan subscription),
// 		unregister: make(chan subscription),
// 		rooms:      make(map[string]map[*connection]bool),
// 	}
// }
var H = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

func (h *hub) Run(us repository.GrupoUseCase) {
	ctx := context.Background()
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
		case s := <-h.unregister:
			log.Println("Unregister connection")
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			// log.Println("Sending data")
			log.Println(string(m.data))
			event := &repository.GrupoEvent{}
			err := json.Unmarshal(m.data,event)
			if err != nil {
				log.Println("Fail to unmarshall",err)
			}
			if event.Type == "message"{	
				grupoId,_ :=strconv.Atoi(m.room)
				event.Message.GrupoId = grupoId
				err = us.SaveGrupoMessage(ctx,&event.Message)
				if err != nil{
					log.Println("fail to save message")
					log.Println(err)
				}
			}

			jsonStr ,err := json.Marshal(event)
			if err != nil {
				log.Println(err)
			}
			// log.Println(msgData)
			//   <- msgData
			// log.Println("Message Saved",msg.Content)
			// log.Println(msgData)
			// log.Println(connections)
			for c := range connections {
				select {
				case c.send <-jsonStr:
				default:
					// log.Println("close - delete")
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						// log.Println("delete room")
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}

// create table messages(
// 	caso_id text,
// 	content text,
// 	created timestamp,
// 	from_user text,
// 	from_user_id text,
// 	to_user text,
// 	to_user_id text,
// 	read int,
// 	PRIMARY KEY (caso_id,content));
