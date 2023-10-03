package helper

import (
	seleniumBase "github.com/tebeka/selenium"
	"log"
	"time"
)

func TrickProtection(wd seleniumBase.WebDriver) {
	log.Println("Try trick protection")

	time.Sleep(time.Second * 10)

	if _, err := wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 1)

	if _, err := wd.ExecuteScript("window.scrollTo(document.body.scrollHeight, 0);", nil); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 2)
}
