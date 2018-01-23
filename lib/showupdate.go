package lib

import (
	"html/template"
	"net/http"
	"strconv"
)

// 更新画面の表示
func ShowUpdate(w http.ResponseWriter, r *http.Request) {
	//===
	//=== postで渡された値の格納
	//===
	id := r.FormValue("ID")
	year := r.FormValue("Year")
	month := r.FormValue("Month")
	day := r.FormValue("Day")
	category := r.FormValue("Category")
	categoryIndex, _ := strconv.Atoi(r.FormValue("CategoryIndex"))
	detail := r.FormValue("Detail")
	price := r.FormValue("Price")
	//===
	//=== 表示用の構造体の作成
	//===
	// htmlTemplateに渡すパラメータ
	paramToUpdate := ParamToShowUpdate{
		CategoryList:  make(map[int]string),
		HasAnyError:   false,
		DatastoreId:   id,
		Date:          year + "-" + month + "-" + day,
		Category:      category,
		CategoryIndex: categoryIndex,
		Detail:        detail,
		Price:         price,
	}
	// 選択肢用の費目
	for key, val := range Categories {
		paramToUpdate.CategoryList[key] = val.Name
	}
	//===
	//=== 画面表示
	//===
	// htmlTemplateを読み込み
	var html = template.Must(template.ParseFiles(DIR_HTML + "update.html"))
	// getで渡されたindexのエントリを表示
	if err := html.ExecuteTemplate(w, "update.html", paramToUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
