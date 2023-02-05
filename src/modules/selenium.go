package modules

import (
	"github.com/tebeka/selenium"
)

type SeleniumConfig struct {
	Capabilities selenium.Capabilities
	SeleniumHub  string
}

func NewSelenium(config SeleniumConfig) (error, *selenium.WebDriver) {
	// Connect to the WebDriver instance running locally.
	wd, err := selenium.NewRemote(config.Capabilities, config.SeleniumHub)
	if err != nil {
		return err, nil
	}
	return nil, &wd
}
