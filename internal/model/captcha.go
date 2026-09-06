package model

import (
	"github.com/skvdmt/captcha"
)

// Captcha Создать каптчу и получить ее изображение и значение.
func Captcha() (image []byte, value string, err error) {
	capt, err := captcha.New(
		captcha.WithFontFiles(Config.Fonts),
		captcha.WithFontSizes(80, 120),
	)
	if err != nil {
		return nil, "", err
	}
	img, err := capt.GetImage()
	if err != nil {
		return nil, "", err
	}
	return img, capt.Value, nil
}
