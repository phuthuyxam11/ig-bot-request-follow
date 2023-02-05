package modules

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/tebeka/selenium"
	"go-selenium.com/component/db"
	"time"
)

func SendFollowRequestManual(store *sql.DB, seleniumInstance *selenium.WebDriver) error {
	sendList, err := db.GetAllUser(store)

	if err != nil {
		return err
	}
	if len(sendList) == 0 {
		return errors.New("empty data finish loop")
	}
	bar := progressbar.Default(int64(len(sendList)))

	// chunk data
	sendListChunk := ArrayChunk[db.UserModel](20, sendList)
	for _, list := range sendListChunk {
		_ = bar.Add(1)
		if err := ExeSendFollowManual(list, seleniumInstance); err != nil {
			continue
		}
	}
	return nil
}

func ExeSendFollowManual(listUser []db.UserModel, seleniumInstance *selenium.WebDriver) error {
	var wd selenium.WebDriver = *seleniumInstance
	for _, user := range listUser {
		userPage := fmt.Sprintf("https://www.instagram.com/%s", user.UserName)
		err := wd.Get(userPage)
		if err != nil {
			return err
		}
		if err = wd.Wait(waitUserPageIsDisplay); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
		// click follower
		btnFollow, err := wd.FindElement(selenium.ByXPATH, "//button[@class='_acan _acap _acas _aj1-']")
		err = btnFollow.Click()
		if err != nil {
			fmt.Println("❌: Error send request follow : " + user.UserName)
			return err
		}
		time.Sleep(time.Second * 1)
		fmt.Println("✅: Following : " + user.UserName)
	}

	return nil
}

var waitUserPageIsDisplay = func(wd selenium.WebDriver) (bool, error) {

	commentElm, err := wd.FindElement(selenium.ByXPATH, "//button[@class='_acan _acap _acas _aj1-']")
	if err != nil {
		return false, err
	}
	e, ok := commentElm.IsDisplayed()
	if ok != nil {
		panic(ok)
	}
	return e, nil

}

func ArrayChunk[T any](chunkSize int, arr []T) [][]T {
	var divided [][]T

	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize

		if end > len(arr) {
			end = len(arr)
		}

		divided = append(divided, arr[i:end])
	}
	return divided
}
