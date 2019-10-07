# 第14章 並行・並列処理の手法と設計のパターン

> 本章では、並行・並列処理の手法に関する一般的なパターンと、Go言語における並行・並列処理の設計のパターンを取り上げます。

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

イベント駆動

* 目的：シングルスレッドあたりの性能向上
* システムプログラミングの文脈では I/O多重化のこと

マルチプロセス、マルチスレッド

* OSのプロセス/スレッドを使った並列化

ストリーミングプロセッシング

* 並行・並列処理を実現するプログラミング手法のひとつ
* 現代では、ストリームプロセッサとしてGPU, ストリーム言語としてCUDAが有名
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
* Pを改善するには、逐次処理をなるべく分解して、同じ粒度のシンプルなたくさんのジョブに分ける

### 14.2.1 同期→非同期化

重い処理をばらして並列実行するためにgoroutine, channelを使う

#### 重い関数をgoroutineでラップするパターン

* goroutineの寿命（リーク）に注意

```go

// 重い関数, 同期
func Heavy() int {...}

// 重い関数をラップ、チャネルを返す関数
func HeavyAsync() (chan int) {
    chanResult := make(chan int)

    go func() {
        chanResult <- Heavy()
    }

    return chanResult
}

func main() {
    // HeavyAsyncの返り値をreadしていれば, 完了を待ってgoroutineを抜けるが
    n := <- HeavyAsync()
    fmt.Println(n)
}

func main() {
    // HeavyAsyncの返り値をreadしていなければ、goroutineがいつまでも終わらない
    // 途中でキャンセルもされない
    select {
    case n := <-HeavyAsync(10):
        break
    case n := <-HeavyAsync(11):
        break
    }
}
```

#### 重い関数 (同期) を、タイムアウト付きにする

タイムアウトが先に来た場合、Heavyは裏で実行したままになる。
（→完了次第、resultに書き込んで終わる）

```go
// 同期
func HeavyWithTimeout(timeout time.Duration) (int, error) {
    result := make(chan int, 1)

    go func() {
        result <- Heavy()
    }()

    select {
    case k := <-result:
        return k, nil
    case <-time.After(timeout):
        return 0, errors.New("timed out")
    }
}
```

#### goroutineのリークに注意

```
Goroutineハンターが過労死する前に - Qiita
https://qiita.com/i_yudai/items/3336a503079ac5749c35
> 心得
> goroutine安易に作らない
> goroutineを作るときは同じ関数で終了させる
> channelよりもsync.Mutexとsync.WaitGroupで簡単に実装できないか考える
> ...
```

次のように終わらないgoroutineは簡単に書けてしまう.
プロセスの寿命が長いと、そのうちOOMでプロセスごとOSに殺される.

```go
func leak() {
    // 読み込みでブロックされる例
    ch := make(chan int)
    <-ch
}
// go leak() でリーク

func leak2() chan bool {
    // 書き込みでブロックされる例
    // （この場合はバッファありチャネルにすれば一応ブロックされない）
    ch := make(chan bool)
    go func() {
        ch <- true // ここで止まり永遠にgoroutineが終わらない
        fmt.Println("done")
    }()
    return ch
}        err := s0.Activate(cctx)
        if err != nil {
            log.Printf("s0 stopped: %s", err)
        }

// _ = leak2() としてchan boolを捨ててしまうと、leak2の中のgoroutineがいつまでも終わらない
```

### 14.2.2 非同期→同期化

非同期にしたらどこかで待ち合わせる必要がある。さもないとmain()を抜けてしまう

* channel
    * 他の単一スレッドの、待ち合わせ
* select (p.278)
    * 他の複数のスレッドの、待ち合わせ
    * default節なしで、イベント駆動
        * I/Oのシステムコールでいうと, epoll
        * vs. フロー駆動
        * イベントをキューイング, 対応するサブルーチンを実行
    * default節ありで、非同期I/O として扱えます
* syncパッケージ

### Producer-Consumer

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

* Future/Promise
    * Goではバッファなしチャネルの受信/送信
* ReactiveX
    * 高級なオブザーバーパターン
    * Observable
        * chan interface{}
    * Observer
        * ハンドラ関数を持つ構造体
* アクターモデル
    * 現在はErlang/OTPのsupervision treeがセットになったものを呼ぶことがほとんど？
        * Akkaもそうらしい
        * Supervision and Monitoring • Akka Documentation
        * https://doc.akka.io/docs/akka/2.5/general/supervision.html
        * 他のアクターモデルの実装は知らない
    * https://github.com/AsynkronIT/protoactor-go
        * Ultra fast distributed actors for Go, C# and Java/Kotlin
        * Akka.NETのGo/Java
        * アクター
            * Receive()メソッドを持つ構造体
    * https://github.com/teivah/gosiris
        * An actor framework for Go

## まとめ

Goの並列処理パターン集

* 非同期化 (14.2.1)
    * goroutineを作るときはリークしないように
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

<!--
## channelよりもsyncパッケージ (13章) を使った方が追いやすいコードが書けそう？

```
Goroutineハンターが過労死する前に - Qiita
https://qiita.com/i_yudai/items/3336a503079ac5749c35

    Goといえばgoroutineとchannel。しかし、channelを適切に扱い続けることは、実際のと
    ころ比較的難しい。キャパシティの考慮漏れでwriteが出来ずにgoroutineがスタックする
    ことはよくあるし、イベントループ内に巨大なselectが乱立してコードの見通しが悪くな
    ることもよくある。channelの使い方に起因したgoroutineリークは事実として非常に多
    い。
```

上のように、channelのキャパシティ考慮漏れ、巨大selectの見通しの悪さ、goroutineのリーク
を気をつけて並列化しないとひどいことになる。

対策「Activateパターン」

* type Statefull
* method
    * Activate(context.Context) error
        * 無限ループ. for { select }
        * contextでキャンセルできる
    * Put(context.Context, value []byte) error
        * leaseIDを元に値を書き込む
    * Current() int
        * leaseIDを返す
    * withTimeout(base context.Context) context.Context

使う側

```go
ctx, cancel := ..
defer cancel()
```
 -->