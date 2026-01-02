package messages

// messages
const (
	Start = `😎 Hello %s.
Welcome to %s telegram bot.
You can send messages after you confirm that you are a human using the %s command.
After that, using the %s command, you can send a message that will be saved.
You can add no more than one message per day.
You can see your message on address %s`

	Commands = `🙄 The bot supports the following commands:
%s to pass the verification that you are a human.
%s use the command to send messages.`

	AuthTitle = `😀 Please enter symbols from picture.
You have %d times left.`

	AuthWrong = `😔 Wrong letters.`

	AuthComplete = `👍 Auth complete.`

	AuthIncomplete = `😌 Auth incomplete.
You need send %s and enter letters from the image correctly.`

	IsEmpty = `You can only send text messages.`

	AuthOver = `🤭 Auth attempts ended, try through %s.`

	AlreadyAuth = `😎 You already authorized.`

	SendMessageTitle = `😃 Now you can send %s command to enter text message.`

	SendMessageEnter = `🙂 Please enter a text message`

	SendMessageSaved = `🙃 Thank you. Message are saved.
You can see your message on address %s
The next message can be sent no sooner than %s`

	WillSaved = `🧐 You have already sent a message.
You can see your message on address %s
The next message can be sent no sooner than %s`

	UnknownCommand = `🤔 Unknown command %s`
)
