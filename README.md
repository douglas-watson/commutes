# Strava Commute Uploader
- Douglas Watson
- March 2016

I commute the same road to work everyday, and I like to use Strava to keep track of how much distance I ride. To simplify the process I wrote this small program that creates a private manual entry on Strava for every day within a specified time range. Weekends and holidays can be excluded.


## Usage
I wrote this entirely for my own use, and thus did not include many configuration options. If you want to adapt the code to your own needs, feel free to fork!

First of all, copy `config-template.yaml` to `config.yaml`.

### Get an API access token

Create an application key on Strava, from your user settings page: https://www.strava.com/settings/api. Then, run Strava's official oath_example.go to get a token with write privilege. This requires a couple steps:

1. Install the official go.strava client library:

 	​ ```
	 go get github.com/strava/go.strava
	 ​```

2. Find oath_example.go and open it with your favourite editor:
  ​```
  cd $GOPATH/src/github.com/strava/go.strava/examples
  vi oath_example.go
  ​```
3. Modify the authorization request to ask for Write privileges. Change the following line in `indexHandler` from this:
   ```
   fmt.Fprintf(w, `<a href="%s">`, 	authenticator.AuthorizationURL("state1", strava.Permissions.Public, true))
   ```
   to this:
   ```
   fmt.Fprintf(w, `<a href="%s">`, authenticator.AuthorizationURL("state1", strava.Permissions.Write, true))
   ```

4. Run the example, using the ID and secret from the Strava settings page:
   ```
   go run oath_example.go -id=YOUR_ID -secret=YOUR_SECRET
   ```
   Follow the instructions, authorize the app from your browser, and obtain an access token. Copy this access token to your `config.yaml` file. On the same result page as the access token, find the id of your commuter bike, and copy it as the "Gear ID" value in `config.yaml`. The bike entry looks like this:

  ```
  ...
  "bikes": [
        ...
    {
     "id": "XXXX",
     "name": "Commuter bike",
     "primary": false,
     "distance": 3.695307e+06
    },
        ...
  ]
  ...
  ```

### Fill in commute details

Continue filling in `config.yaml`. Specify the commute distance and time, start and end date of the upload, and excluded periods (vacation days or days you did not bike). You must include at least one excluded date range, but it can be outside of the upload range.

### Upload your commutes!

Compile and run the program, with the config file you just filled in:

```
go build && ./commutes config.yaml
```

For every day in the date range, this will create a private manual entry on Strava named "Commute", with the specified distance and time. The uploader skips days in the each of the "Exclude dates" ranges, and optionally on weekends.
