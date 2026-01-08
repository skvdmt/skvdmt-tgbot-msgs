package model

import (
	"github.com/skvdmt/captcha"
	"github.com/skvdmt/captcha/config"
)

// Captcha Создать каптчу и получить ее изображение и значение.
func Captcha() (image []byte, value string, err error) {
	capt, err := captcha.New(&captcha.Config{
		FontFiles: Config.Fonts,
		FontSizes: &config.FontSizes{
			Min: 80,
			Max: 120,
		},
	})
	if err != nil {
		return nil, "", err
	}
	img, err := capt.GetImage()
	if err != nil {
		return nil, "", err
	}
	return img, capt.Value, nil
}
