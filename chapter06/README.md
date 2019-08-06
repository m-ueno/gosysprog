# TCPソケットとHTTPの実装

ここからはGo言語の低レベルなAPIを紹介.

この章はソケット通信とHTTP.

## 6.1 プロトコルとレイヤー

## 6.2 HTTPとその上のプロトコルたち

* ネットワークの低レベルな部分をGo言語で見る前に、HTTPの基本的な仕組みと、更に上位のレイヤーをいくつかピックアップして紹介
    * HTTP
    * RPC
        * GoではJSON-RPC (1.0) ライブラリが標準添付
        * JSON-RPC 2.0 はサードパーティのがたくさん <https://godoc.org/?q=json-rpc+2.0>
    * REST
    * GraphQL
* p.103 コラム プロトコルにおけるシステムコールともいえるRFC

## 6.3 ソケットとは

アプリケーション層からトランスポート層のプロトコルを利用するときのAPIとして、
ほとんどのOSで提供されている仕組み

* ソケットは、数多くあるプロセス間通信の一種、でもある
    * ソケットが他と違う点は、はIPとポートで外部とつながる点

ソケットの種類 本書では3つ

* TCP
* UDP
* Unixドメインソケット

## 6.4 ソケット通信の基本構造

* サーバー
    * Listen()してAccept()
* クライアント
    * Dial()

net.Connインタフェース：通信のための共通機能が実装されたインタフェース

p.107
最低限のTCPサーバ

```go
ln, _ := net.Listen("tcp", ":8080")
for {
    conn, _ :=  ln.Accept()
    go func() {
        //...
    }
}
```

しれっと書いているがC10K問題はおきない

## 6.5 Go言語でHTTPサーバーを実装する

> ここでは Go 言語に組み込まれている TCP の機能 (net.Conn) だけを使って HTTP による通信を実現してみましょう。

6.5〜6.9まで.

HTTPサーバとHTTPクライアントそれぞれ1つずつつくる。

だんだん速度改善していく。

### 6.5.1 TCP ソケットを使った HTTP サーバー

基本構造は6.4章のソケット通信そのままです。（「最低限のTCPサーバ」）

`httpserver00/server/main.go`

```
Server is running at localhost:8888
Accept 127.0.0.1:63694
GET / HTTP/1.1
Host: localhost:8888
Accept: */*
User-Agent: curl/7.58.0
```

※goroutineの最後にconn.Close()がある。毎回Accept()してClose()する

### 6.5.2 TCP ソケットを使った HTTP クライアント

## 6.6 速度改善( 1 ) : HTTP/1.1 の Keep-Alive に対応させる

### 6.6.1 Keep-Alive 対応の HTTP サーバー

> このコードで重要なのは、 Accept() を受信したあとに for ループがある点です。
これにより、 TCP のコネクションが張られたあとに何度もリクエストを受けられるよ
うにしています。

※forループのたびにReadRequest()を呼ぶ

## 6.7 速度改善( 2 ) : 圧縮

> リクエスト生成部を改造して、自分が対応しているアルゴリズムを宣言する
ようにします。サーバーから自分が理解できない圧縮フォーマットでデータを送りつ
けられても、クライアントではそれを読み込めないからです。下記のように、リクエ
ストヘッダーの "Accept-Encoding" に「このクライアントは gzip 圧縮を処理でき
ます」という表明を入れます。

```go
request, err := http.NewRequest(
    "POST",
    "http://localhost:8888",
    strings.NewReader(sendMessages[current]))
if err != nil {
    panic(err)
}
request.Header.Set("Accept-Encoding", "gzip")
```

## 6.8 速度改善( 3 ) : チャンク形式のボディー送信

## 6.9 速度改善( 4 ) : パイプライニング

> この機能はパイプライニングと呼ばれ、 HTTP/1.1 の規格にも含まれています。パ
イプライニングでは、レスポンスがくる前に次から次にリクエストを多重で飛ばすこ
とで、最終的に通信が完了するまでの時間を短くします(図 6.9 )

* 後方互換性は崩れる (HTTP/1.0しか解釈しないプロキシが経路上にあると使えない)

### 6.9.1 パイプライニングのサーバー実装

* 仕様
    * サーバー側の状態を変更しないメソッド (GET, HEAD) であれば、サーバー側で並列処理を行う
        * ※これまでと同様
    * リクエストの順序でレスポンスを返さなければならない
        * ※キューとしてバッファ付きチャネルを使う

```go
// 順番に従って conn に書き出しをする (goroutine で実行される )
func writeToConn(sessionResponses chan chan *http.Response, conn net.Conn) {
    defer conn.Close()
    // 順番に取り出す
    for sessionResponse := range sessionResponses {
        // 選択された仕事が終わるまで待つ
        response := <-sessionResponse
        response.Write(conn)
        close(sessionResponse)
    }
}
```

挙動の説明



### 6.9.3 パイプライニングと HTTP/2

## まとめ

表 6.2
本章で紹介してきた高速化手法

|手法|効果|
|--|--|
|Keep-Alive|再接続のコストを削減|
|圧縮|通信時間の短縮|
|チャンク|レスポンスの送信開始を早める|
|パイプライニング|通信の多重化|