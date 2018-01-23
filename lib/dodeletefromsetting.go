package lib

import (
	"html/template"
	"net/http"
	"strconv"
)

// 設定画面でデータの一括削除をした場合に実行する関数
func DoDeleteFromSetting(w http.ResponseWriter, r *http.Request) {
	var err error
	//===
	//=== postされた値の受け取り
	//===
	// 各キーごとの値を格納する変数
	year := r.FormValue("Year")
	month := r.FormValue("Month")
	//===
	//=== 入力チェック
	//===
	// 削除対象の年と月
	var yearToDelete, monthToDelete int
	if yearToDelete, err = strconv.Atoi(year); err != nil {
		// yearが数値に変換できない場合、エラー
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if monthToDelete, err = strconv.Atoi(month); err != nil {
		// monthが数値に変換できない場合、エラー
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//===
	//=== 削除処理
	//===
	err = DeleteMonth(r, yearToDelete, monthToDelete)
	//===
	//=== ページ遷移
	//===
	// エラーの場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// htmlに渡すパラメータを作成
	src := "delete"
	// htmlファイルを読み込み
	html := template.Must(template.ParseFiles(DIR_HTML + "succeed.html"))
	if err = html.ExecuteTemplate(w, "succeed.html", src); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
