package lib

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
	"time"
)

// 入力されたデータをDBに登録する関数
func InsertKakeibo(r *http.Request, p []ParamToInsert) error {
	var err error
	//===
	//=== DB登録用のパラメータの作成
	//===
	// 登録するエントリ数
	putNum := len(p)
	// DBに登録する家計簿のエントリ
	var entryToInsert = make([]Kakeibo, putNum)
	// 登録用エントリへの値のコピー
	for i := 0; i < putNum; i++ {
		// 曜日の算出
		dayOfWeek := time.Date(p[i].Year, time.Month(p[i].Month), p[i].Day, 0, 0, 0, 0, jst).Weekday()
		entryToInsert[i] = Kakeibo{
			Day:            p[i].Day,
			Month:          p[i].Month,
			Year:           p[i].Year,
			DayOfWeek:      int(dayOfWeek),
			Category:       p[i].Category,
			Type:           EXPENSE,
			Price:          p[i].Price,
			Detail:         p[i].Detail,
			TermsOfPayment: CASH,
		}
	}
	//===
	//=== DB登録
	//===
	// Contextの作成
	ctx := appengine.NewContext(r)
	// 登録用のKeyをputNum分作成
	keys := make([]*datastore.Key, putNum)
	for i := 0; i < putNum; i++ {
		keys[i] = datastore.NewIncompleteKey(ctx, KAKEIBO_ENTRY,
			getMonthKey(ctx, familyKakeibo, entryToInsert[i].Year, entryToInsert[i].Month))
	}
	// トランザクションで更新
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		_, e := datastore.PutMulti(c, keys, entryToInsert)
		return e
	}, nil)
	if err != nil {
		// トランザクションがエラーならerrを返却
		return err
	}
	return nil
}

// 一覧画面に表示するデータの検索を行う関数
func ReadList(r *http.Request, year, month int) (*ParamToShowList, error) {
	var err error
	//===
	//=== エントリの取得
	//===
	// contextの作成
	ctx := appengine.NewContext(r)
	// queryの作成
	// 指定した家計簿中の年・月で指定したエントリを全て表示
	query := datastore.NewQuery(KAKEIBO_ENTRY).
		Ancestor(getMonthKey(ctx, familyKakeibo, year, month)).
		Order("-Day")
	// queryの結果格納用のsliceの作成
	var tempArray []Kakeibo
	// queryの実行結果のキー格納用のsliceの作成
	var keys []*datastore.Key

	// トランザクションでクエリを実行
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		var e error
		keys, e = query.GetAll(c, &tempArray)
		return e
	}, nil)
	if err != nil {
		// エラーならnilを返す
		return nil, err
	}
	//===
	//=== エントリの出力
	//===
	// 出力用の構造体の作成
	p := new(ParamToShowList)
	p.Year = strconv.Itoa(year)
	if month < 10 {
		// 10以下の場合は0詰め
		p.Month = "0" + strconv.Itoa(month)
	} else {
		p.Month = strconv.Itoa(month)
	}
	p.KakeiboToShow = make([]WebKakeibo, len(keys))
	// tempの内容とkeyを出力用の配列に格納
	for i := 0; i < len(keys); i++ {
		// Key
		p.KakeiboToShow[i].DatastoreId = strconv.FormatInt(keys[i].IntID(), 10)
		// 年
		p.KakeiboToShow[i].Year = strconv.Itoa(tempArray[i].Year)
		// 月
		if tempArray[i].Month < 10 {
			// 10以下の場合は0詰め
			p.KakeiboToShow[i].Month = "0" + strconv.Itoa(tempArray[i].Month)
		} else {
			p.KakeiboToShow[i].Month = strconv.Itoa(tempArray[i].Month)
		}
		// 日
		if tempArray[i].Day < 10 {
			// 10以下の場合は0詰め
			p.KakeiboToShow[i].Day = "0" + strconv.Itoa(tempArray[i].Day)
		} else {
			p.KakeiboToShow[i].Day = strconv.Itoa(tempArray[i].Day)
		}
		// 曜日
		p.KakeiboToShow[i].DayOfWeek = getWeekString(tempArray[i].DayOfWeek)
		// 費目
		if v, ok := Categories[tempArray[i].Category]; ok {
			p.KakeiboToShow[i].Category = v.Name
		} else {
			// 対応する費目が見つからなかった場合、空白
			p.KakeiboToShow[i].Category = ""
		}
		// 費目のインデックス
		p.KakeiboToShow[i].CategoryIndex = tempArray[i].Category
		p.KakeiboToShow[i].Detail = tempArray[i].Detail
		p.KakeiboToShow[i].Price = strconv.Itoa(tempArray[i].Price)
	}
	return p, nil
}

// まとめ画面に表示するデータの検索を行う関数
func ReadSummary(r *http.Request, year int) (*ParamToShowSummary, error) {
	var err error
	//===
	//=== 表示用のパラメータの初期化
	//===
	p := new(ParamToShowSummary)
	p.SumOfMonth = make([]int, 12)
	p.OthersOfMonth = make([]int, 12)
	p.Results = make(map[int]*ResultOfMonth)
	for k, v := range Categories {
		if v.IsCalcSummary {
			// 計算対象となっている費目だけResults[]を作成
			p.Results[k] = &ResultOfMonth{
				Name:    v.Name,
				Summary: make([]int, 12),
			}
		}
	}
	p.Year = year
	p.YearList = make([]int, 0)
	// いったんUTCで現在時刻をフォーマットして、JSTに変換
	thisYear := time.Now().UTC().In(jst).Year()
	// 年のリストに表示する年を格納
	for y := thisYear; y >= STARTYEAR; y-- {
		p.YearList = append(p.YearList, y)
	}
	//===
	//=== エントリの取得
	//===
	// contextの作成
	ctx := appengine.NewContext(r)
	// queryの作成
	// 指定した家計簿中の年・月で指定したエントリを全て表示
	query := datastore.NewQuery(KAKEIBO_ENTRY).
		Ancestor(getKakeiboKey(ctx, familyKakeibo)).
		Filter("Year =", p.Year)
	// queryの結果格納用のsliceの作成
	var tempArray []Kakeibo
	// トランザクションでクエリを実行
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		_, e := query.GetAll(c, &tempArray)
		return e
	}, nil)
	if err != nil {
		// エラーの場合errを返却
		return nil, err
	}
	//===
	//=== 集計
	//===
	// queryの結果を月と費目ごとに集計
	for _, v := range tempArray {
		if v.Type == EXPENSE {
			// 金額を合計金額に加算
			p.SumOfMonth[v.Month-1] += v.Price
			// 費目別に金額を加算
			if _, ok := Categories[v.Category]; ok {
				if Categories[v.Category].IsCalcSummary {
					// 集計対象としている費目であれば、それに加算
					p.Results[v.Category].Summary[v.Month-1] += v.Price
				} else {
					// 集計対象以外の費目であれば、それらをまとめる
					p.OthersOfMonth[v.Month-1] += v.Price
				}
			}
		}
	}
	return p, nil
}

// DBのデータ更新を行う関数
func Update(r *http.Request, p *ParamToUpdate) error {
	var err error
	//===
	//=== 更新対象のエンティティの取得
	//===
	// Contextの作成
	ctx := appengine.NewContext(r)
	// 更新対象のエンティティを格納するkakeiboEntry
	var entryToUpdate Kakeibo
	// 更新対象のエンティティのキーの取得
	keyForUpdate := datastore.NewKey(ctx, KAKEIBO_ENTRY, "", int64(p.ID),
		getMonthKey(ctx, familyKakeibo, p.YearBefore, p.MonthBefore))
	// 更新対象のエンティティの取得
	if err = datastore.Get(ctx, keyForUpdate, &entryToUpdate); err != nil {
		return err
	}
	// 更新対象のエンティティを格納するkakeiboEntry
	dayOfWeek := time.Date(p.Year, time.Month(p.Month), p.Day, 0, 0, 0, 0, jst).Weekday()
	// 各値の更新
	entryToUpdate.Year = p.Year
	entryToUpdate.Month = p.Month
	entryToUpdate.Day = p.Day
	entryToUpdate.DayOfWeek = int(dayOfWeek)
	entryToUpdate.Category = p.Category
	entryToUpdate.Detail = p.Detail
	entryToUpdate.Price = p.Price

	//===
	//=== DB更新
	//===
	// エンティティの更新
	if p.Year == p.YearBefore && p.Month == p.MonthBefore {
		// 年月が変更になっていない場合、キーは変更せず更新
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			// 更新をトランザクションで実行
			_, e := datastore.Put(ctx, keyForUpdate, &entryToUpdate)
			return e
		}, nil)
	} else {
		// 年月が変更になってる場合、キーを削除して新たに作成
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			// 削除と作成をトランザクションで実行
			e := datastore.Delete(ctx, keyForUpdate)
			if e != nil {
				return e
			}
			newKey := datastore.NewIncompleteKey(ctx, KAKEIBO_ENTRY, getMonthKey(ctx, familyKakeibo, p.Year, p.Month))
			_, e = datastore.Put(ctx, newKey, &entryToUpdate)
			return e
		}, nil)
	}
	return err
}

// DBからデータ削除を行う関数
func Delete(r *http.Request, id, year, month int) error {
	var err error
	//===
	//=== エントリの削除
	//===
	// Contextの作成
	ctx := appengine.NewContext(r)
	// 削除対象のエンティティのキーの取得
	keyForDelete := datastore.NewKey(ctx, KAKEIBO_ENTRY, "", int64(id), getMonthKey(ctx, familyKakeibo, year, month))
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		e := datastore.Delete(c, keyForDelete)
		return e
	}, nil)
	if err != nil {
		return err
	}
	// エラーがなければnilを返す
	return nil
}

// データを実際に削除する(容量確保のため)
// データの一括削除をした場合に実行する関数
func DeleteMonth(r *http.Request, year, month int) error {
	var err error
	//===
	//=== 削除対象のキーの取得
	//===
	// Contextの作成
	ctx := appengine.NewContext(r)
	// 削除対象のエンティティのキーを格納するスライス
	var keyForDelete []*datastore.Key
	if month == ALL {
		// 一括削除の場合
		queries := make([]*datastore.Query, 12)
		for m := 0; m < 12; m++ {
			// 各月のキーを検索するクエリの作成
			queries[m] = datastore.NewQuery(KAKEIBO_ENTRY).
				Ancestor(getMonthKey(ctx, familyKakeibo, year, m+1)).KeysOnly()
		}
		// 各月のキーを検索
		keys := make([][]*datastore.Key, 12)
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			var e error
			for m := 0; m < 12; m++ {
				keys[m], e = queries[m].GetAll(c, nil)
				if e != nil {
					// GetAllがうまくいかなかったらbreakしてすぐにエラー値を返却
					break
				}
			}
			return e
		}, nil)
		if err != nil {
			// 検索がエラーの場合
			return err
		} else {
			for _, k := range keys {
				// 各月のキーを削除対象のキーに追加
				keyForDelete = append(keyForDelete, k...)
			}
		}
	} else {
		// 各月削除の場合
		// 削除対象のキーを検索するクエリの作成
		query := datastore.NewQuery(KAKEIBO_ENTRY).
			Ancestor(getMonthKey(ctx, familyKakeibo, year, month)).KeysOnly()
		// 削除する月のキーを検索
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			var e error
			keyForDelete, e = query.GetAll(c, nil)
			return e
		}, nil)
		if err != nil {
			return err
		}
	}
	//===
	//=== 削除処理
	//===
	// トランザクションによるキーの削除
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		e := datastore.DeleteMulti(c, keyForDelete)
		return e
	}, nil)
	if err != nil {
		return err
	}
	// エラーがなければnilを返す
	return nil
}

// 家計簿の名前に対応するキーを返す関数
func getKakeiboKey(c context.Context, kakeiboName string) *datastore.Key {
	return datastore.NewKey(c, KAKEIBO_ID, kakeiboName, 0, nil)
}

// 年に対応するキーを返す関数
func getYearKey(c context.Context, kakeiboName string, year int) *datastore.Key {
	return datastore.NewKey(c, KAKEIBO_YEAR, "", int64(year), getKakeiboKey(c, kakeiboName))
}

// 月に対応するキーを返す関数
func getMonthKey(c context.Context, kakeiboName string, year, month int) *datastore.Key {
	return datastore.NewKey(c, KAKEIBO_MONTH, "", int64(month), getYearKey(c, kakeiboName, year))
}

