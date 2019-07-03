# Chapter03: 低レベルアクセスへの入口2：io.Reader

```go
// Implementations must not retain p.
type Reader interface {
	Read(p []byte) (n int, err error)
}
```
## 3.2 補助関数

* ioutil.ReadAll(r io.Reader) ([]byte, err)
* io.Copy(w, r)
* io.CopyN(w, r, n) nバイトコピー
* io.CopyBuffer(w, r, buf) バッファを自分で用意しコピー


## 3.3 複合インタフェース

```go
// ReadWriteCloser is the interface that groups the basic Read, Write and Close methods.
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}
```

## 3.4

> `bytes.Buffer` はio.Writerとしてもio.Readerとしても使えます。
> 読み出しに使えるものとしてはこれだけ覚えておけば（ほぼ）問題ないでしょう。

## 3.5 バイナリ解析

## 3.6 テキスト解析

* bufio.ReaderのReadString(), ReadBytes()
* bufio.Scanner
* fmt.Fscan