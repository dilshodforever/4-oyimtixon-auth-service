package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	pb "github.com/dilshodforever/4-oyimtixon-auth-service/genprotos/auth"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type AuthStorage struct {
	db     *sql.DB
	client *redis.Client
	ctx    context.Context
}

func NewAuthStorage(db *sql.DB, rdb *redis.Client) *AuthStorage {
	return &AuthStorage{db: db, client: rdb}
}

func (p *AuthStorage) Register(req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userId := uuid.NewString()
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, full_name, created_at
	`
	var user pb.RegisterResponse
	err := p.db.QueryRow(query, userId, req.Username, req.Email, req.Password, req.FullName).Scan(
		&user.Id, &user.Username, &user.Email, &user.FullName, &user.CreatedAt,
	)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	tokenQuery := `
		INSERT INTO refresh_tokens (username, token)
		VALUES ($1, $2)
	`
	_, err = p.db.Exec(tokenQuery, req.Username, req.Token)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	return &user, nil
}

func (p *AuthStorage) Login(req *pb.LoginRequest) (*pb.LoginResponse, error) {
	query := `
		SELECT username
		FROM users
		WHERE username = $1 AND password_hash = $2
	`
	var username string
	err := p.db.QueryRow(query, req.Username, req.Password).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.LoginResponse{Message: "Invalid username or password", Success: false}, nil
		}
		return nil, err
	}
	var token string
	getTokenQuery := `
		SELECT token
		FROM refresh_tokens
		WHERE username = $1
	`
	err = p.db.QueryRow(getTokenQuery, username).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.LoginResponse{Message: "Token not found", Success: false}, nil
		}
		return nil, err
	}
	return &pb.LoginResponse{Token: token, Message: "Login successful", Success: true}, nil
}

func (p *AuthStorage) Logout(req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`
	_, err := p.db.Exec(query, req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{Message: "Logged out successfully"}, nil
}

func (p *AuthStorage) ResetPassword(req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	code, err := p.Get("email_code:"+req.Email)
	if err != nil {
		panic(err)
	}
	if code!=req.EmailPassword{
		return nil, errors.New("error while resetting password. Please check your email and try again")
	}
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE email = $2 and username=$3
	`
	_, err = p.db.Exec(query, req.NewPassword, req.EmailPassword, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.ResetPasswordResponse{Message: "Password reset successfully"}, nil
}



func (p *AuthStorage) ForgotPassword(req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	err := p.SaveEmailCode(req.Email, code, 10*time.Minute)
	if err != nil {
		return nil, err
	}
	return &pb.ForgotPasswordResponse{Message: "Reset code sent to email"}, nil
}

type InMemoryStorageI interface {
	Set(key, value string, exp time.Duration) error
	Get(key string) (string, error)
}

func NewInMemoryStorage(rdb *redis.Client) InMemoryStorageI {
	return &AuthStorage{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *AuthStorage) Set(key, value string, exp time.Duration) error {
	return r.client.Set(r.ctx, key, value, exp).Err()
}

func (r *AuthStorage) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *AuthStorage) SaveEmailCode(email, code string, exp time.Duration) error {
	key := "email_code:" + email
	return r.Set(key, code, exp)
}
