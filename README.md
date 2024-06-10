# Android TV Controller Telegram Bot

## Overview

This is a Go project that allows you to control your Android TV using a Telegram bot. With this bot, you can control your Chromecast / TV just as you do it with its remote control!

## Getting Started

### Prerequisites

- Telegram bot tken (you can get it from [BotFather](https://telegram.me/BotFather))
- Docker installed on your machine (should be on the same network as TV is)
- Android TV / Google TV / Chromecast with Android version >= 5

### Installation

```bash
docker run --rm --add-host=tv:<TV_IP_ADDRESS> --name telegram-tv-remote -e TELEGRAM_APITOKEN=<TELEGRAM_BOT_TOKEN> aminsaedi/t-bot
```

## Usage

The first time that you `/start` the bot, you have to send `/pair` to pair the bot with your tv, you will receive a verification code on your tv and you have to send to the bot.
After a successful pair, you will be able to use `/status` and `/remote`. Enjoy :)
