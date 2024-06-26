package main

import (
	"fmt"
	"net/http"
	"time"
)

func GetHttpLocation(Url string) string {
	resp, err := http.Get(Url)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("Headers for", url)

	if resp.Request == nil || resp.Request.Response == nil || resp.Request.Response.Header == nil {
		return ""
	}
	locationHeader := resp.Request.Response.Header["Location"][0]
	if len(locationHeader) > 0 {
		return locationHeader
	}
	return ""
}

func RefreshHttpLocation(Url string) string {
	defer func() {
		if r := recover(); r != nil {
			go RefreshHttpLocation(*url)
		}
	}()

	for {
		time.Sleep(60 * time.Second)
		location := GetHttpLocation(Url)
		if len(location) > 0 {
			TargetUrl = location
		} else {
			TargetUrl = *url
		}
	}
}
