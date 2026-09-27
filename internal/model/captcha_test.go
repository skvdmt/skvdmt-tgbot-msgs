//go:build unit

package model

import (
	"testing"
)

func TestCaptcha(t *testing.T) {
	if err := CreateLogger(); err != nil {
		t.Error(err)
	}
	if err := CreateConfig(); err != nil {
		t.Error(err)
	}
	i, v, err := Captcha()
	if err != nil {
		t.Error(err)
	}
	if len(i) == 0 {
		t.Errorf("Captcha image is empty")
	}
	if len(v) == 0 {
		t.Errorf("Captcha value is empty")
	}
}
