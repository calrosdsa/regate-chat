package repository

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	r "message/domain/repository"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type MessagePublishRequest struct {
	Message  r.Message  `json:"message"`
	TypeChat r.TypeChat `json:"type_chat"`
	ChatId   int        `json:"chat_id"`
}

type WsChatUseCase interface {
	// PublishMessage(msg MessagePublishRequest)
}

type WsServer struct {
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	SubscriberMessageBuffer int

	// PublishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	PublishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	Logf func(f string, v ...interface{})

	// serveMux routes the various endpoints to the appropriate handler.
	// serveMux http.ServeMux

	SubscribersMu sync.Mutex
	Subscribers   map[int]map[int]*Subscriber

	// GrupoUseCase r.GrupoUseCase
}

type Subscriber struct {
	Msgs      chan []byte
	CloseSlow func()
}

func (cs *WsServer) Subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request,
	id int, profileId int) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	s := &Subscriber{
		Msgs: make(chan []byte, cs.SubscriberMessageBuffer),
		CloseSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
	}
	cs.AddSubscriber(s, id, profileId)
	defer cs.DeleteSubscriber(s, id, profileId)

	c2, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}
	c = c2
	mu.Unlock()
	defer c.CloseNow()
	// s1 := Subscriber{}
	// sd,dsd  :=  cs.Subscribers[&s1]
	ctx = c.CloseRead(ctx)
	for {
		select {
		case msg := <-s.Msgs:
			// log.Println(string(msg))
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (cs *WsServer) Publish(msg []byte, id int) {
	cs.SubscribersMu.Lock()
	defer cs.SubscribersMu.Unlock()
	cs.PublishLimiter.Wait(context.Background())
	log.Println(len(cs.Subscribers))
	suscribers := cs.Subscribers[id]
	for _, s := range suscribers {
		select {
		case s.Msgs <- msg:
		default:
			go s.CloseSlow()
		}
	}
	// subscriber,isPresent:=cs.Subscribers[id]
	// if isPresent {
	// 	select{
	// 	case subscriber.Msgs <- msg:
	// 	default:
	// 		go subscriber.CloseSlow()
	// 	}
	// }

}

func (cs *WsServer) AddSubscriber(s *Subscriber, id int, profileId int) {
	log.Println(len(cs.Subscribers), "chats")
	suscribers := cs.Subscribers[id]
	log.Println(len(suscribers), "ADD SUSCRIBER", s)
	if len(suscribers) == 0 {
		log.Println(suscribers, "MAP IS NIL")
		suscribers = make(map[int]*Subscriber)
		cs.Subscribers[id] = suscribers
	}
	log.Println(suscribers, "MAP IS NOT NIL")
	cs.SubscribersMu.Lock()
	// prevSuscription,isPresent := cs.Subscribers[id][profileId]
	// if isPresent {
	// 	log.Println("DELETING PREVIOUS CONNECTIONS")
	// 	cs.DeleteSubscriber(prevSuscription,id,profileId)
	// }

	cs.Subscribers[id][profileId] = s
	cs.SubscribersMu.Unlock()
}

// deleteSubscriber deletes the given Subscriber.
func (cs *WsServer) DeleteSubscriber(s *Subscriber, id int, profileId int) {
	suscribers := cs.Subscribers[id]
	if suscribers != nil {
		if _, ok := suscribers[profileId]; ok {
			cs.SubscribersMu.Lock()
			delete(suscribers, profileId)
			// close(s.Msgs)
			cs.SubscribersMu.Unlock()
			if len(suscribers) == 0 {
				delete(cs.Subscribers, id)
			}
		}
		// delete(cs.Subscribers, id)
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
