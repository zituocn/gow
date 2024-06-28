package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type (
	Session struct {
		sessionID        string
		lastTimeAccessed time.Time
		values           map[interface{}]interface{}
	}

	Manager struct {
		cookieName  string
		mu          sync.RWMutex
		maxLifeTime int64
		session     map[string]*Session
	}
)

// NewSessionManager return a session manager
func NewSessionManager(cookieName string, maxLifeTime int64) *Manager {
	mgr := &Manager{
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime,
		session:     nil,
		mu:          sync.RWMutex{},
	}
	go mgr.GC()
	return mgr
}

// Start  session return sessionID
func (m *Manager) Start(w http.ResponseWriter, r *http.Request) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var (
		sessionID string
	)

	ck, err := r.Cookie(m.cookieName)
	if err == nil && ck != nil && ck.Value != "" {
		sessionID = ck.Value
	}
	if sessionID != "" && m.session != nil {
		return sessionID
	}

	sessionID = url.QueryEscape(m.makeNewSessionID())
	session := &Session{
		sessionID:        sessionID,
		lastTimeAccessed: time.Now(),
		values:           make(map[interface{}]interface{}),
	}

	m.session = make(map[string]*Session)
	m.session[sessionID] = session

	cookie := http.Cookie{
		Name:     m.cookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(m.maxLifeTime),
	}
	http.SetCookie(w, &cookie)
	return sessionID
}

// End  session
//
//	delete sessionID from cookie and map
func (m *Manager) End(w http.ResponseWriter, r *http.Request) {
	ck, err := r.Cookie(m.cookieName)
	if err != nil || ck.Value == "" {
		return
	} else {
		m.mu.Lock()
		defer m.mu.Unlock()

		//delete map key
		delete(m.session, ck.Value)
		exp := time.Now()
		cookie := http.Cookie{
			Name:     m.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  exp,
			MaxAge:   -1,
		}
		http.SetCookie(w, &cookie)
	}
}

// Extension ext func
func (m *Manager) Extension(w http.ResponseWriter, r *http.Request) string {
	ck, err := r.Cookie(m.cookieName)
	if err != nil || ck == nil {
		return ""
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID := ck.Value

	if session, ok := m.session[sessionID]; ok {
		session.lastTimeAccessed = time.Now()
		return sessionID
	}

	return ""
}

// Get get session value
func (m *Manager) Get(sessionID string, key interface{}) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, ok := m.session[sessionID]; ok {
		if val, ok1 := session.values[key]; ok1 {
			return val, ok1
		}
	}
	return nil, false
}

// Set  session
func (m *Manager) Set(sessionID string, key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, ok := m.session[sessionID]; ok {
		session.values[key] = value
	}
}

// Delete delete value by key
func (m *Manager) Delete(sessionID string, key interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, ok := m.session[sessionID]; ok {
		delete(session.values, key)
	}
}

// GC session gc
func (m *Manager) GC() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for sessionID, session := range m.session {
		if session.lastTimeAccessed.Unix()+m.maxLifeTime < time.Now().Unix() {
			delete(m.session, sessionID)
		}
	}

	time.AfterFunc(time.Duration(m.maxLifeTime)*time.Second, func() {
		m.GC()
	})
}

// makeNewSessionID
func (m *Manager) makeNewSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		nano := time.Now().UnixNano()
		return strconv.FormatInt(nano, 10)
	}
	return base64.URLEncoding.EncodeToString(b)
}
