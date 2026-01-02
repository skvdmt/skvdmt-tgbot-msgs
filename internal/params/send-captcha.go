package params

import (
	"bytes"
	"io"
	"mime/multipart"
	"strconv"
)

const (
	caption         = "caption"
	photo           = "photo"
	methodSendPhoto = "sendPhoto"
	captchaFilename = "captcha.png"
)

// SendCaptcha params for send message request
type SendCaptcha struct {
	chatId      int
	caption     string
	captcha     []byte
	contentType string
}

// NewSendCaptcha конструктор
func NewSendCaptcha(chatId int, caption string, captcha []byte) *SendCaptcha {
	return &SendCaptcha{
		chatId:  chatId,
		caption: caption,
		captcha: captcha,
	}
}

// Method
func (s *SendCaptcha) Method() string {
	return methodSendPhoto
}

// Body
func (s *SendCaptcha) Body() (io.Reader, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	defer func() error {
		return w.Close()
	}()
	part, err := w.CreateFormField(chatId)
	if err != nil {
		return nil, err
	}
	if _, err = part.Write([]byte(strconv.Itoa(s.chatId))); err != nil {
		return nil, err
	}
	part, err = w.CreateFormField(caption)
	if err != nil {
		return nil, err
	}
	if _, err = part.Write([]byte(s.caption)); err != nil {
		return nil, err
	}
	part, err = w.CreateFormFile(photo, captchaFilename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(s.captcha); err != nil {
		return nil, err
	}
	s.contentType = w.FormDataContentType()
	return &body, nil
}

// ContentType
func (s *SendCaptcha) ContentType() string {
	return s.contentType
}
