package authn

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	webAuthn *webauthn.WebAuthn
	err      error
)

// User represents the user model
type User struct {
	id          uint64
	name        string
	displayName string
	credentials []webauthn.Credential
}

// NewUser creates and returns a new User
func NewUser(name string, displayName string) *User {

	user := &User{}
	user.id = randomUint64()
	user.name = name
	user.displayName = displayName
	user.credentials = []webauthn.Credential{}

	return user
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.id))
	return buf
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.name
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.displayName
}

// WebAuthnIcon is not (yet) implemented
func (u User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

var datastore *Datastore

type Datastore struct {
	users    map[string]*User
	sessions map[string]*webauthn.SessionData
	mu       sync.Mutex
}

func NewDatasotre() *Datastore {
	if datastore != nil {
		return datastore
	}
	datastore = &Datastore{users: map[string]*User{}, sessions: map[string]*webauthn.SessionData{}}
	return datastore
}

func (db *Datastore) GetUser(name string) (*User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	user, ok := db.users[name]
	if !ok {
		return &User{}, fmt.Errorf("error getting user '%s': does not exist", name)
	}

	return user, nil
}

func (db *Datastore) PutUser(user *User) {

	db.mu.Lock()
	defer db.mu.Unlock()
	db.users[user.name] = user
}

func (db *Datastore) SaveWebAuthnSession(key string, sessionData *webauthn.SessionData) {
	db.sessions[key] = sessionData
}

func (db *Datastore) GetWebAuthnSession(key string) (*webauthn.SessionData, error) {
	val, ok := db.sessions[key]
	if !ok {
		return nil, fmt.Errorf("error getting session '%s': does not exist", key)
	}
	return val, nil
}

func NewWebAuthNService() (*webauthn.WebAuthn, error) {
	wconfig := &webauthn.Config{
		RPDisplayName: "yogo webauthn",                   // Display Name for your site
		RPID:          "localhost",                       // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:8888"}, // The origin URLs allowed for WebAuthn requests
	}

	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Println(err)
	}
	return webAuthn, nil
}
