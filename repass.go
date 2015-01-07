package repass

import (
	"crypto/md5"
	"fmt"
	"github.com/plimble/errs"
	"github.com/plimble/mailba"
	"time"
)

//go:generate mockgen -destination=mock.go --self_package=github.com/plimble/repass -package=repass github.com/plimble/repass Store,Interface

type Interface interface {
	SendResetPasswordMail(userID string, mail *mailba.Mail, d time.Duration) (*Token, error)
	GetToken(tokenID string) (*Token, error)
	UseToken(tokenID string) error
}

type service struct {
	sender mailba.Sender
	store  Store
}

func NewService(sender mailba.Sender, store Store) *service {
	return &service{sender, store}
}

func encode(email string) string {
	hash := md5.New()
	hash.Write([]byte(time.Now().UTC().Format(time.RFC3339) + "," + email))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (s *service) SendResetPasswordMail(userID string, mail *mailba.Mail, d time.Duration) (*Token, error) {
	if len(mail.To) == 0 {
		return nil, errs.NewErrors("To is required")
	}

	if len(mail.To) > 1 {
		return nil, errs.NewErrors("To should only one")
	}

	token := Token{
		ID:     encode(mail.To[0].Email),
		UserID: userID,
		Email:  mail.To[0].Email,
		Expire: time.Now().Add(d),
	}

	if err := s.store.Insert(token); err != nil {
		return nil, err
	}

	mail.AddGlobalVar("token", token.ID)

	return &token, s.sender.Send(mail, nil)
}

func (s *service) UseToken(tokenID string) error {
	token, err := s.store.Get(tokenID)
	if err != nil {
		return err
	}

	if token.Used == true {
		return nil
	}

	token.Used = true
	return s.store.Update(tokenID, map[string]interface{}{"used": true})
}

func (s *service) GetToken(tokenID string) (*Token, error) {
	return s.store.Get(tokenID)
}
