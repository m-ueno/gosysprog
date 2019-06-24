# chapter2

* `io.Writer`：出力の抽象化
    * ファイル
    * ファイルディスクリプタ
        * プロセスごと
        * 0,1,2 非負整数
    * Goでは、OSのファイルディスクリプタを直接使うのではなく、io.Writerインタフェースをかぶせている
* Note
    * Goで直接fdいじれる：`file, err = os.NewFile(ファイルディスクリプタ, 名前)`

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

* 2.4.1 ファイル出力
* 2.4.2 画面出力
* 2.4.3 bytes.Buffer
* 2.4.4 strings.Builder
* 2.4.5 net.TCPConn
* 2.4.6 io.Writerのデコレータ
    * `io.MultiWriter`
    * バッファ付き出力 `bufio.Writer`
        * おまけ：バッファ無しとの速度比較

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

* 2.4.7 フォーマットしてデータをio.Writerに書き出す fmt.Fprintf

## 問題
