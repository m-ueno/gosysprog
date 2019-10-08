# 第13章

GoのランタイムはミニOS

* OSとGoランタイムの対応（相似）
    * 物理コア：M
    * スケジューラ（ランキュー）：P
    * プロセス：G
* M:N
    * OSのスレッド数Mに対してgoroutine N

syncパッケージ

* sync.Mutexとチャネルの使いわけ [Go Wiki](https://github.com/golang/go/wiki/MutexOrChannel)
    * キャッシュ, 状態, を（競合せずに）扱いたいときはMutexで
    * まずトラディショナルなsyncパッケージで考えて、チャネルの方がシンプルになるならチャネルにする

条件変数 sync.Cond [script](./main_test.go)

* 用途
    * ブロードキャスト
    * リソースを待っているスレッドにリソース準備できたら通知
        * これはGoではチャネルでも代用できる
* 用途
    * [スレッドセーフなキューの実装](https://yohhoy.hatenablog.jp/entry/2014/09/23/193617)


sync/atomicパッケージ

* 不可算操作
    * AddUint64, AddUint64ptr, ...