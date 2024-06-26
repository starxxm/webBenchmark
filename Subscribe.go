package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Subscribe(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()
	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(readAll)
}

func subscribeUpdate(url string) {
	defer func() {
		if r := recover(); r != nil {
			go subscribeUpdate(url)
		}
	}()
	for {
		time.Sleep(30 * time.Second)
		subUrl := Subscribe(url)
		if len(subUrl) > 0 {
			if *detectLocation {
				location := GetHttpLocation(subUrl)
				if len(location) > 0 {
					TargetUrl = location
				} else {
					TargetUrl = subUrl
				}
			} else {
				TargetUrl = subUrl
			}
		}
	}
}
