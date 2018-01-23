package app

import (
	"net/http"
	"time"
	"strings"
	"strconv"
	"html/template"
)

// 一覧画面の表示
func ShowList(w http.ResponseWriter, r *http.Request) {

	// なぜかhtmlが表示されず、HTTP 500が返ってくる
	// ShowListまでリダイレクトはされる
	// ShowListでfprintfで単純な文字列を表示させるだけならできる
	// yearToShow, monthToShowの表示までもできる
	// ReadListもerrorは返ってきてない
	// どうやらhtmlとcssが読み込めていないらしい(このファイルに直書きしたらOKだった)
	// app.yamlをトップレベルに持ってきたらhtmlとcssが読み込めた
	// 他のファイルもhtmlを同じようにする

	var err error
	//===
	//=== postで渡された値の格納
	//===
	requestedMonth := r.FormValue("selected_month")
	//===
	//=== 入力チェックと表示する年月の設定
	//===
	// 表示する月を格納する変数
	var yearToShow, monthToShow int
	if requestedMonth == "" {
		// requestedMonthがセットされていない場合、現在の月をセット
		// いったんUTCで現在時刻をフォーマットして、JSTに変換
		now := time.Now().UTC().In(jst)
		// 年・月の取り出し
		yearToShow = now.Year()
		monthToShow = int(now.Month())
	} else {
		// requestMonthがセットされている場合、その値を設定
		if splitString := strings.Split(requestedMonth, "-"); len(splitString) != 2 {
			// 年・月の2要素のみでなかった場合、エラー
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			// 年・月の設定
			yearToShow, _ = strconv.Atoi(splitString[0])
			monthToShow, _ = strconv.Atoi(splitString[1])
		}
	}
	//===
	//=== 表示するデータの取得
	//===
	// contextの作成
	paramToShowList, err := ReadList(r, yearToShow, monthToShow)
//	_, err = ReadList(r, yearToShow, monthToShow)
	//===
	//=== ページ遷移
	//===
	// エラーの場合
	if err != nil {
		// test
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// htmlファイルを読み込み
	html := template.Must(template.ParseFiles("html/list.html"))
	// htmlの出力
	err = html.ExecuteTemplate(w, "list.html", paramToShowList)
	if  err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
