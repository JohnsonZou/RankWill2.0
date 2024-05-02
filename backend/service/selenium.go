package service

import (
	"github.com/tebeka/selenium"
	"log"
	"time"
)

func Selenium() {

	caps := selenium.Capabilities{
		"browserName": "MicrosoftEdge",
		"goog:chromeOptions": map[string]interface{}{
			"args": []string{"--headless", "--disable-gpu"},
		},
	}

	wd, err := selenium.NewRemote(caps, "http://localhost:9515")
	if err != nil {
		log.Println("Failed to open session:", err)
		return
	}
	defer wd.Quit()

	err = wd.Get("https://leetcode.cn/contest/api/ranking/weekly-contest-385/?pagination=1&region=global")
	if err != nil {
		log.Println("Failed to navigate:", err)
		return
	}

	// 等待一些时间以确保页面加载完成
	time.Sleep(2 * time.Second)

	pageSource, err := wd.ExecuteScript("return document.documentElement.outerHTML", nil)
	if err != nil {
		log.Fatal("Failed to get page source:", err)
		return
	}

	// 将结果转为字符串
	pageSourceStr, ok := pageSource.(string)
	if !ok {
		log.Fatal("Failed to convert page source to string")
		return
	}

	log.Println("页面 HTML 内容:", pageSourceStr)
}
