package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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
