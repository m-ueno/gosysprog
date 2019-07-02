# Chapter03: 低レベルアクセスへの入口2：io.Reader

```go
// Implementations must not retain p.
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

## 3.3 複合インタフェース

```go
// ReadWriteCloser is the interface that groups the basic Read, Write and Close methods.
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}
```
