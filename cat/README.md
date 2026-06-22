# cat - シンプルなファイル表示ツール

Goで実装されたシンプルな`cat`コマンドです。指定されたファイルの内容を標準出力に表示します。

## 機能

- ファイルの内容を表示
- Unix/Linuxの`cat`コマンドと同様の基本機能

## 前提条件

- Go 1.26.4以上がインストールされていること
- `GOPATH\bin`（通常は`C:\Users\<ユーザー名>\go\bin`）がPATHに含まれていること

## インストール方法

### 1. プロジェクトディレクトリに移動

```cmd
cd C:\Users\zeroz\dev\command-line-tool\cat
```

### 2. バイナリをビルド・インストール

```cmd
go install .
```

これで`C:\Users\zeroz\go\bin\cat.exe`が作成されます。

### 3. インストール確認

```cmd
where cat.exe
```

**期待される出力：**
```
C:\Users\zeroz\go\bin\cat.exe
```

## 使い方

```cmd
cat <ファイル名>
```

### 使用例

```cmd
# READMEファイルを表示
cat README.md

# Goのソースコードを表示
cat main.go

# go.modを表示
cat go.mod

# 絶対パスで指定
cat C:\temp\test.txt
```

## 動作確認

プロジェクトディレクトリで以下を実行：

```cmd
cat go.mod
```

**期待される出力：**
```
module github.com/ue555/cat

go 1.26.4
```

## プログラムの更新

main.goを修正した後、以下を実行：

```cmd
cd C:\Users\zeroz\dev\command-line-tool\cat
go install .
```

これで`cat.exe`が最新版に更新されます。

## アンインストール

```cmd
del C:\Users\zeroz\go\bin\cat.exe
```

## 開発者向け

### ローカルでテスト実行

ビルドせずに直接実行する場合：

```cmd
go run main.go <ファイル名>
```

**例：**
```cmd
go run main.go go.mod
```

### 手動ビルド

特定の場所にバイナリを作成する場合：

```cmd
go build -o cat.exe .
```

## トラブルシューティング

### 「コマンドが見つかりません」エラー

**確認1: バイナリが存在するか**
```cmd
dir C:\Users\zeroz\go\bin\cat.exe
```

**確認2: PATHに含まれているか**
```cmd
echo %PATH% | findstr "go\bin"
```

**確認3: 現在のセッションでPATHを更新**
```cmd
set PATH=%PATH%;C:\Users\zeroz\go\bin
```

その後、新しいコマンドプロンプトを開いて再度試してください。

### PowerShellでの注意点

PowerShellには組み込みの`cat`エイリアス（`Get-Content`のエイリアス）があります。
自作のcat.exeを使用する場合は、明示的に`.exe`を付けて実行してください：

```powershell
cat.exe go.mod
```

または相対パスで指定：

```powershell
.\cat.exe go.mod
```

## よくある間違い

### ❌ ファイルを二重指定
```cmd
go run main.go main.go  # エラー: case-insensitive file name collision
```

### ✅ 正しい実行方法
```cmd
go run main.go <読み込むファイル>
```

## ライセンス

MIT License

## 作者

@ue555
