# go-twitter-bot
A simple Twitter bot written in Go

## Setup

1. Create an app for your Twitter account at https://apps.twitter.com and generate an AccessToken and Secret. Use these to populate config.json.
2. Put your own tweets in tweets.json.
3. Run `go build && ./bot -config "config.json" -data "tweets.json"`.

## To Run on a Linux Server

Follow steps 1 & 2 above, then:

1. Build for Linux with `env GOOS=linux GOARCH=amd64 go build`.
2. Copy the binary, `config.json`, and `tweets.json` to your server.
3. Place the `go-bot.conf` file in the `/etc/init` directory. Modify the install dir as required.
4. (Using upstart init daemon that ships with Ubuntu) Run the cmd `start go-bot`. Check it's running with `ps -ef`. Stop the process with `stop go-bot`.
