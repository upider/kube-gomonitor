package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

func busy(stopChan chan os.Signal) {
	for {
		select {
		case <-stopChan:
			return
		default:
			fmt.Println("busy loop is running...")
			cnt := 0
			for i := 0; i < 1000000; i++ {
				cnt++
			}
			time.Sleep(500 * time.Millisecond)
			//Init jar
			j, _ := cookiejar.New(nil)
			// Create client
			client := &http.Client{Jar: j}
			// Create request
			req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
			// Fetch Request
			_, err := client.Do(req)
			if err != nil {
				fmt.Println("Failure : ", err)
			}
			//开始修改缓存jar里面的值
			var clist []*http.Cookie
			clist = append(clist, &http.Cookie{
				Name:    "BDUSS",
				Domain:  ".baidu.com",
				Path:    "/",
				Value:   "Dc2cG5McjNzZlJxMi00SHM4eWJxRWp3elpiT0hoVEhWYjJCTWh6dDIxc2pYODlaSVFBQUFBJCQAAAAAAAAAAAEAAABSgP0BQUxBTE1OwLa6~AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACPSp1kj0qdZU",
				Expires: time.Now().AddDate(1, 0, 0),
			})
			urlX, _ := url.Parse("http://zhanzhang.baidu.com")
			j.SetCookies(urlX, clist)

			fmt.Printf("Jar cookie : %v", j.Cookies(urlX))
			// Fetch Request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Failure : ", err)
			}

			// Read Response Body
			respBody, _ := ioutil.ReadAll(resp.Body)

			// Display Results
			fmt.Println("response Status : ", resp.Status)
			fmt.Println("response Headers : ", resp.Header)
			fmt.Println("response Body : ", string(respBody))
			fmt.Printf("response Cookies :%v", resp.Cookies())
		}
	}
}
