package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/repository"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}
func (r *TokenRepository) Save(chatID int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(token))
	})
}

func (r *TokenRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("token is empty")
	}
	return token, nil
}
func intToBytes(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
}
