package game

import (
	"duel-masters/game/match"
	"duel-masters/internal"
	"duel-masters/server"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var Matchmaker = &Matchmaking{
	requests:    internal.NewConcurrentDictionary[MatchRequest](),
	broadcaster: func(msg interface{}) { logrus.Warn("Use of default matchmaking broadcaster") },
}

type MatchUser struct {
	ID       string
	Username string
	Color    string
	SocketID string
}

type MatchRequest struct {
	ID           string
	Host         MatchUser
	Guest        *internal.Option[MatchUser]
	Name         string
	Format       match.Format
	BlockedUsers internal.ConcurrentDictionary[bool]
}

func (r *MatchRequest) Serialize() server.MatchRequestMessage {
	msg := server.MatchRequestMessage{
		ID:        r.ID,
		HostID:    r.Host.ID,
		HostName:  r.Host.Username,
		HostColor: r.Host.Color,
		Format:    string(r.Format),
	}

	guest, ok := r.Guest.Unwrap()

	if !ok {
		return msg
	}

	msg.GuestID = guest.ID
	msg.GuestName = guest.Username
	msg.GuestColor = guest.Color

	return msg
}

type Matchmaking struct {
	sync.RWMutex
	requests    internal.ConcurrentDictionary[MatchRequest]
	broadcaster func(msg interface{})
	matchSystem *match.MatchSystem
}

func (m *Matchmaking) Initialize(f func(msg interface{}), sys *match.MatchSystem) {
	m.broadcaster = f
	m.matchSystem = sys
}

func (m *Matchmaking) NewRequest(s *server.Socket, name string, format match.Format) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.requests.Find(s.User.UID); ok {
		return fmt.Errorf("you already have a duel request open")
	}

	r := &MatchRequest{
		ID: s.User.UID,
		Host: MatchUser{
			ID:       s.User.UID,
			Username: s.User.Username,
			Color:    s.User.Color,
			SocketID: s.UID,
		},
		Guest:        internal.NewOption[MatchUser](),
		Name:         name,
		Format:       format,
		BlockedUsers: internal.NewConcurrentDictionary[bool](),
	}

	m.requests.Add(s.User.UID, r)
	m.BroadcastState()

	return nil
}

func (m *Matchmaking) Join(s *server.Socket, id string) error {
	m.Lock()
	defer m.Unlock()

	// check if guest is in another request
	for _, r := range m.requests.Iter() {
		if guest, ok := r.Guest.Unwrap(); ok && guest.ID == s.User.UID {
			return fmt.Errorf("you are already requesting to join a duel. Leave the current request and try again")
		}
	}

	r, ok := m.requests.Find(id)

	if !ok {
		return fmt.Errorf("the duel request you attempted to join does not exist")
	}

	if r.Guest.Some() {
		return fmt.Errorf("someone has already joined that duel")
	}

	r.Guest.Set(&MatchUser{
		ID:       s.User.UID,
		Username: s.User.Username,
		Color:    s.User.Color,
		SocketID: s.UID,
	})

	m.BroadcastState()

	return nil
}

func (m *Matchmaking) Leave(s *server.Socket) {
	m.Lock()
	defer m.Unlock()

	var broadcast bool

	for _, r := range m.requests.Iter() {
		guest, ok := r.Guest.Unwrap()

		if r.Host.ID == s.User.UID {
			m.requests.Remove(r.ID)

			if ok {
				guestSocket, _ := server.Sockets.Find(guest.SocketID)
				if guestSocket != nil {
					guestSocket.Warn("The match you requested to join has been closed by the host")
				}
			}

			broadcast = true
		} else if ok && guest.ID == s.User.UID {
			r.Guest.Clear()
			broadcast = true
		}
	}

	if broadcast {
		m.BroadcastState()
	}
}

func (m *Matchmaking) Kick(s *server.Socket, requestId string, toKickId string) {
	m.Lock()
	defer m.Unlock()

	r, ok := m.requests.Find(requestId)

	if !ok {
		s.Warn("Failed to kick user from your duel request, it does not exist")
		return
	}

	if len(r.BlockedUsers.Iter()) > 25 {
		s.Warn("You have reached the limit of how many users can be blocked from joining your match")
		return
	}

	r.BlockedUsers.Add(toKickId, nil)

	guest, ok := r.Guest.Unwrap()

	if !ok || guest.ID != toKickId {
		s.Warn("The user you tried to kick had already left, they are now blocked from joining again")
		return
	}

	r.Guest.Clear()
	guestSocket, ok := server.Sockets.Find(guest.SocketID)
	if !ok {
		guestSocket.Warn("You were kicked from the match")
	}

	s.Warn(fmt.Sprintf("Successfully kicked %s, they are now blocked from joining again", guest.Username))
	m.BroadcastState()
}

func (m *Matchmaking) Start(s *server.Socket, requestId string) {
	m.Lock()
	defer m.Unlock()

	r, ok := m.requests.Find(requestId)

	if !ok {
		s.Warn("The match you attempted to start no longer exist")
		return
	}

	if r.Host.ID != s.User.UID {
		s.Warn("Only the host can start the match")
		return
	}

	guest, ok := r.Guest.Unwrap()

	if !ok {
		s.Warn("Can't start the match because there's no opponent")
		return
	}

	guestSocket, ok := server.Sockets.Find(guest.SocketID)

	if !ok {
		s.Warn("Failed to communicate with the opposing player, consider kicking them")
		return
	}

	match := m.matchSystem.NewMatch(r.Name, r.Host.ID, true)

	msg := server.MatchForwardMessage{
		Header: "match_forward",
		ID:     match.ID,
	}

	s.Send(msg)
	guestSocket.Send(msg)

}

func (m *Matchmaking) Serialize() server.MatchReuestsListMessage {
	msg := server.MatchReuestsListMessage{
		Header:   "match_requests",
		Requests: []server.MatchRequestMessage{},
	}

	for _, r := range m.requests.Iter() {
		msg.Requests = append(msg.Requests, r.Serialize())
	}

	return msg
}

func (m *Matchmaking) BroadcastState() {
	m.broadcaster(m.Serialize())
}
