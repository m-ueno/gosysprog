# graceful

`start_server`対応ウェブサーバ

* GETリクエスト到着後3秒後にレスポンス"Hello"を返す
* シグナルSIGTERMを受け取ると、新たな接続要求受付を止め、処理中のリクエストに対するレスポンスを返したら終了する
    * 終了時に、処理を開始したりクエスト数、処理を完了したリクエスト数を出力

## 起動方法

`start_server --port 8080 --pid-file app.pid -- ./graceful` [スクリプト](./run)

## 再起動方法

`pkill -HUP start_server`

## グレースフル・リスタートの挙動

* kill -HUPしたときの挙動
    * リクエスト中に `kill -HUP` すると、今動いているサーバは接続要求を待つのを止める. 処理中のリクエストに対するレスポンスを返したら終了する
    * この間新しいサーバが接続要求を受け入れ、処理する
* 時刻Tに再起動
    * Tまでに届いたリクエストは古いサーバから、T以降に届いたリクエストは新しいサーバから返るはず


### テストシナリオ

* T=0～10まで毎秒5リクエスト送信、T=5, T=13で再起動する [スクリプト](./client.sh)
    * 古いサーバは約25リクエスト、新しいサーバがのこりのリクエスト返すはず

実行結果

```
starting new worker 13777
received HUP (num_old_workers=TODO)
spawning a new worker (num_old_workers=TODO)
starting new worker 13896
new worker is now running, sending TERM to old workers:13777
sleep 0 secs
killing old workers
2019/10/08 11:27:12 main: shutdown start...
2019/10/08 11:27:12 svr: Accept() returned error: accept tcp 0.0.0.0:8080: use of closed network connection
2019/10/08 11:27:12 svr: channel shutdown is closed, so waiting all goroutines finish
2019/10/08 11:27:12 main: waiting server thread shutdown
2019/10/08 11:27:15 svr: bye pid=13777 accepted=28 processing=28 processed=28    # リクエスト完了を待って終了するので数が一致する
2019/10/08 11:27:15 main: shutdown done
old worker 13777 died, status:0
received HUP (num_old_workers=TODO)
spawning a new worker (num_old_workers=TODO)
starting new worker 14007
new worker is now running, sending TERM to old workers:13896
sleep 0 secs
killing old workers
2019/10/08 11:27:23 main: shutdown start...
2019/10/08 11:27:23 svr: Accept() returned error: accept tcp 0.0.0.0:8080: use of closed network connection
2019/10/08 11:27:23 svr: channel shutdown is closed, so waiting all goroutines finish
2019/10/08 11:27:23 svr: bye pid=13896 accepted=22 processing=22 processed=22    #
2019/10/08 11:27:23 main: waiting server thread shutdown
2019/10/08 11:27:23 main: shutdown done
old worker 13896 died, status:0
```