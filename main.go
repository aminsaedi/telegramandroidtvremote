package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	v2c "github.com/aminsaedi/atvremote/pkg/v2/command"
	v2com "github.com/drosocode/atvremote/pkg/common"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\x8B", "POWER"),
		tgbotapi.NewInlineKeyboardButtonData("\xE2\xAC\x86", "UP"),
		tgbotapi.NewInlineKeyboardButtonData("\xE2\x86\xA9", "BACK"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("\xE2\xAC\x85", "LEFT"),
		tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\xB5", "ENTER"),
		tgbotapi.NewInlineKeyboardButtonData("\xE2\x9E\xA1", "RIGHT"),
	),

	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\x89", "VOLUME_DOWN"),
		tgbotapi.NewInlineKeyboardButtonData("\xE2\xAC\x87", "DOWN"),
		tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\x8A", "VOLUME_UP"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\x87", "MUTE"),
		tgbotapi.NewInlineKeyboardButtonData("Play/Pause", "PLAYPAUSE"),
	),
)

var (
	certs, err = tls.LoadX509KeyPair("cert.pem", "key.pem")
	cmd        = v2c.New("tv", 6466, &certs)
)

func connect() {
	err = cmd.Connect()
	if err != nil {
		log.Fatalf("unable to connect to TV: %s", err)
	}
}

func sendKey(k v2com.RemoteKeyCode) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %s", r)
			log.Printf("Reconnecting to TV")
			connect()
			time.Sleep(1 * time.Second)
			log.Printf("Re-Sending key")
			sendKey(k)
		}
	}()

	log.Printf("Sending key: %d", k)
	err = cmd.SendKey(k)
	if err != nil {
		log.Panicf("unable to send key: %s", err)
	}
}

func getStatus() (result string) {
	result = "Unable to get status\nTry /connect first"
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %s", r)
		}
	}()

	data, _ := cmd.GetData()
	bdata, _ := json.MarshalIndent(data, "", "  ")
	result = "Status: \n```json\n" + string(bdata) + "\n```\n"
	return result
}

var commands = [...]tgbotapi.BotCommand{
	{
		Command:     "start",
		Description: "Start the bot",
	},
	{
		Command:     "remote",
		Description: "Enter remote control mode",
	},
	{
		Command:     "pair",
		Description: "Pair with TV",
	},
	{
		Command:     "connect",
		Description: "Connect to TV",
	},
	{
		Command:     "status",
		Description: "Get TV status",
	},
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil && update.Message.IsCommand() {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Command() {
			case "start":
				msg.Text = "Welcome to ATV Remote Bot"
				cmds := tgbotapi.NewSetMyCommands(commands[:]...)
				bot.Send(cmds)
			case "status":
				msg.Text = getStatus()
				msg.ParseMode = "MarkdownV2"
			case "remote":
				msg.Text = "Remote control mode, use /exit to exit remote control mode"
				msg.ReplyMarkup = numericKeyboard
			case "pair":
				msg.Text = "Pairing request received"
			case "connect":
				connect()
				msg.Text = "Connected to TV"
			case "exit":
				msg.Text = "Exiting remote control mode"
			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			// msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			// if _, err := bot.Send(msg); err != nil {
			// 	panic(err)
			// }

			switch update.CallbackQuery.Data {
			case "POWER":
				sendKey(26)
			case "UP":
				sendKey(19)
			case "DOWN":
				sendKey(20)
			case "LEFT":
				sendKey(21)
			case "RIGHT":
				sendKey(22)
			case "ENTER":
				sendKey(23)
			case "VOLUME_UP":
				sendKey(24)
			case "VOLUME_DOWN":
				sendKey(25)
			case "PLAYPAUSE":
				sendKey(85)
			case "MUTE":
				sendKey(91)
			default:
				log.Printf("Unknown command: %s", update.CallbackQuery.Data)
			}

		}
	}
}
