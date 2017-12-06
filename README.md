# gotalks
gotalks is a small chat web app + server to play around with websockets.

## Configuration
Example configuration is located at `config/config.toml.dist`
Copy it to `config/config.toml` and provide proper config values.

### Google API
This app uses google as OAuth2 provider in order to use it you need to register an app in google's developer console and set the `client_id` and `client_secret` in the cofig file

## Building the app
For the depencies management [golang/dep](https://github.com/golang/dep) is used.
Simply run:
```
$ dep ensure
```
to fetch all the dependencies.

Now just run 
```
$ go build -o bin/chat
```
to build the app

## Run
To run the app:
```
$ ./bin/chat
```
