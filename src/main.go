package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
	"go-selenium.com/component/db"
	"go-selenium.com/modules"
	"log"
	"os"
	"time"
)

func main() {
	// load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	appId := os.Getenv("IG_APP_ID")
	seleniumConfig := modules.SeleniumConfig{
		Capabilities: selenium.Capabilities{"browserName": os.Getenv("BROWSER_NAME")},
		SeleniumHub:  os.Getenv("SELENIUM_HUB"),
	}

	instagramAcc := modules.InstagramAcc{
		LoginUrl: "https://www.instagram.com",
		UserName: os.Getenv("IG_USER_NAME"),
		PassWord: os.Getenv("IG_PASSWORD"),
	}

	listPage := []modules.PublicPage{
		{PageName: "xxx"},
	}

	fmt.Println("Register database")
	rootDir, _ := os.Getwd()
	err, database := db.Connection(fmt.Sprintf("%s/component/db/database.db", rootDir))
	if err != nil {
		panic(err)
	}
	if err = db.CreateTable(database); err != nil {
		panic(err)
	}

	fmt.Println("connected selenium hub")
	err, instance := modules.NewSelenium(seleniumConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println("Capture cookie")
	err, cookie := modules.Login(instance, instagramAcc)
	if err != nil {
		panic(err)
	}
	authInfo := modules.AuthInfo{
		AppId:     appId,
		CookieMap: cookie,
	}

	//get page follower
	err = modules.Execute(listPage, authInfo, database)

	fmt.Println("🚀. Get list follower of page done..")
	fmt.Println("🔥. Hang on~ Send follow request...")

	// Alternative Version
	for true {
		if err := modules.SendFollowRequestManual(database, instance); err != nil {
			panic(err)
			return
		}
		time.Sleep(time.Second)
	}

	// close webDriver
	var wd selenium.WebDriver = *instance
	defer func(wd selenium.WebDriver) {
		err = wd.Quit()
		if err != nil {
			panic(err)
		}
	}(wd)
	return
}
