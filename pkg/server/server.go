package server

import (
	"fmt"
	"github.com/mirshodNasilloyev/tg-bot-youtube-go/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/http"
	"strconv"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketClient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	authResp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.tokenRepository.Save(chatID, authResp.AccessToken, repository.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Printf("chat_id %d\n request_token %s\n access_token %s", chatID, requestToken, authResp.AccessToken)
	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)

}
