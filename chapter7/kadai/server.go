package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apbgo/go-study-group/chapter7/kadai/model"
)

func userFortuneHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// リクエストBodyの内容を取得
	var req model.Request

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 返ってきたレスポンスの内容を表示
	fmt.Println(req)

	// レスポンスの作成
	response := model.Response{
		Status: http.StatusOK,
		Data:   fmt.Sprintf("ID:%vの%sさんの運勢は%sです！", req.UserID, req.Name, doFortune()),
	}

	var res bytes.Buffer
	enc := json.NewEncoder(&res)
	if err := enc.Encode(response); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(res.Bytes())
}

// 処理ハンドラ
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, server.")
}

// 処理ハンドラ
func fortuneHandler(w http.ResponseWriter, r *http.Request) {

	p := r.FormValue("p")
	if p == "cheat" {
		fmt.Fprint(w, "大吉")
		return
	}
	fmt.Fprint(w, doFortune())
}

func doFortune() (fortune string) {

	switch rand.Intn(4) {
	case 0:
		fortune = "大吉"
	case 1:
		fortune = "中吉"
	case 2:
		fortune = "吉"
	case 3:
		fortune = "凶"
	}
	return fortune
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/fortune", fortuneHandler)
	mux.HandleFunc("/user_fortune", userFortuneHandler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// OSからのシグナルを待つ
	go func() {
		// SIGTERM: コンテナが終了する時に送信されるシグナル
		// SIGINT: Ctrl+c
		sigCh := make(chan os.Signal, 1)
		// 受け取るシグナルを指定
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		// チャネルでの待受、シグナルを受け取るまで以降は処理されない
		<-sigCh

		log.Println("start graceful shutdown server.")
		// タイムアウトのコンテキストを設定
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// Graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
			// 接続されたままのコネクションも明示的に切る
			srv.Close()
		}
		log.Println("HTTPServer shutdown.")
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Print(err)
	}
}
