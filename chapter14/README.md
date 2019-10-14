# 第14章 並行・並列処理の手法と設計のパターン

* 14.1 並行・並列処理の手法のパターン
* 14.2 Goにおける並行・並列処理のパターン集

## 正誤表

```
260ページ
シングルコアでも時分割でCPU時間を分け合って並列で動作します。
↓
シングルコアでも時分割でCPU時間を分け合って並行で動作します。

274ページ
誤：…グローバルVMロック（GVL）という機構があるため，並列動作はできても並行動作はせず
正：…グローバルVMロック（GVL）という機構があるため，並行動作はできても並列動作はせず

287ページ
第 3 刷 14.2.10
誤：各アクターは自律しており，並行動作するものとして考えます
正：各アクターは自律しており，並列動作するものとして考えます
```

## 14.1 並行・並列処理の手法のパターン

> 並行・並列処理の実現手法には、おおまかに区分すると、マルチプロセス、イベン
> ト駆動、マルチスレッド、ストリーミング・プロセッシングの4つのパターンがあり
> ます

マルチプロセス

* OSのプロセスによる並行・並列化
* メモリ空間が別れているので安全
* オーバーヘッドは大きい

マルチスレッド

* OSのスレッド/軽量スレッドを使った並行・並列化
* プロセスほどではないが、OSのスレッドは比較的大きなスタックメモリ（1～2メガバイト）を必要とする
* RubyやPythonは同時に1スレッドしか動かないので性能は上がらない（並行だが並列でない）

イベント駆動

* 目的：シングルスレッドあたりの性能向上
* システムプログラミングの文脈では I/O多重化のこと

ストリーミングプロセッシング

* 並行・並列処理を実現するプログラミング手法のひとつ
* 今は「ストリームプロセッサ」としてGPU、「ストリーム言語」としてCUDAが有名
    * 他にも色々あったらしい
```
これからの並列計算のためのGPGPU連載講座(VI) 様々なGPGPUプログラミング環境
https://www.cc.u-tokyo.ac.jp/public/VOL12/No6/201011_gpgpu.pdf

    3.3 ストリーミング言語

    GPGPUの発展とは独立した平行処理/並列処理のパラダイムとして、ストリーミング処理お
    よびそれを記述するストリーミング言語の研究が行われてきた。GPUによる高性能並列計
    算の可能性とプログラミングの難しさが知られるようになると、ストリーミング言語を用
    いたGPUプログラミングの研究が注目された。残念ながらCUDAが登場したことでストリー
    ミング言語を用いたGPUプログラミングへの注目度は下がってしまったが、本節ではGPUに
    対応したストリーミング言語のいくつかを簡単に紹介する。

    ストリーミング処理は、計算対象のデータをストリーム(入力ストリーム/出力ストリー
    ム)、計算内容をカーネルと定義し、プロセッサがカーネルを次々にストリームへと作用
    させていくという概念を持っている。これに対応するストリーミング言語も、言語仕様と
    してストリームとカーネルを明示的に記述する言語仕様を採用している。なお、CUDAもス
    トリーミング言語の一種と見なされることがある。

    特に2005年頃にはGPUに対応した複数のストリーミング言語が登場した。当時注目された
    ストリーミング言語の例としては、BrookGPU[12]やRapidMind[14]などが挙げられる。


NVIDIA GPUの構造とCUDAスレッディングモデル
https://www.softek.co.jp/SPG/Pgi/TIPS/public/accel/gpu-accel2.html

    GPU はそもそも、データ依存性がなく、データの再利用性も少ないアプリケーション（グ
    ラフィック処理等）で、どんどん結果を出してゆくようなストリーム型、あるいはスルー
    プット型計算を志向したものです
```

## 14.2 Goにおける並行・並列処理のパターン集

* アムダールの法則より、プログラム全体の性能向上率は、プログラム全体のうち並列化できる部分の割合Pに依存
* Pを改善するには、逐次処理を分解して、なるべく同じ実行コストの、依存関係のない処理に分ける


### 14.2.1 同期→非同期化

重い処理をばらして並列実行するためにgoroutine, channelを使う

goroutineの寿命（リーク）に注意する★補足


### 14.2.2 非同期→同期化

非同期にしたらどこかで待ち合わせる必要がある（さもないとmain関数を抜けてしまう）

* channel
    * 他の1スレッドの待ち合わせ
* select (p.278)
    * 他の複数のスレッドの待ち合わせ
    * default節なしで、イベント駆動
        * I/Oのシステムコールでいうと, epoll
    * default節ありで、非同期I/O として扱えます
* syncパッケージ (13章)
  * 使い分け https://github.com/golang/go/wiki/MutexOrChannel)


`default節ありで非同期I/Oとして扱えます`の意味がわからないがイディオムで覚えてしまえば良い。

```go
for {
	select {
	case <-ctx.Done():
		return
	default:
	}

	// なんか処理
}
```

### 14.2.3 Producer-Consumer

キューを仲立ちに非同期化。よくあるやつ。

以下バリエーション。

* 開始した順で（続く処理を）おこなう：チャネルのチャネル
    * HTTPサーバの並列化の例 -> 6章
* バックプレッシャー
    * メッセージキューのサイズを固定し、過度のキューイングをブロック
    * Goではバッファありチャネル
* 並列forループ
    * forループの内部を全てgoroutineとして並列化
* ワーカープール
    * CPUコア数を超えてワーカーを生成しても無駄な場合
        * コア数は `runtime.NumCPU()`
            * Hyper Threadingが有効だと2で割ったほうがいい
    * ワーカーの生成コストが馬鹿にならない場合?
        * goroutineだとあまりないのでは

### 他の言語の紹介

Future/Promise (14.2.8)

* Goではバッファなしチャネルの受信/送信



ReactiveX (14.2.9)

* オブザーバーパターンが少し賢くなったもの
* Observable
    * chan interface{}
* Observer
    * ハンドラ関数を持つ構造体

アクターモデル (14.2.10)

* 現在はErlang/OTPのsupervision treeがセットになったものを呼ぶことがほとんど？
    * Akkaもそうらしい
    * Supervision and Monitoring • Akka Documentation
    * https://doc.akka.io/docs/akka/2.5/general/supervision.html
    * （他のアクターモデルの実装は知らない）
* https://github.com/AsynkronIT/protoactor-go
    * Ultra fast distributed actors for Go, C# and Java/Kotlin
    * Akka.NETのGo/Java版
    * アクター：Receive()メソッドを持つ構造体
* https://github.com/teivah/gosiris
    * An actor framework for Go

## まとめ

Goの並列処理パターン集

* 非同期化 (14.2.1)
    * goroutineを作るときはリークしないように ★補足
    * Producer-Consumer (14.2.3)
        * 開始した順で（続く処理を）おこなう：チャネルのチャネル (14.2.4)
        * バックプレッシャー (14.2.5)
        * 並列forループ (14.2.6)
        * ワーカープール (14.2.7)
* 同期化 (14.2.2)
    * channel, select
        * channelよりもsyncパッケージやerrgroup (13章) を使った方が追いやすいコードが書けそう？
* 他言語の紹介
    * Future/Promise (14.2.8)
    * ReactiveX (14.2.9)
    * アクターモデル (14.2.10)

## 補足


### goroutineのリークに注意

次のように終わらないgoroutineは簡単に書けてしまう.
プロセスの寿命が長いと、そのうちOOMでランタイムごとOSに止められる.

```go
func leak1() {
    // 読み込みでブロックされる例
    ch := make(chan int)
    <-ch
}
// go leak1() でリーク

func leak2() chan bool {
    // 書き込みでブロックされる例
    ch := make(chan bool)
    go func() {
        ch <- true // ここで止まり永遠にgoroutineが終わらない
        fmt.Println("done")
    }()
    return ch
}
// _ = leak2() としてchan boolを捨ててしまうと、leak2の中のgoroutineがいつまでも終わらない
// （この場合はバッファありチャネルにすれば一応ブロックされない）
```

どうやってgoroutineの終了を保証するか。おそらくgoroutineを作ったところで終了させるのが基本。

とはいえ、goroutineは他のgoroutineからkillできないので、途中で中断するかもしれないgoroutineは、他のgoroutineから通知を受けたら処理を終えて抜けるようにするしかない。
* contextを第一引数に受け取るのが一般的
  * 例
	  * https://golang.org/pkg/net/http/#NewRequestWithContext
	  * https://golang.org/pkg/net/http/#Server.Shutdown
* 他にも、同じ構造体の別メソッドに呼応して止まる標準APIもある
  * 例 https://golang.org/pkg/net/http/#Transport.CancelRequest (deprecated)
  * (timer.Stop ?)

次のコードは、Heavy関数を途中中断できるように書き換えて、goroutineを起動したら必ずそのgoroutineに通知を送る例。

```go
func Heavy() {
	time.Sleep(time.Hour)
}

func HeavyWithContext(ctx context.Context) {
  timer := time.NewTimer(time.Hour)
  select {
  case <-ctx.Done():
    timer.Stop()
    return
  case <-timer.C:
    return
  }
}

func Foo() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  go func() {
    HeavyWithContext(ctx)
    fmt.Println("done")
  }()

  // なんかする
}

```

参考

Goroutineハンターが過労死する前に - Qiita
https://qiita.com/i_yudai/items/3336a503079ac5749c35

```
Goroutineハンターが過労死する前に - Qiita
https://qiita.com/i_yudai/items/3336a503079ac5749c35

	心得
	1.goroutine安易に作らない
	2.goroutineを作るときは同じ関数で終了させる
	3.channelよりもsync.Mutexとsync.WaitGroupで簡単に実装できないか考える
	...

	Goといえばgoroutineとchannel。しかし、channelを適切に扱い続けることは、実際のと
	ころ比較的難しい。キャパシティの考慮漏れでwriteが出来ずにgoroutineがスタックする
	ことはよくあるし、イベントループ内に巨大なselectが乱立してコードの見通しが悪くな
	ることもよくある。channelの使い方に起因したgoroutineリークは事実として非常に多
	い。
```

### 非同期APIにしないほうがいい？

Go標準APIはほとんど同期型APIで（つまりチャネルを返さない）、ユーザが必要に応じて非同期化することになっている。

ライブラリ関数も同期型にしたほうがいいのか。

```
アプリ：同期＋非同期
---同期型API---
3rd partyライブラリ
  内部は非同期
---同期型API---
標準ライブラリ
  内部は非同期
```


* チャネルを返す代わりに、コールバック関数を受け取るパターン
`os.Walk(walkFn)`
* 時間が来たらチャネルに書き込む代わりに、一定時間ブロックするパターン

