package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (dict *Dict) GetInfo() error {
	dict.IsUpdate = false
	url := fmt.Sprintf("https://pinyin.sogou.com/dict/detail/index/%d", dict.ID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	dict.Name = doc.Find(".dict_detail_title").Text()
	date := doc.Find(".dict_info_list > ul:nth-child(1) > li:nth-child(4) > div:nth-child(1)").Text()
	date, _ = strings.CutPrefix(date, "更   新：")
	if dict.Date != date {
		dict.Date = date
		dict.IsUpdate = true
	}

	return nil
}

func (dict *Dict) Download() error {
	f, err := os.Create(fmt.Sprintf("scel/%d-%s.scel", dict.ID, dict.Name))
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://pinyin.sogou.com/d/dict/download_cell.php?id=%d&name=%s&f=detail", dict.ID, dict.Name)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}
