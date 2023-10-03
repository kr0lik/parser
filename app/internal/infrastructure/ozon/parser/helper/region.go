package helper

import (
	"fmt"
	seleniumBase "github.com/tebeka/selenium"
	"log"
	"parser/internal/config/ozon"
	"strings"
	"time"
)

type Region struct {
	Config ozon.Region
}

func (p *Region) RegionCheck(wd seleniumBase.WebDriver) error {
	log.Println("Region check")

	if err := wd.WaitWithTimeout(func(wd seleniumBase.WebDriver) (bool, error) {
		return HasElement(wd, p.Config.CheckRegion.Wait.Element)
	}, time.Duration(p.Config.CheckRegion.Wait.Timeout)*time.Second); err != nil {
		return err
	}

	checkRegionHolder, err := wd.FindElement(seleniumBase.ByCSSSelector, p.Config.CheckRegion.Check.Element)
	if err != nil {
		return err
	}

	expectedRegion, err := checkRegionHolder.Text()
	if err != nil {
		return err
	}

	if strings.ToLower(expectedRegion) == strings.ToLower(p.Config.CheckRegion.Check.IsName) {
		log.Println("Region is correct")
		return nil
	}

	//log.Println("choose region")
	//return p.Choose(wd)

	return nil
}

func (p *Region) Choose(wd seleniumBase.WebDriver) error {
	regionPopupButton, err := wd.FindElement(seleniumBase.ByCSSSelector, p.Config.CallChoosePopup.Click)
	if err != nil {
		return err
	}

	if err = regionPopupButton.Click(); err != nil {
		return err
	}

	if err := wd.WaitWithTimeout(func(wd seleniumBase.WebDriver) (bool, error) {
		return HasElement(wd, p.Config.CallChooser.Wait.Element)
	}, time.Duration(p.Config.CallChooser.Wait.Timeout)*time.Second); err != nil {
		return err
	}

	chooseRegionButton, err := wd.FindElement(seleniumBase.ByCSSSelector, p.Config.CallChooser.Click)
	if err != nil {
		return err
	}

	err = chooseRegionButton.Click()
	if err != nil {
		return err
	}

	if err := wd.WaitWithTimeout(func(wd seleniumBase.WebDriver) (bool, error) {
		return HasElement(wd, p.Config.ChooseRegion.Wait.Element)
	}, time.Duration(p.Config.ChooseRegion.Wait.Timeout)*time.Second); err != nil {
		return err
	}

	regions, err := wd.FindElements(seleniumBase.ByCSSSelector, p.Config.ChooseRegion.List.Element)
	if err != nil {
		return err
	}

	for _, region := range regions {
		nameText, err := GetText(region, p.Config.ChooseRegion.List.Name)
		if err != nil {
			return err
		}

		if strings.ToLower(nameText) == strings.ToLower(p.Config.CheckRegion.Check.IsName) {
			if err := region.Click(); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("region no choosen")
}
