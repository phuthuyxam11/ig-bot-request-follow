# ig-bot-request-follow
 A bot. Get users who are following a public page. And send them a follow-up request
## Contents

1. ### Create .env file
```go
    IG_APP_ID=xxxxxxxx #follow example: https://www.youtube.com/watch?v=izeYkVZydxQ
    BROWSER_NAME=xxxxxx #browser node on selenium ex: chrome, edge, firefox ...
    SELENIUM_HUB=xxxxxxxx #selenium hub url: http://127.0.0.1:4444/wd/hub
    IG_USER_NAME=xxxxxx #username of instagram
    IG_PASSWORD=xxxxxxx #password of instagram
```
2. ### Add public page name
    go to /src/main.go
```go
    listPage := []modules.PublicPage{
        {PageName: "page_name_1"},
        {PageName: "page_name_2"},
		...
        {PageName: "page_name_xxx"},
    }
```
3. ### Run bot
```
cd /src
go run main.go
```
