package selenium

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
)

const (
	webDriverPort   = 4444
	chromeDebugPort = 9222
)

var (
	devMode = "n" // set during the build using ldflags (eg.: go build -ldflags="-X 'parser/internal/infrastructure/selenium.devMode=y'")

	webDriverCaps = selenium.Capabilities{
		"browserName": "chrome",
		"platform":    "Windows 11",
	}

	chromeCaps = chrome.Capabilities{
		Path: "",
		Args: []string{
			"--blink-settings=imagesEnabled=false",
			"--disable-dev-shm-usage",
			"--disable-gpu",
			"--no-sandbox",
			"--disable-infobars",
			"--disable-extensions",
			"--disable-popup-blocking",
			"--ignore-certificate-errors",
			"--disable-web-security",
			"--allow-running-insecure-content",
			"--disable-blink-features=AutomationControlled",
			"--disable-breakpad",
			"--disable-client-side-phishing-detection",
			"--disable-default-apps",
			"--disable-translate",
			"--disable-password-manager-reauthentication",
			"--disable-save-password-bubble",
			"--disable-single-click-autofill",
			"--disable-sync",
			"--no-default-browser-check",
			"--no-first-run",
			"--use-mock-keychain",
			"--lang=ru-RU",
			"--disable-background-timer-throttling",
			"--disable-background-networking",
			"--disable-backgrounding-occluded-windows",
			"--disable-popup-blocking",
			"--disable-notifications",
			"--disable-component-extensions-with-background-pages",
			"--disable-ipc-flooding-protection",
			"--disable-renderer-backgrounding",
			fmt.Sprintf("--remote-debugging-port=%d", chromeDebugPort),
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.170 Safari/537.36",
		},
		W3C: true,
	}
)

type Chrome struct {
	service     *selenium.Service
	webDrivers  []selenium.WebDriver
	seleniumUrl string
}

type ServerOptions struct {
	Host string
}

func NewChrome(opts *ServerOptions, ctx context.Context) (*Chrome, error) {
	result := &Chrome{webDrivers: make([]selenium.WebDriver, 0)}

	if "" != opts.Host {
		result.seleniumUrl = fmt.Sprintf("http://%s:%d/wd/hub", opts.Host, webDriverPort)
	}

	if err := result.startService(ctx); err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			for _, wd := range result.webDrivers {
				wd.Quit()
			}
			log.Println("Stop chrome drivers on context done")
		}
	}()

	return result, nil
}

func (c *Chrome) NewWindow(url string) (selenium.WebDriver, error) {
	wd, err := c.init()
	if err != nil {
		return wd, err
	}

	c.webDrivers = append(c.webDrivers, wd)

	return wd, wd.Get(url)
}

func (c *Chrome) init() (selenium.WebDriver, error) {
	chromeCaps := chromeCaps

	if "n" == devMode {
		chromeCaps.Args = append(chromeCaps.Args, "--headless=true") // Без графического интеррфейса (or --headless=new)
	}

	webDriverCaps := webDriverCaps
	webDriverCaps.AddChrome(chromeCaps)

	webDriver, err := selenium.NewRemote(webDriverCaps, c.seleniumUrl)
	if err != nil {
		return nil, err
	}

	return webDriver, nil
}

func (c *Chrome) startService(ctx context.Context) error {
	if nil != c.service || "" != c.seleniumUrl {
		return nil
	}

	service, err := selenium.NewChromeDriverService("chromedriver", webDriverPort)
	if err != nil {
		return err
	}

	c.service = service

	go func() {
		select {
		case <-ctx.Done():
			service.Stop()
			log.Println("Stop chrome driver service on context done")
		}
	}()

	return nil
}
