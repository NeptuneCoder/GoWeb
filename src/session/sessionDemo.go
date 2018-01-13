package main

import "sync"

var globalSessions *session.Manager
//然后在init函数中初始化
func init() {
	globalSessions, _ = NewManager("memory","gosessionid",3600)
}
func main() {

}
type Manager struct {
	cookieName string
	lock sync.Mutex // protects session
	provider Provider
	maxlifetime int64
}
type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}