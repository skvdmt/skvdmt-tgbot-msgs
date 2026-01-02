package model

import (
	"path/filepath"

	"github.com/skvdmt/captcha"
	"github.com/skvdmt/captcha/config"
)

// font files paths for captcha
const (
	fontsFolder                  = "/usr/local/share/fonts/"
	fileFontOpenSansBoldItalic   = "OpenSans-BoldItalic.ttf"
	fileFontOswaldBold           = "Oswald-Bold.ttf"
	fileFontGoBold               = "Go-Bold.ttf"
	fileFontRobotoMonoBoldItalic = "RobotoMono-BoldItalic.ttf"
)

// Captcha getting image and meaning
func Captcha() (image []byte, value string, err error) {
	capt, err := captcha.New(&captcha.Config{
		FontFiles: []string{
			filepath.Join(fontsFolder, fileFontGoBold),
			filepath.Join(fontsFolder, fileFontRobotoMonoBoldItalic),
			filepath.Join(fontsFolder, fileFontOswaldBold),
			filepath.Join(fontsFolder, fileFontOpenSansBoldItalic),
		},
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
