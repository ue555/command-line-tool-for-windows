package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func main() {
	// 引数があればそのパス、なければカレントディレクトリ
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		os.Exit(1)
	}

	// 隠しファイルを除外（Windowsの隠し属性 + ドット始まり）
	var names []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if isHidden(path, e.Name()) {
			continue
		}
		names = append(names, e.Name())
	}

	// 名前でソート（大文字小文字を区別しない）
	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(names[i]) < strings.ToLower(names[j])
	})

	// ターミナル幅を取得（失敗したら80をデフォルト）
	termWidth := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		termWidth = w
	}

	// 最長ファイル名の長さを取得してカラム幅を決定
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}
	colWidth := maxLen + 2 // 2文字分の余白
	cols := termWidth / colWidth
	if cols < 1 {
		cols = 1
	}

	// ANSIカラーを有効化（Windows 10以降）
	enableANSI()

	// 複数列で表示
	for i, name := range names {
		// ファイルの情報を取得します。ディレクトリかどうかの判定に使います。
		fullPath := path + string(os.PathSeparator) + name
		info, err := os.Stat(fullPath)
		if err != nil {
			// 情報取得失敗 → 色なしで表示
			fmt.Print(name)
		} else if info.IsDir() {
			// ディレクトリは青色・太字
			fmt.Printf("\033[1;34m%-*s\033[0m", colWidth, name)
		} else {
			// 通常ファイル → デフォルト色
			fmt.Printf("%-*s", colWidth, name)
		}

		// 列の末尾か最後の要素なら改行
		// (i+1)%cols == 0 => 列の末尾に達した
		// i == len(names)-1 => 最後の要素に達した
		if (i+1)%cols == 0 || i == len(names)-1 {
			fmt.Println()
		}
	}
}

// Windowsの隠し属性チェック
// フルパスを作成
//
//	↓
//
// UTF-16に変換（Windows API用）
//
//	↓
//
// Windows APIでファイル属性を取得
//
//	↓
//
// 隠し属性のビットが立っているか確認
//
//	└── true  → 隠しファイル
//	└── false → 通常ファイル
func isHidden(dir, name string) bool {
	// ディレクトリパスとファイル名を結合してフルパスを作ります。
	// os.PathSeparator はWindowsでは \ になります。
	// "C:\Users" + "\" + "secret.txt" → "C:\Users\secret.txt"
	fullPath := dir + string(os.PathSeparator) + name
	// WindowsのAPIは文字列をUTF-16形式で受け取るため、Go の文字列（UTF-8）を変換しています。(ptr はWindows APIに渡すためのポインタです)
	ptr, err := windows.UTF16PtrFromString(fullPath)
	if err != nil {
		return false
	}
	// WindowsAPIを呼び出してファイルの属性を取得します。(属性はビットフラグの形式で返ってきます)
	// attrs                : 0b00100010
	// FILE_ATTRIBUTE_HIDDEN: 0b00000010  （値は2）
	// AND演算結果          : 0b00000010  → 0以外 → true（隠しファイル）

	// attrs                : 0b00100000
	// FILE_ATTRIBUTE_HIDDEN: 0b00000010
	// AND演算結果          : 0b00000000  → 0 → false（隠しファイルでない）
	attrs, err := windows.GetFileAttributes(ptr)
	if err != nil {
		return false
	}
	return attrs&windows.FILE_ATTRIBUTE_HIDDEN != 0
}

// WindowsターミナルのANSIエスケープコードを有効化
// LinuxやMacはデフォルトでANSIが有効ですが、Windowsは明示的に有効化しないと
// \033[1;34m のような文字がそのまま表示されてしまいます。
// この関数を呼ぶことで色付き表示が機能するようになります。
func enableANSI() {
	// Windows APIはファイルやターミナルをハンドルという識別番号で管理しています。
	// os.Stdout.Fd() で標準出力の識別番号を取得し、Windows API用の型に変換しています。
	handle := windows.Handle(os.Stdout.Fd())
	// 現在のターミナルの設定を mode に取得します。
	// &mode はmodeのメモリアドレスを渡しており、関数がそのアドレスに値を書き込む形になっています。
	// 設定はビットフラグで管理されています。
	var mode uint32
	windows.GetConsoleMode(handle, &mode)
	// ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	// 0x0004 が ENABLE_VIRTUAL_TERMINAL_PROCESSING、つまりANSIを有効にするビットです。
	// OR演算で既存の設定を壊さずにANSIのビットだけを追加しています。
	windows.SetConsoleMode(handle, mode|0x0004)
}
