package helper

import (
	seleniumBase "github.com/tebeka/selenium"
)

func GetText(we seleniumBase.WebElement, element string) (string, error) {
	if "" == element {
		return we.Text()
	}

	el, err := we.FindElement(seleniumBase.ByCSSSelector, element)
	if err != nil {
		return "", err
	}

	return el.Text()
}

func GetHref(we seleniumBase.WebElement, element string) (string, error) {
	if "" == element {
		return we.GetAttribute("href")
	}

	el, err := we.FindElement(seleniumBase.ByCSSSelector, element)
	if err != nil {
		return "", err
	}

	return el.GetAttribute("href")
}
