package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/dto"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/messages"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/params"
)

const (
	cmdStart   = "/start"
	cmdAuth    = "/auth"
	cmdMessage = "/message"
)

// Update
type Update struct {
	Id      int          `json:"update_id"`
	Message *dto.Message `json:"message"`
	bot     *Bot
	user    *User
}

// isMessage
func (u *Update) isMessage() bool {
	return u.Message != nil
}

// isCommand
func (u *Update) isCommand() bool {
	return strings.HasPrefix(u.Message.Text, "/")
}

// Hadndling telegram bot updates
func (u *Update) Handling(ctx context.Context) error {
	// set params.Offset to next update
	if u.Id >= u.bot.paramsGetUpdates.Offset {
		u.bot.paramsGetUpdates.Offset = u.Id + 1
	}
	user, err := u.bot.users.GetOrCreateUserByIdAndName(
		ctx, u.Message.From.Id, u.Message.From.Username,
	)
	if err != nil {
		return err
	}
	u.user = user
	if u.isMessage() {
		switch u.user.IsWant() {
		case BotWantCommand:
			return u.cmdHandling(ctx)
		case BotWantCaptcha:
			return u.captchaHandling(ctx)
		case BotWantMessage:
			return u.messageHandling(ctx)
		default:
			return fmt.Errorf("unknown bot want %d", u.user.IsWant())
		}
	}
	return nil
}

// Handling search for entered command
func (u *Update) cmdHandling(ctx context.Context) error {
	if !u.isCommand() {
		return u.unknown(ctx)
	}
	switch u.Message.Text {
	case cmdStart:
		return u.start(ctx)
	case cmdAuth:
		return u.auth(ctx)
	case cmdMessage:
		return u.message(ctx)
	default:
		return u.unknown(ctx)
	}
}

// message restonse to message command
func (u *Update) message(ctx context.Context) error {
	ok, err := u.user.CanSendMessage()
	if !ok {
		switch err.Error() {
		case ErrNotAuthorized:
			if err := u.send(ctx, fmt.Sprintf(messages.AuthIncomplete,
				cmdAuth)); err != nil {
				return err
			}
		case ErrMessageAlreadySended:
			if err := u.send(ctx, fmt.Sprintf(messages.WillSaved,
				msgsUrl, u.user.SendMessageCooldownLeft())); err != nil {
				return err
			}
		default:
			return err
		}
		u.user.SetWant(BotWantCommand)
		return nil
	}
	if err := u.send(ctx, messages.SendMessageEnter); err != nil {
		return err
	}
	u.user.SetWant(BotWantMessage)
	return nil
}

// unknown response to unknown command
func (u *Update) unknown(ctx context.Context) error {
	if err := u.send(ctx, fmt.Sprintf(messages.UnknownCommand,
		u.Message.Text)); err != nil {
		return err
	}
	if err := u.send(ctx, fmt.Sprintf(messages.Commands,
		cmdAuth, cmdMessage)); err != nil {
		return err
	}
	return nil
}

// start response to start command
func (u *Update) start(ctx context.Context) error {
	if err := u.send(ctx, fmt.Sprintf(messages.Start,
		u.Message.From.Username, name, cmdAuth, cmdMessage, msgsUrl)); err != nil {
		return err
	}
	return nil
}

// auth response to auth command
func (u *Update) auth(ctx context.Context) error {
	ok, err := u.user.CanAuth()
	if !ok {
		switch err.Error() {
		case ErrAlreadyAuthorized:
			if err := u.send(ctx, messages.AlreadyAuth); err != nil {
				return err
			}
			if err := u.send(ctx, fmt.Sprintf(messages.SendMessageTitle,
				cmdMessage)); err != nil {
				return err
			}
		case ErrAuthAttemptsLeft:
			if err := u.send(ctx, fmt.Sprintf(messages.AuthOver,
				u.user.AuthCooldownLeft())); err != nil {
				return err
			}
		default:
			return err
		}
		u.user.SetWant(BotWantCommand)
		return nil
	}
	i, v, err := model.Captcha()
	if err != nil {
		return err
	}
	_, err = u.bot.do(ctx,
		params.NewSendCaptcha(
			u.Message.Chat.Id,
			fmt.Sprintf(
				messages.AuthTitle,
				u.user.AuthAttemptsLeft(),
			),
			i,
		))
	if err != nil {
		return err
	}
	u.user.SetCaptcha(v)
	u.user.SetWant(BotWantCaptcha)
	return nil
}

// captchaHandling
func (u *Update) captchaHandling(ctx context.Context) error {
	if !u.user.Auth(u.Message.Text) {
		if err := u.send(ctx, messages.AuthWrong); err != nil {
			return err
		}
		return u.auth(ctx)
	}
	if err := u.send(ctx, fmt.Sprintf(messages.AuthComplete)); err != nil {
		return err
	}
	if err := u.send(ctx, fmt.Sprintf(messages.SendMessageTitle,
		cmdMessage)); err != nil {
		return err
	}
	u.user.SetWant(BotWantCommand)
	return nil
}

// messageHandling обработка отправленного текстового сообщения
func (u *Update) messageHandling(ctx context.Context) error {
	if err := u.user.SaveMessage(ctx, u.Message.Text); err != nil {
		return err
	}
	if err := u.send(ctx, fmt.Sprintf(messages.SendMessageSaved,
		msgsUrl, u.user.SendMessageCooldownLeft())); err != nil {
		return err
	}
	u.user.SetWant(BotWantCommand)
	return nil
}

// send алиас отправки текстового сообщения
func (u *Update) send(ctx context.Context, message string) error {
	if _, err := u.bot.do(ctx, params.NewSendMessage(
		u.Message.Chat.Id,
		message,
	)); err != nil {
		return err
	}
	return nil
}
