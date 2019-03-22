# go_twitch
Web client on Go for watching Twitch streams. 

It's currently running on Heroku, so you can play around without any installation. Simply open in browser: http://stark-harbor-28675.herokuapp.com

## Setup and Run
Client utilizes Twitch API, therefore in order to get it running, you need to first:
1. Create and login to Twitch account
2. Register a [new App](https://dev.twitch.tv/console/apps) on Twitch Developers site.
* *Owner*: name for your app
* *OAuth Redirect URL*: this must be your client instance URL where Twitch will redirects once it got consent from user, e.g. "https://stark-harbor-28675.herokuapp.com:80" for Heroku instance or "http://localhost:80" if you want to run it locally.
* *Category*: whatever

Copy your *"Client identifier"* and *"Client secret key"* (hit "New secret key" button to get it) and save it somewhere. This will be your Twitch client login and password :-)


### Now you can run client itself. 
Following environments variables must be set to set up client:

TWITCH_CLIENT_ID - Mandatory. Use *"Client identifier"* copied before as value.
TWITCH_CLIENT_SECRET - Mandatory. Use *"Client secret key"* copied before as value.
PORT - Optional. By default, if not provided, server will try to bind listener on port 8080. Heroku will automatically set PORT env var for you.
HOST - Optional. It's used for creating Redirect URL inside client. By default, if not provided, client will try to guess your hostname, however it doesn't work well in case of Heroku instance. So, it's good to set this env var beforeahead.

#### How to set variables?
*On Linux*:
either use `export key=value` format command in your shell or set them before start server command:
```bash
➜  : TWITCH_CLIENT_ID="ghhxr1234567890sqhz" TWITCH_CLIENT_SECRET=owhowqheqwheoq HOST=localhost go run cmd/main.go
2019-03-22T16:53:36.054+0300    INFO    cmd/main.go:22  Host name is set        {"Host": "localhost"}
2019-03-22T16:53:36.054+0300    INFO    cmd/main.go:28  Listening Port is set   {"Port": "8080"}
2019-03-22T16:53:36.054+0300    INFO    cmd/main.go:34  Client ID is set, OK    {"ID": "ghhxr1234567890sqhz"}
2019-03-22T16:53:36.054+0300    INFO    cmd/main.go:40  Client Secret is set, OK
2019-03-22T16:53:36.054+0300    INFO    api/router.go:39        HTTP service started    {"Incoming": "http://localhost:80", "Outgoing": "http://localhost:8080"}
2019-03-22T16:53:40.940+0300    INFO    api/auth.go:44  Inside ShowLoginPage() view function    {"Method": "GET"}
```

*On Heroku*
```bash
heroku config:set HOST=stark-harbor-28675.herokuapp.com
heroku config:set TWITCH_CLIENT_ID=foo1234
heroku config:set TWITCH_CLIENT_SECRET=bar5678
```

#### How to run?
*Locally* (it's assumed you'e already set env vars):

`go run cmd/main.go`

*Heroku*:

It's currently running on my Heroku instance, so you can play around without any installation. Simply open in browser: http://stark-harbor-28675.herokuapp.com

If you want to run it on your own Heroku instance, simply follow this guideline:
https://devcenter.heroku.com/articles/getting-started-with-go#deploy-the-app

## Known issues
Once you get logged in to Twitch, your session keeps alive for around one hour. There is no automatic refresh for your auth token, so (at least locally) you have to manually clean cookies in order to get a new one. I will fix it if I find the time, also feel free to send PR with fix, I'll be more than happy accept it.
