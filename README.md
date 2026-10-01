# command-line-tool-for-windows

Linuxでおなじみのコマンドを、Windowsのコマンドプロンプトで使えるようにGoで実装したものです。
`install.bat` を実行すると、すべてのコマンドをまとめてインストールできます。

## コマンド一覧

| コマンド | 機能 | 使い方 |
|---|---|---|
| `cat` | ファイルの内容を表示する | `cat <ファイル名>` |
| `ls` | ファイルとフォルダの一覧を表示する（フォルダは青色、隠しファイルは除外） | `ls [パス]` |
| `pwd` | カレントディレクトリを表示する | `pwd` |
| `rm` | ファイル・フォルダを削除する | `rm <パス> [パス2 ...]` |
| `stat` | ファイルのサイズ・日時・属性を表示する | `stat <ファイル名>` |
| `touch` | 空のファイルを作成する（既にある場合は更新日時を現在時刻にする） | `touch <ファイル名>` |

> **注意:** `rm` はフォルダを指定すると、確認なしで中身ごと削除します。

## 前提条件

- Go 1.26.4以上がインストールされていること
- Goのインストール先フォルダ（通常は `C:\Users\<ユーザー名>\go\bin`）がPATHに含まれていること

インストール先フォルダは次のコマンドで確認できます。`GOBIN` が空の場合は `GOPATH` の下の `bin` フォルダです。

```cmd
go env GOBIN GOPATH
```

## インストール方法

### 1. リポジトリを取得

```cmd
git clone https://github.com/ue555/command-line-tool-for-windows.git
cd command-line-tool-for-windows
```

### 2. install.bat を実行

```cmd
install.bat
```

エクスプローラーから `install.bat` をダブルクリックしても実行できます。

**期待される出力：**
```
[OK] cat
[OK] ls
[OK] pwd
[OK] rm
[OK] stat
[OK] touch

Install result: 6 succeeded, 0 failed
Install folder: C:\Users\<ユーザー名>\go\bin
```

### 3. インストール確認

```cmd
where cat
```

**期待される出力：**
```
C:\Users\<ユーザー名>\go\bin\cat.exe
```

複数行表示された場合は、1行目のものが実行されます。

## プログラムの更新

ソースコードを修正した後、もう一度 `install.bat` を実行します。すべてのコマンドが最新版に更新されます。

```cmd
install.bat
```

## アンインストール

```cmd
install.bat uninstall
```

インストール先フォルダから、このリポジトリのコマンドの実行ファイルだけを削除します。

## install.bat の表示の意味

| 表示 | 意味 |
|---|---|
| `[OK] <コマンド名>` | インストール（または削除）に成功 |
| `[NG] <コマンド名>` | 失敗。直前に表示されるGoのエラーメッセージを確認してください |
| `[--] <コマンド名> : not installed` | アンインストール時、もともとインストールされていなかった |
| `[WARN] another "<コマンド名>" comes first on PATH` | 同名の別コマンドがPATHの手前にあり、そちらが優先される |
| `[WARN] The install folder is not on PATH` | インストール先フォルダがPATHに含まれていない |

## コマンドを追加する場合

リポジトリのルートに、コマンド名のフォルダを作って `go.mod` と `main.go` を置きます。
`install.bat` は `go.mod` があるフォルダを自動で対象にするため、`install.bat` の修正は不要です。

```cmd
mkdir head
cd head
go mod init github.com/ue555/head
```

## 開発者向け

以下は `cat` を例にしています。他のコマンドも同じ手順です。

### ビルドせずに実行する

```cmd
cd cat
go run . go.mod
```

### 1つのコマンドだけインストールする

```cmd
cd cat
go install .
```

## トラブルシューティング

### `install` と入力すると「missing file operand」と表示される

MSYS2などに含まれる別の `install.exe` が実行されています。拡張子を付けて `install.bat` と入力してください。

### 「Go was not found」と表示される

Goがインストールされていないか、PATHに含まれていません。<https://go.dev/dl/> からインストールしてください。

### 「The install folder is not on PATH」と表示される

インストール自体は成功していますが、コマンド名だけでは実行できない状態です。
`Install folder:` に表示されたフォルダをユーザー環境変数 `Path` に追加し、コマンドプロンプトを開き直してください。

### 「another "..." comes first on PATH」と表示される

MSYS2やGit for Windowsなど、同名のコマンドを含むフォルダがPATHの手前にあります。
このリポジトリのコマンドを使いたい場合は、環境変数 `Path` でインストール先フォルダをそれらより上に移動してください。

### PowerShellでの注意点

PowerShellでは `cat` `ls` `pwd` `rm` が組み込みのエイリアスとして定義されているため、そちらが優先されます。
このリポジトリのコマンドを使う場合は、`.exe` を付けて実行してください。

```powershell
cat.exe go.mod
```

### install.bat を編集する場合

`install.bat` には日本語を書かないでください。コマンドプロンプトがファイルを正しく読めなくなり、エラーになります。

## ライセンス

MIT License

## 作者

@ue555
