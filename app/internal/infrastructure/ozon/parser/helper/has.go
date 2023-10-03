package helper

import (
	seleniumBase "github.com/tebeka/selenium"
	"strings"
)

func HasElement(wd seleniumBase.WebDriver, element string) (bool, error) {
	els, err := wd.FindElements(seleniumBase.ByCSSSelector, element)
	if err != nil {
		return false, err
	}

	if len(els) > 0 {
		return true, nil
	}

	return false, nil
}

func HasText(we seleniumBase.WebElement, element string, findText string) (bool, error) {
	els, err := we.FindElements(seleniumBase.ByCSSSelector, element)
	if err != nil {
		return false, err
	}

	for _, el := range els {
		text, err := el.Text()
		if err != nil {
			return false, err
		}

		if strings.ToLower(text) == strings.ToLower(findText) {
			return true, nil
		}
	}

	return false, nil
}
