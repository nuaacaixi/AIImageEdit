package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Simple JSON file-based store for Phase 1.
// Can be replaced with SQLite/PostgreSQL for production.

type Store struct {
	mu       sync.RWMutex
	dataDir  string
	sessions map[string]*Session
	turns    map[string][]Turn
	context  map[string][]ContextMessage
}

var store *Store

func InitDB(dataDir string) error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	store = &Store{
		dataDir:  dataDir,
		sessions: make(map[string]*Session),
		turns:    make(map[string][]Turn),
		context:  make(map[string][]ContextMessage),
	}

	// Load existing data from disk
	store.loadAll()

	return nil
}

func DB() *Store {
	return store
}

// loadAll loads sessions, turns, and context from JSON files
func (s *Store) loadAll() {
	s.loadSessions()
	s.loadTurns()
	s.loadContext()
}

func (s *Store) sessionsPath() string {
	return filepath.Join(s.dataDir, "sessions.json")
}

func (s *Store) turnsPath() string {
	return filepath.Join(s.dataDir, "turns.json")
}

func (s *Store) contextPath() string {
	return filepath.Join(s.dataDir, "context.json")
}

func (s *Store) loadSessions() {
	data, err := os.ReadFile(s.sessionsPath())
	if err != nil {
		return
	}
	var sessions []*Session
	if json.Unmarshal(data, &sessions) == nil {
		for _, sess := range sessions {
			s.sessions[sess.ID] = sess
		}
	}
}

func (s *Store) loadTurns() {
	data, err := os.ReadFile(s.turnsPath())
	if err != nil {
		return
	}
	json.Unmarshal(data, &s.turns)
}

func (s *Store) loadContext() {
	data, err := os.ReadFile(s.contextPath())
	if err != nil {
		return
	}
	json.Unmarshal(data, &s.context)
}

func (s *Store) saveSessions() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []*Session
	for _, sess := range s.sessions {
		list = append(list, sess)
	}
	data, _ := json.MarshalIndent(list, "", "  ")
	os.WriteFile(s.sessionsPath(), data, 0o644)
}

func (s *Store) saveTurns() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, _ := json.MarshalIndent(s.turns, "", "  ")
	os.WriteFile(s.turnsPath(), data, 0o644)
}

func (s *Store) saveContext() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, _ := json.MarshalIndent(s.context, "", "  ")
	os.WriteFile(s.contextPath(), data, 0o644)
}

// Session CRUD

func (s *Store) CreateSession(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	session.CreatedAt = now
	session.UpdatedAt = now
	s.sessions[session.ID] = session

	go s.saveSessions()
	return nil
}

func (s *Store) GetSession(id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", id)
	}

	// Return a copy
	cp := *sess
	return &cp, nil
}

func (s *Store) UpdateSessionCurrent(sessionID, imageID, imageURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	sess.CurrentImageID = imageID
	sess.CurrentImageURL = imageURL
	sess.UpdatedAt = time.Now().UTC()

	go s.saveSessions()
	return nil
}

// Turn CRUD

func (s *Store) CreateTurn(turn *Turn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	turn.CreatedAt = time.Now().UTC()
	s.turns[turn.SessionID] = append(s.turns[turn.SessionID], *turn)

	go s.saveTurns()
	return nil
}

func (s *Store) UpdateTurnResult(turnID, status, resultImageID, resultImageURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for sid, turns := range s.turns {
		for i := range turns {
			if turns[i].ID == turnID {
				turns[i].Status = status
				turns[i].ResultImageID = resultImageID
				turns[i].ResultImageURL = resultImageURL
				s.turns[sid] = turns
				go s.saveTurns()
				return nil
			}
		}
	}
	return fmt.Errorf("turn not found: %s", turnID)
}

func (s *Store) GetTurnsBySession(sessionID string) ([]Turn, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	turns := s.turns[sessionID]
	if turns == nil {
		return []Turn{}, nil
	}

	// Return a copy
	cp := make([]Turn, len(turns))
	copy(cp, turns)
	return cp, nil
}

// ContextMessage CRUD

func (s *Store) AddContextMessage(msg *ContextMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg.CreatedAt = time.Now().UTC()
	s.context[msg.SessionID] = append(s.context[msg.SessionID], *msg)

	go s.saveContext()
	return nil
}

func (s *Store) GetContextMessages(sessionID string, limit int) ([]ContextMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msgs := s.context[sessionID]
	if msgs == nil {
		return []ContextMessage{}, nil
	}

	start := 0
	if len(msgs) > limit {
		start = len(msgs) - limit
	}

	cp := make([]ContextMessage, len(msgs)-start)
	copy(cp, msgs[start:])
	return cp, nil
}
