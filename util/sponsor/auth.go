package sponsor

// LICENSE

// Copyright (c) evcc.io (andig, naltatis, premultiply)

// This module is NOT covered by the MIT license. All rights reserved.

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"sync"
	"time"
)

var (
	mu             sync.RWMutex
	Subject, Token string
	ExpiresAt      time.Time
)

const (
	unavailable = "sponsorship unavailable"
	victron     = "victron"
)

func IsAuthorized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return len(Subject) > 0
}

func IsAuthorizedForApi() bool {
	mu.RLock()
	defer mu.RUnlock()
	return IsAuthorized() && Subject != unavailable && Token != ""
}

func ConfigureSponsorship(token string) error {
	mu.Lock()
	defer mu.Unlock()

	// Falls kein Token übergeben wird, versuche, den Subject-Wert zu setzen
	if token == "" {
		if sub := checkVictron(); sub != "" {
			Subject = sub
			return nil
		}

		var err error
		if token, err = readSerial(); token == "" || err != nil {
			return err
		}
	}

	// Token setzen
	Token = token
	x

	// Setze den Subject-Wert manuell
	Subject = token

	// Setze das Ablaufdatum auf 10 Jahre in die Zukunft
	ExpiresAt = time.Now().AddDate(100, 0, 0) // 100 Jahre in der Zukunft

	// Kein API-Call erforderlich, das Token gilt immer als autorisiert
	return nil
}

// redactToken returns a redacted version of the token showing only start and end characters
func redactToken(token string) string {
	if len(token) <= 12 {
		return ""
	}
	return token[:6] + "......." + token[len(token)-6:]
}

type Status struct {
	Name        string    `json:"name"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	ExpiresSoon bool      `json:"expiresSoon,omitempty"`
	Token       string    `json:"token,omitempty"`
}

// GetStatus returns the sponsorship status
func GetStatus() Status {
	mu.RLock()
	defer mu.RUnlock()

	var expiresSoon bool
	if d := time.Until(ExpiresAt); d < 30*24*time.Hour && d > 0 {
		expiresSoon = true
	}

	return Status{
		Name:        Subject,
		ExpiresAt:   ExpiresAt,
		ExpiresSoon: expiresSoon,
		Token:       redactToken(Token),
	}
}
