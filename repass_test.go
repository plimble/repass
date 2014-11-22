package repass

import (
	"code.google.com/p/gomock/gomock"
	"github.com/plimble/mailba"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRePass(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	store := NewMockStore(ctrl)
	sender := mailba.NewMockSender(ctrl)

	s := NewService(sender, store)

	mail := mailba.NewMail("test@plimble.com", "Test")
	mail.AddTo("witooh@icloud.com", "Witoo")
	mail.AddGlobalVar("name", "Witoo")
	mail.SetSubject("Reset password *|name|*")
	mail.SetContent("<h1>*|name|*</h1><p>Token: *|token|*</p>")

	store.EXPECT().Create(gomock.Any()).Return(nil)
	sender.EXPECT().Send(gomock.Any(), gomock.Any())

	token, err := s.SendResetPasswordMail(mail, time.Hour*24*10)

	assert.NoError(err)
	assert.NotNil(token)
	assert.Equal(time.Now().Add(time.Hour*24*10).Day(), token.Expire.Day())
	assert.False(token.Used)

	store.EXPECT().Get(token.ID).Return(token, nil)
	store.EXPECT().Update(token.ID, gomock.Any()).Return(nil)
	err = s.UseToken(token.ID)
	assert.NoError(err)

	store.EXPECT().Get(token.ID).Return(token, nil)
	token, err = s.GetToken(token.ID)
	assert.NoError(err)
	assert.True(token.Used)
}
