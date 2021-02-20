# snowfall-discord
This is a way to verify users based on their discord handle, with no need for email, or phone number.

## How it works
When a user registers on a website, the bot receives the discord username and denominator, 
if the given user was found within a given server, it then tries to send a direct message with the activation string.

## Usage
Since the code relies heavily on environment variables, it's important to set them either before or along with the code.
#### Example launch command using bash
```
go build -o snowfall-discord
BOT_TOKEN=asdfgh12345 RABBIT_URL=amqp://guest:guest@localhost:5672/ GUILD_ID=1234567890987654321 CHANNEL_ID=1234567890987654321 ./snowfall-discord
```
Make sure that your bot has permissions to list users on servers it's on, as well as message permissions in its respective channel.
