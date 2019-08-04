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

### 6.5.1

```
Server is running at localhost:8888
Accept 127.0.0.1:63694
GET / HTTP/1.1
Host: localhost:8888
Accept: */*
User-Agent: curl/7.58.0
```

