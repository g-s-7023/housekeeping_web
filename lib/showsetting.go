package lib

import (
	"html/template"
	"net/http"
	"time"
)

// 設定画面に遷移する際に実行する関数
func ShowSetting(w http.ResponseWriter, r *http.Request) {
	//===
	//=== HTMLに渡すパラメータの設定
	//===
	// 設定画面を表示するためにhtmlemplateに渡す構造体
	var ParamToShowSetting struct {
		// データの一括削除に使用する年のリスト
		YearList []int
	}
	// いったんUTCで現在時刻をフォーマットして、JSTに変換
	thisYear := time.Now().UTC().In(jst).Year()
	// 年のリストに表示する年を格納
	for y := thisYear; y >= STARTYEAR; y-- {
		ParamToShowSetting.YearList = append(ParamToShowSetting.YearList, y)
	}
	//===
	//=== HTMLの出力
	//===
	// htmlTemplateを読み込み
	var html = template.Must(template.ParseFiles(DIR_HTML + "setting.html"))
	// getで渡されたindexのエントリを表示
	if err := html.ExecuteTemplate(w, "setting.html", ParamToShowSetting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
