package sponsor

import (
	"sync"
	"time"
)

var (
	mu             sync.RWMutex
	Subject, Token string
	ExpiresAt      time.Time
)

const victron = "victron"

func IsAuthorized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return Subject == victron
}

func IsAuthorizedForApi() bool {
	mu.RLock()
	defer mu.RUnlock()
	return Subject == victron
}

func ConfigureSponsorship(token string) error {
	mu.Lock()
	defer mu.Unlock()

	// Immer Victron autorisiert, kein Token nötig
	Subject = victron
	Token = ""
	ExpiresAt = time.Now().AddDate(100, 0, 0) // praktisch nie ablaufend
	return nil
}

func redactToken(token string) string {
	return ""
}

type Status struct {
	Name        string    `json:"name"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	ExpiresSoon bool      `json:"expiresSoon,omitempty"`
	Token       string    `json:"token,omitempty"`
}

func GetStatus() Status {
	mu.RLock()
	defer mu.RUnlock()
	return Status{
		Name:        Subject,
		ExpiresAt:   ExpiresAt,
		ExpiresSoon: false,
		Token:       "",
	}
}
