package modules

import (
	"fmt"
	"github.com/tebeka/selenium"
	"os"
	"time"
)

type InstagramAcc struct {
	LoginUrl string
	UserName string
	PassWord string
}

type AuthInfo struct {
	AppId     string
	CookieMap []selenium.Cookie
}

func NewAuthInfo(appId string, cookieMap []selenium.Cookie) *AuthInfo {
	return &AuthInfo{
		AppId:     appId,
		CookieMap: cookieMap,
	}
}

func Login(instance *selenium.WebDriver, instagramAcc InstagramAcc) (error, []selenium.Cookie) {

	fmt.Println("go to Instagram")
	var wd selenium.WebDriver = *instance
	err := wd.Get(instagramAcc.LoginUrl)
	if err != nil {
		return err, nil
	}

	err = wd.SetImplicitWaitTimeout(time.Second * 10)
	if err != nil {
		return err, nil
	}

	elementUserName, err := wd.FindElement(selenium.ByXPATH, "//input[@name='username']")
	if err != nil {
		return err, nil
	}

	err = elementUserName.SendKeys(instagramAcc.UserName)
	if err != nil {
		return err, nil
	}

	elementPassword, err := wd.FindElement(selenium.ByXPATH, "//input[@name='password']")
	if err != nil {
		return err, nil
	}
	err = elementPassword.SendKeys(instagramAcc.PassWord)
	if err != nil {
		return err, nil
	}

	btnLogin, err := wd.FindElement(selenium.ByCSSSelector, "button[type=submit]")
	err = btnLogin.Click()
	if err != nil {
		return err, nil
	}

	time.Sleep(time.Second * 15)

	cookie, err := wd.GetCookies()
	fmt.Println("login success")
	return nil, cookie
}

func saveCookiesFile(cookieMap []selenium.Cookie) error {
	var stringResult string
	for _, s := range cookieMap {
		convertString := fmt.Sprintf("%s=%s", s.Name, s.Value)
		stringResult += convertString + "; "
	}

	f, err := os.Create("cookie.txt")

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(stringResult)

	if err != nil {
		return err
	}
	return nil
}
