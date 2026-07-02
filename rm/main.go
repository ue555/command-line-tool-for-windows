package main

import (
	"fmt"
	"os"
)

func main() {
	// 引数チェック
	args := os.Args
	if len(args) < 2 {
		fmt.Println("使い方: rm <path> [path2 ...]")
		os.Exit(1)
	}

	// 複数ファイルをループ処理
	for _, path := range args[1:] {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println("ファイルが見つかりません:", path)
			continue
		}

		if info.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				fmt.Println("削除失敗:", err)
			}
		} else {
			// ファイルの場合
			if err := os.Remove(path); err != nil {
				fmt.Println("削除失敗:", err)
			}
		}
	}
}
