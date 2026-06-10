package auth

import (
	"brotchaschte2/internal/database"
	query "brotchaschte2/internal/database/generated"
	"context"
	crypt_rand "crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/argon2"

	_ "github.com/mattn/go-sqlite3"
)

type AuthService interface {
	SessionManager() *scs.SessionManager
	AuthenticateUser(ctx context.Context, email, password string) (bool, error)
}

type authservice struct {
	db      *sql.DB
	queries *query.Queries

	sessionManager *scs.SessionManager

	hashParams argon2id_params
}

type argon2id_params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ErrIncompatibleAlgo    = errors.New("incompatible hash algorithm (not argon2id)")

	ErrUserAlreadyExists = errors.New("the user already exists")
)

func New() AuthService {
	db_service := database.New()
	db := db_service.DB()
	queries := db_service.Queries()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)

	sessionManager.Lifetime = 7 * 24 * time.Hour
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	hashParams := argon2id_params{
		memory:      256 * 1024,
		iterations:  3,
		parallelism: 1,
		saltLength:  16,
		keyLength:   32,
	}

	return &authservice{
		db:             db,
		queries:        queries,
		sessionManager: sessionManager,
		hashParams:     hashParams,
	}
}

func (s *authservice) SessionManager() *scs.SessionManager {
	return s.sessionManager
}

func (s *authservice) AuthenticateUser(ctx context.Context, email, password string) (bool, error) {
	user, err := s.queries.GetUser(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("Failed to fetch user with email %s. Got error %s\n", email, err.Error())
		return false, err
	}

	ok, err := comparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		log.Printf("Failed to check password hash for %s. Got error %s\n", email, err.Error())
		return false, err
	}
	return ok, nil
}

func (s *authservice) CreateUser(ctx context.Context, email, password string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("failed to obtain db transaction to create user")
		return err
	}
	defer tx.Rollback()
	q := s.queries.WithTx(tx)

	_, err = q.GetUser(ctx, email)
	if err != sql.ErrNoRows {
		return ErrUserAlreadyExists
	}
	err = q.CreateUser(ctx, query.CreateUserParams{
		Email:        email,
		PasswordHash: s.HashPassword(password),
	})
	if err != nil {
		log.Printf("failed create new user %s\n", email)
		return err
	}

	return tx.Commit()

}

func (s *authservice) HashPassword(password string) string {
	salt := genRandomBytes(s.hashParams.saltLength)
	hash := argon2.IDKey([]byte(password), salt, s.hashParams.iterations, s.hashParams.memory, s.hashParams.parallelism, s.hashParams.keyLength)

	b64_salt := base64.RawStdEncoding.EncodeToString(salt)
	b64_hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded_hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		s.hashParams.memory,
		s.hashParams.iterations,
		s.hashParams.parallelism,
		b64_salt, b64_hash)

	return encoded_hash
}

func genRandomBytes(n uint32) []byte {
	b := make([]byte, n)
	// never returns an error according docs
	_, err := crypt_rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("Failed to generate crypto secure random bytes: %w", err))
	}
	return b
}

func decodeHashParams(encoded_hash string) (p *argon2id_params, salt []byte, hash []byte, err error) {
	vals := strings.Split(encoded_hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}
	if vals[1] != "argon2id" {
		return nil, nil, nil, ErrIncompatibleAlgo
	}
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("version scan: %w", err)
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}
	p = &argon2id_params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("param scan: %w", err)
	}
	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("salt decode: %w", err)
	}
	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("hash decode: %w", err)
	}
	p.keyLength = uint32(len(hash))
	p.saltLength = uint32(len(salt))

	return p, salt, hash, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHashParams(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}
