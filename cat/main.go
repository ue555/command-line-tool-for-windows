package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// 1. コマンドライン引数を取得する
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("ファイルパスを指定してください")
		return
	}

	// 2. ファイルを開く
	filePath := args[0]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("ファイルを開けません: %v\n", err)
		return
	}
	// main関数が終わる前に必ずファイルを閉じる
	defer file.Close()

	// 3. ファイルの内容を標準出力(os.Stdout)に書き出す
	// io.Copy を使うと、ファイルを一行ずつ処理する手間を省いて一気に書き出せます
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		fmt.Printf("読み込み/書き出しエラー: %v\n", err)
	}
}
