package app

import (
	"net/http"
	"housekeeping_web/housekeeping_web"
)

// todo
// スマホアプリとの連携

func init() {
	http.HandleFunc("/", housekeeping_web.Entry)
	http.HandleFunc("/logout", housekeeping_web.ShowLogout)
	http.HandleFunc("/user/list", housekeeping_web.ShowList)
	http.HandleFunc("/user/input", housekeeping_web.ShowInput)
	http.HandleFunc("/user/summary", housekeeping_web.ShowSummary)
	http.HandleFunc("/user/setting", housekeeping_web.ShowSetting)
	http.HandleFunc("/user/update", housekeeping_web.ShowUpdate)
	http.HandleFunc("/user/doinput", housekeeping_web.DoInput)
	http.HandleFunc("/user/dodeletemonth", housekeeping_web.DoDeleteFromSetting)
	http.HandleFunc("/user/doupdate", housekeeping_web.DoUpdate)
	http.HandleFunc("/user/dodelete", housekeeping_web.DoDelete)
	http.HandleFunc("/dologout", housekeeping_web.DoLogout)
	http.Handle("/css", http.FileServer(http.Dir(".")))
}
