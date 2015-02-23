package repass

import (
	"code.google.com/p/gomock/gomock"
	"github.com/plimble/mailba/mock_mailba"
	"github.com/plimble/moment/mock_moment"
	"github.com/plimble/unik/mock_unik"

	"github.com/plimble/mailba"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRePass(t *testing.T) {
	assert := assert.New(t)

	//mock
	ctrl := gomock.NewController(t)
	store := NewMockStore(ctrl)
	sender := mock_mailba.NewMockSender(ctrl)
	unik := mock_unik.NewMockGenerator(ctrl)
	moment := mock_moment.NewMockTime(ctrl)

	s := NewService(&Config{
		Sender: sender,
		Store:  store,
		Unik:   unik,
		Moment: moment,
	})

	mail := mailba.NewMail("test@plimble.com", "Test")
	mail.AddTo("witooh@icloud.com", "Witoo")
	mail.AddGlobalVar("name", "Witoo")
	mail.SetSubject("Reset password *|name|*")
	mail.SetContent("<h1>*|name|*</h1><p>Token: *|token|*</p>")

	now := time.Now()
	duration := time.Hour * 24 * 10
	expire := now.Add(duration)

	token := &Token{
		ID:     "1111",
		UserID: "1",
		Email:  mail.To[0].Email,
		Expire: expire,
	}

	moment.EXPECT().Now().Return(now)
	moment.EXPECT().Add(now, duration).Return(expire)
	unik.EXPECT().Generate().Return("1111")
	store.EXPECT().Insert(token).Return(nil)
	sender.EXPECT().Send(mail, nil).Return(nil)

	token, err := s.SendResetPasswordMail("1", mail, duration)

	assert.NoError(err)
	assert.NotNil(token)
	assert.Equal(time.Now().Add(time.Hour*24*10).Day(), token.Expire.Day())
	assert.False(token.Used)

	token.Used = true
	store.EXPECT().Get(token.ID).Return(token, nil)
	store.EXPECT().Update(token.ID, token).Return(nil)
	err = s.UseToken(token.ID)
	assert.NoError(err)

	store.EXPECT().Get(token.ID).Return(token, nil)
	token, err = s.GetToken(token.ID)
	assert.NoError(err)
	assert.True(token.Used)
}
