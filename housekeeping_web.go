package develop

import (
	"net/http"
	"housekeeping_web/lib"
)

// (開発サーバ)
// app.yamlとこのファイルをdevelopフォルダにおく
// DIR_HTMLの値を"../html/"、cssのstatic_dirを"../css"

// (デプロイ)
// app.yamlとこのファイルをトップレベルに置く
// デプロイ DIR_HTML1の値を"html/"、cssのstatic_dirを"css"

func init() {
	http.HandleFunc("/", lib.Entry)
	http.HandleFunc("/logout", lib.ShowLogout)
	http.HandleFunc("/user/list", lib.ShowList)
	http.HandleFunc("/user/input", lib.ShowInput)
	http.HandleFunc("/user/summary", lib.ShowSummary)
	http.HandleFunc("/user/setting", lib.ShowSetting)
	http.HandleFunc("/user/update", lib.ShowUpdate)
	http.HandleFunc("/user/doinput", lib.DoInput)
	http.HandleFunc("/user/dodeletemonth", lib.DoDeleteFromSetting)
	http.HandleFunc("/user/doupdate", lib.DoUpdate)
	http.HandleFunc("/user/dodelete", lib.DoDelete)
	http.HandleFunc("/dologout", lib.DoLogout)
	http.Handle("/css", http.FileServer(http.Dir(".")))
}
