/*
gow session

sam
2021/6/7

r:=gow.New()
r.Use(Session())

 */

package gow

import (
	"fmt"
	"github.com/zituocn/gow/session"
	"strconv"
)

var (
	cookieName     = "gow_session_id"
	sessionManager *session.Manager
	sessionID      string
)

// InitSession   init gow session
//	before using session,please call this function first
func InitSession() {
	sessionManager = session.NewSessionManager(cookieName, 3600)
}

// Session session middleware
//		r := gow.Default()
//		r.Use(gow.Session())
func Session() HandlerFunc {
	return func(c *Context) {
		if sessionManager == nil {
			panic("Please call gow.InitSession() first")
		}
		sessionID = sessionManager.Start(c.Writer, c.Request)
		sessionManager.Extension(c.Writer, c.Request)

		c.Next()
	}
}

// SetSession set session
func (c *Context) SetSession(key string, v interface{}) {
	setSession(key, v)
}

// GetSession return interface
func (c *Context) GetSession(key string) interface{} {
	return getSession(key)
}

// SessionString return string
func (c *Context) SessionString(key string) string {
	ret := c.GetSession(key)
	v, ok := ret.(string)
	if ok {
		return v
	}
	return ""
}

// SessionInt return int
//		default 0
func (c *Context) SessionInt(key string) int {
	v := c.SessionInt64(key)
	return int(v)
}

// SessionInt64 return int64
//		default 0
func (c *Context) SessionInt64(key string) int64 {
	ret := c.GetSession(key)
	v, err := strconv.ParseInt(fmt.Sprintf("%v", ret), 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// SessionBool return bool
//		default false
func (c *Context) SessionBool(key string) bool {
	ret := c.GetSession(key)
	v, ok := ret.(bool)
	if ok {
		return v
	}
	return false
}

// DeleteSession delete session key
func (c *Context) DeleteSession(key string) {
	deleteSession(key)
}

// getSession getSession
func getSession(key interface{}) interface{} {
	v, ok := sessionManager.Get(sessionID, key)
	if ok {
		return v
	}
	return nil
}

// setSession setSession
func setSession(key, value interface{}) {
	sessionManager.Set(sessionID, key, value)
}

// deleteSession deleteSession
func deleteSession(key interface{}) {
	sessionManager.Delete(sessionID, key)
}
