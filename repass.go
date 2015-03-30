package repass

import (
	"github.com/plimble/mailba"
	"github.com/plimble/moment"
	"github.com/plimble/unik"
	"github.com/plimble/utils/errors2"
	"time"
)

//go:generate mockgen -destination=mock_repass/mock_service.go github.com/plimble/repass Service

type Service interface {
	SendResetPasswordMail(userID string, mail *mailba.Mail, d time.Duration) (*Token, error)
	GetToken(tokenID string) (*Token, error)
	UseToken(tokenID string) error
}

type RepassService struct {
	*Config
}

type Config struct {
	Sender mailba.Sender
	Store  Store
	Unik   unik.Generator
	Moment moment.Time
}

func NewService(config *Config) *RepassService {
	return &RepassService{config}
}

func (s *RepassService) SendResetPasswordMail(userID string, mail *mailba.Mail, d time.Duration) (*Token, error) {
	if len(mail.To) == 0 {
		return nil, errors2.NewBadReq("To is required")
	}

	if len(mail.To) > 1 {
		return nil, errors2.NewBadReq("To should only one")
	}

	token := &Token{
		ID:     s.Unik.Generate(),
		UserID: userID,
		Email:  mail.To[0].Email,
		Expire: s.Moment.Add(s.Moment.Now(), d),
	}

	if err := s.Store.Insert(token); err != nil {
		return nil, err
	}

	mail.AddGlobalVar("token", token.ID)

	return token, s.Sender.Send(mail)
}

func (s *RepassService) UseToken(tokenID string) error {
	var token *Token
	token, err := s.Store.Get(tokenID)
	if err != nil {
		return err
	}

	if token.Used == true {
		return nil
	}

	token.Used = true
	return s.Store.Update(tokenID, map[string]interface{}{"used": true})
}

func (s *RepassService) GetToken(tokenID string) (*Token, error) {
	return s.Store.Get(tokenID)
}
