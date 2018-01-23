package lib

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func DoUpdate(w http.ResponseWriter, r *http.Request) {
	var err error
	//===
	//=== postされた値の受け取り
	//===
	// r.Formに各値を格納
	r.ParseForm()
	// 各キーごとの値を格納する配列
	idToUpdate := r.FormValue("ID")
	dateBefore := r.FormValue("DateBefore")
	dateToUpdate := r.FormValue("Date")
	categoryToUpdate := r.FormValue("Category")
	detailToUpdate := r.FormValue("Detail")
	priceToUpdate := r.FormValue("Price")
	//===
	//=== 入力チェック
	//===
	// 更新に必要なパラメータ
	paramToUpdate := new(ParamToUpdate)
	// エラーチェック用の関数(エラー時falseを返却)
	errorCheck := func() bool {
		if dateBefore == "" || dateToUpdate == "" || categoryToUpdate == "0" || priceToUpdate == "" {
			// 必須入力項目のいずれかが空白ならエラー
			return false
		}
		// 日付の入力値チェック
		if date := strings.Split(dateToUpdate, "-"); len(date) != 3 {
			// 年月日の3要素以外の場合、エラー
			return false
		} else {
			if _, err = time.Parse("2006-01-02", dateToUpdate); err != nil {
				// 日付が現実にあるものでなければエラー
				return false
			}
			paramToUpdate.Year, _ = strconv.Atoi(date[0])
			paramToUpdate.Month, _ = strconv.Atoi(date[1])
			paramToUpdate.Day, _ = strconv.Atoi(date[2])
		}
		// 更新前の日付の入力値チェック
		if date := strings.Split(dateBefore, "-"); len(date) != 3 {
			// 年月日の3要素以外の場合、エラー
			return false
		} else {
			if _, err = time.Parse("2006-01-02", dateBefore); err != nil {
				// 日付が現実にあるものでなければエラー
				return false
			}
			// 各値をintに変換
			paramToUpdate.YearBefore, _ = strconv.Atoi(date[0])
			paramToUpdate.MonthBefore, _ = strconv.Atoi(date[1])
		}
		// 費目のチェック
		if category, err := strconv.Atoi(categoryToUpdate); err != nil {
			// 費目が数値でなければエラー
			return false
		} else {
			if _, ok := Categories[category]; ok {
				// 費目が定義されていれば、DB登録用のパラメータにコピー
				paramToUpdate.Category = category
			} else {
				// 費目が定義されていなければエラー
				return false
			}
		}
		// 価格のチェック
		if price, err := strconv.Atoi(priceToUpdate); err != nil {
			// 価格が数値でなければエラー
			return false
		} else {
			if price < 0 {
				// 価格が0以上でなければエラー
				return false
			} else {
				// 価格が0以上なら、DB登録用のパラメータにコピー
				paramToUpdate.Price = price
			}
		}
		return true
	}
	if !errorCheck() {
		// エラーチェックに引っかかった場合、元のページに戻る
		// htmlTemplateに渡すパラメータ
		cIndex, _ := strconv.Atoi(categoryToUpdate)
		paramToReturn := ParamToShowUpdate{
			CategoryList:  make(map[int]string),
			HasAnyError:   true,
			DatastoreId:   idToUpdate,
			Date:          dateToUpdate,
			Category:      categoryToUpdate,
			CategoryIndex: cIndex,
			Detail:        detailToUpdate,
			Price:         priceToUpdate,
		}
		// 選択肢用の費目
		for key, val := range Categories {
			paramToReturn.CategoryList[key] = val.Name
		}
		// htmlファイルを読み込み
		html := template.Must(template.ParseFiles(DIR_HTML + "update.html"))
		if err := html.ExecuteTemplate(w, "update.html", paramToReturn); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// エラーが出なかった場合、残りのパラメータをセット
	paramToUpdate.ID, _ = strconv.Atoi(idToUpdate)
	paramToUpdate.Detail = detailToUpdate
	//===
	//=== DBのデータ更新
	//===
	err = Update(r, paramToUpdate)
	//===
	//=== ページ遷移
	//===
	// データ更新エラーの場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// htmlに渡すパラメータを作成
	src := "Update"
	// htmlファイルを読み込み
	var html = template.Must(template.ParseFiles(DIR_HTML + "succeed.html"))
	if err = html.ExecuteTemplate(w, "succeed.html", src); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
