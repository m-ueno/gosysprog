# 2. 低レベルアクセスの入口1：io.Writer

* `io.Writer`：出力の抽象化
    * Cではファイル読み書きに、OSのファイルディスクリプタを使うが、
    * Goではio.Writerインタフェースを使う
* NOTE
    * Goでも直接ファイルディスクリプタ指定できる：`file, err = os.NewFile(ファイルディスクリプタ, 名前)`

## 2.3 io.Writerはインタフェース

<https://golang.org/src/io/io.go#L89>

```go
// [io/io.go]
// Implementations must not retain p.
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

## 2.4 io.Writerを使う構造体の例

* 2.4.1 ファイル出力 os.File
* 2.4.2 画面出力 os.Stdout
* 2.4.3 バッファ(1) bytes.Buffer
* 2.4.4 バッファ(2) strings.Builder
  * Go 1.10から追加された。書き出し専用。bytes.Bufferの代わりに使える。
* 2.4.5 net.TCPConn
* 2.4.6 io.Writerのデコレータ
    * `func MultiWriter(writers ...Writer) Writer`
    * バッファ付き出力 `bufio.Writer`

* 2.4.7 フォーマットしてデータをio.Writerに書き出す fmt.Fprintf

## 2.5 インタフェースの実装状況・利用状況を調べる

`godoc -http ":6060" -analysis type` と `analysis type`オプションをつけるとインタフェースの分析を行ってくれます。

## 2.7 Tips

データの入出力や加工を行う関数を書く場合、
ファイルを読み書きするコードを書く場合でも、`io.Writer / io.Reader`を受け取る関数を書くのが望ましい

## 問題

### Q2.3

* JSONをgzip化してクライアントに返す
* gzipする前のJSONをstdoutに吐き出す
* io.MultiWriter
* Flush()を忘れない

* map -> json -> stdout
              -> gzip -> w


## おまけ：バッファ付き出力の速度比較

```
% go test -bench .
goos: linux
goarch: amd64
pkg: github.com/m-ueno/gosysprog/chapter2
BenchmarkWriteFile-4                  10         182093279 ns/op
BenchmarkWriteFileBuffered-4        1000           2167811 ns/op
BenchmarkReadFile-4                   10         132938858 ns/op
BenchmarkReadFileBuffered-4          200           7933353 ns/op
```

## おまけ：ポインタか値か

* https://stackoverflow.com/questions/23542989/pointers-vs-values-in-parameters-and-return-values

* レシーバーは、副作用があるならポインタ
* スライス・マップ・インタフェース・stringsは、内部ではポインタなのであえてポインタにしなくていい
* その他は、でかいオブジェクトを書き換えるならポインタ。そうでなければ値
