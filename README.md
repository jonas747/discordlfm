#discordlfm - Discord Last FM 

Sets your "game" to the last song played on your last.fm

Usage:
```
  -g string
        Game to set to if there hasn't been a new song for a while (default "", clears playing status)
  -l string
        Last.fm api key
  -n int
        Number of seconds without a new song for it to be considered nothing. (default 600)
  -t string
        Discord token
  -u string
        Last.fm username
```

Example: `./discordlfm -t "my secret discord token" -l "my secret last.fm api key" -u "jonasr747"`

Guide on retrieving the api lfm api key and discord token: https://github.com/jonas747/discordlfm/issues/2#issuecomment-278950579
