package repository

type Bucket string

const (
	AccessToken  Bucket = "access_token"
	RequestToken Bucket = "request_token"
)

type TokenRepository interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
