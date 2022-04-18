package email

import (
	"context"
	"testing"

	"project_layout/model/config"

	"go.uber.org/zap"
)

func TestSend(t *testing.T) {
	from := ""
	e := NewEmail(config.Email{
		Host:     "",
		Port:     587,
		Username: from,
		Password: "",
	}, zap.S())
	if err := e.Send(context.Background(), Message{
		To:      "jdlau@qq.com",
		Title:   "Error has happened!",
		Dear:    "Mr.maintainer",
		Content: "报错啦，报错啦，报错啦！！！!",
	}); err != nil {
		t.Fatal(err)
	}
}
