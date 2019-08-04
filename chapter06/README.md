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

プロセス間通信の一種

IPとポートで外部とつながる

ソケットの種類

* TCP
* UDP
* Unixドメインソケット

## 6.4 ソケット通信の基本構造

* サーバー
    * Listen()してAccept()
* クライアント
    * Dial()

net.Connインタフェース

最低限のTCPサーバ

for { Accept(); go func }

    しれっと書いているがC10Kはおきない
    スレッド
    コア

## 6.5 Go言語でHTTPサーバーを実装する

> ここでは Go 言語に組み込まれている TCP の機能 (net.Conn) だけを使って HTTP による通信を実現してみましょう。

6.5〜6.9までHTTPサーバとHTTPクライアント1つずつつくる。

だんだん速度改善していく。

### 6.5.1 TCP ソケットを使った HTTP サーバー

`httpserver00/server/main.go`

```
Server is running at localhost:8888
Accept 127.0.0.1:63694
GET / HTTP/1.1
Host: localhost:8888
Accept: */*
User-Agent: curl/7.58.0
```

### 6.5.2 TCP ソケットを使った HTTP クライアント

## 6.6 速度改善( 1 ) : HTTP/1.1 の Keep-Alive に対応させる

### 6.6.1 Keep-Alive 対応の HTTP サーバー

## 6.7 速度改善( 2 ) : 圧縮

## 6.8 速度改善( 3 ) : チャンク形式のボディー送信

## 6.9 速度改善( 4 ) : パイプライニング