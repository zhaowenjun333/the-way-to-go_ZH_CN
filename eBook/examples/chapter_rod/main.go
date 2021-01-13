package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var locChina, _ = time.LoadLocation("Asia/Shanghai")
var lock = &sync.Mutex{}

// newBrowser 浏览器初始化，同时设置用户数据目录
func newBrowser(headless bool) *rod.Browser {
	url := launcher.New().Headless(headless).UserDataDir("tmp/user").MustLaunch()
	return rod.New().ControlURL(url).MustConnect()
}

// login 登录淘宝，自动选择扫码模式
func login() {
	lock.Lock()
	defer lock.Unlock()

	browser := newBrowser(true)
	defer browser.Close()

	// 检查页面 a标签是否包含 待收货/待付款/待发货/待评价
	loginPage := browser.MustPage("https://login.taobao.com/member/login.jhtml").MustWaitLoad()
	loginPage.MustElement(".icon-qrcode").MustClick()

	// 重试，直到找到指定的选择器中包含指定的字符
	loginPage.Race().ElementR("a", "待评价").
		MustHandle(func(e *rod.Element) {
			log.Println("登录成功，已经找到登录后页面元素内容 ", e.MustText())
		}).MustDo()

}

// isLoggedIn 检查是否登录，使用无头模式
func isLoggedIn() bool {
	lock.Lock()
	defer lock.Unlock()

	browser := newBrowser(true)
	defer browser.Close()

	islog := browser.MustPage("https://taobao.com/").MustWaitLoad().MustHasR("a", "免费注册")
	return !islog
}

func main() {

	// 检查是否已经登录
	if !isLoggedIn() {
		login()
	}

	// 搜索指定关键词
	keyword := "打火机"
	searchURL := fmt.Sprintf("http://s.taobao.com/search?q=%s", keyword)
	log.Println(searchURL)

	// 获取指定关键词的页面宝贝名称
	browser := newBrowser(true).Timeout(time.Second * 30)
	defer browser.Close()

	// searchPage := browser.MustPage(searchURL).MustWaitLoad()
	searchPage := browser.MustPage(searchURL)

	searchPage.MustScreenshot("search.png") // 根据截图查看，页面未登录。

	// items := searchPage.MustElements(".J_MouserOnverReq")
	// items := searchPage.MustElements("m-itemlist")

	// for _, item := range items {
	// 	title := item.MustElement("a")
	// 	log.Println(title.MustBlur().Text())
	// }

}
