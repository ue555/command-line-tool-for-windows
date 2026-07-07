package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	// 引数チェック
	args := os.Args
	if len(args) < 2 {
		fmt.Println("使い方: stat <filename>")
		os.Exit(1)
	}

	path := args[1] // 引数のみ取得する
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("Name: %s\n", info.Name())
		fmt.Printf("Size: %d\n", info.Size())
		fmt.Printf("ModTime: %s\n", info.ModTime())
		fmt.Printf("IsDir: %v\n", info.IsDir())
		fmt.Printf("Mode: %s\n", info.Mode())

		// info.Sys() は any 型なので、Windows固有の詳細情報を得るために
		// *syscall.Win32FileAttributeData 型へ型アサーションする
		// 第2戻り値の ok で成功/失敗を安全に判定できる（失敗してもパニックしない）
		sys, ok := info.Sys().(*syscall.Win32FileAttributeData)
		if !ok {
			// アサーション失敗時（Windows以外の環境や想定外の型の場合）はエラーを出して終了
			fmt.Fprintln(os.Stderr, "failed to get Windows file attributes")
			return
		}

		// CreationTime は syscall.Filetime 型（Windows FILETIME形式：1601年基準・100ナノ秒単位）
		// Nanoseconds() で Unixエポック基準のナノ秒に変換し、
		// time.Unix(0, ns) で time.Time 型に変換する
		creationTime := time.Unix(0, sys.CreationTime.Nanoseconds())

		// 最終アクセス日時も同様に time.Time 型へ変換
		lastAccessTime := time.Unix(0, sys.LastAccessTime.Nanoseconds())

		// 最終更新日時も同様に time.Time 型へ変換
		lastWriteTime := time.Unix(0, sys.LastWriteTime.Nanoseconds())

		// time.Time は Stringer を実装しているため %s で読みやすい日時文字列として出力される
		fmt.Printf("CreationTime:   %s\n", creationTime)
		fmt.Printf("LastAccessTime: %s\n", lastAccessTime)
		fmt.Printf("LastWriteTime:  %s\n", lastWriteTime)

		// FileAttributes は複数の属性情報を1つの uint32 にビットフラグとして詰め込んだ値
		// %x で16進数表示し、生のビットパターンを確認できるようにする
		fmt.Printf("FileAttributes: 0x%x\n", sys.FileAttributes)

		// FileAttributes と各属性定数（1ビットのみ立っている値）をビットANDすることで
		// 該当ビットが立っているかどうかを判定し、!= 0 で bool値に変換している

		// 読み取り専用フラグが立っているか判定
		fmt.Printf("  ReadOnly:  %v\n", sys.FileAttributes&syscall.FILE_ATTRIBUTE_READONLY != 0)

		// 隠しファイルフラグが立っているか判定
		fmt.Printf("  Hidden:    %v\n", sys.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0)

		// システムファイルフラグが立っているか判定
		fmt.Printf("  System:    %v\n", sys.FileAttributes&syscall.FILE_ATTRIBUTE_SYSTEM != 0)

		// ディレクトリフラグが立っているか判定
		fmt.Printf("  Directory: %v\n", sys.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0)

		// アーカイブ（バックアップ対象マーク）フラグが立っているか判定
		fmt.Printf("  Archive:   %v\n", sys.FileAttributes&syscall.FILE_ATTRIBUTE_ARCHIVE != 0)
	}
}
