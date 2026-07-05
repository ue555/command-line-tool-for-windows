package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// 引数チェック
	args := os.Args
	if len(args) < 2 {
		fmt.Println("使い方: touch <path>")
		os.Exit(1)
	}

	path := args[1] // 引数のみ取得する
	_, err := os.Stat(path)
	if err == nil {
		err := os.Chtimes(path, time.Now(), time.Now())

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	} else if !os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer f.Close()
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}
