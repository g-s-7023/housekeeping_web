package lib

import (
	"time"
)

const (
	//===
	//=== 一覧(Input.html)に表示する行数
	//===
	INPUTLINES = 15
	//===
	//=== htmlファイルのディレクトリ
	DIR_HTML = "html/"
	//===
	//===
	//=== 使いはじめる年
	//===
	STARTYEAR = 2017
	//===
	//=== カインド
	//===
	// ルートエンティティのカインド(ユーザごとに異なる家計簿のID)
	KAKEIBO_ID = "KakeiboID"
	// 年エンティティのカインド(親キー：KakeiboIDカインドのキー)
	KAKEIBO_YEAR = "KakeiboYear"
	// 月エンティティのカインド(親キー：KakeiboYearカインドのキー)
	KAKEIBO_MONTH = "KakeiboMonth"
	// 家計簿の各エントリが入るカインド(親キー：KakeiboMonthカインドのキー)
	KAKEIBO_ENTRY = "KakeiboEntry"
	//===
	//=== typeの値
	//===
	INCOME  = 0
	EXPENSE = 1
	//===
	//=== termsOfPaymentの値
	//===
	CASH = 0
	CARD = 1
	//===
	//=== SQLiteに格納用のフラグの値 (isSynchronized、isDeleted)
	//===
	FALSE = 0
	TRUE  = 1
	//===
	//=== 一括削除する際の月の値
	//===
	ALL = 13
	//===
	//=== 各費目に該当する番号
	//===
	NONE          = 0
	DINEOUT       = 1
	HOMEEATING    = 2
	BUSINESS      = 3
	CLOTHES       = 4
	HOBBY         = 5
	TRAFFIC       = 6
	RENT          = 7
	HOUSEHOLD     = 8
	GOODS         = 9
	ELECTRICITY   = 10
	GAS           = 11
	WATER         = 12
	T_CELLPHONE   = 13
	N_CELLPHONE   = 14
	INTERNET      = 15
	HEALTHCARE    = 16
	HAIRSALON     = 17
	CHILD         = 18
	INSURANCE     = 19
	ASAHI         = 20
	MEMBERSHIP    = 21
	ENTERTAINMENT = 22
	OTHER         = 99
)

// DBに登録・参照するための構造体
type Kakeibo struct {
	Id              int
	Day             int
	Month           int
	Year            int
	DayOfWeek       int
	Category        int
	Type            int
	Price           int
	Detail          string
	TermsOfPayment  int
	IsDeleted       int
	IsSynchronized  int
	LastUpdatedDate string
}

// 一覧ページで出力するために各値を格納する構造体
type WebKakeibo struct {
	DatastoreId    string
	Day            string
	Month          string
	Year           string
	DayOfWeek      string
	Category       string
	CategoryIndex  int
	Type           int
	Price          string
	Detail         string
	TermsOfPayment int
}

// 入力ページで出力するために各値を格納する構造体
type ParamToShowInput struct {
	CategoryList map[int]string
	Lines        []int
	HasAnyError  bool
	HasError     []bool
	Date         []string
	Category     []int
	Detail       []string
	Price        []string
}

// DB登録を行う関数に渡す構造体
type ParamToInsert struct {
	Day      int
	Month    int
	Year     int
	Category int
	Price    int
	Detail   string
}

// 更新ページを出力するために書く値を格納する構造体
type ParamToShowUpdate struct {
	CategoryList  map[int]string
	HasAnyError   bool
	DatastoreId   string
	Date          string
	Category      string
	CategoryIndex int
	Price         string
	Detail        string
}

// DB更新を行う関数に渡す構造体
type ParamToUpdate struct {
	ID          int
	Day         int
	Month       int
	MonthBefore int
	Year        int
	YearBefore  int
	Category    int
	Detail      string
	Price       int
}

// 一覧を表示するためにhtmlTemplateに渡す構造体
type ParamToShowList struct {
	Year          string
	Month         string
	KakeiboToShow []WebKakeibo
}

// 費目とまとめの計算有無
type KakeiboCategory struct {
	Name          string
	IsCalcSummary bool
}

// 費目の名称と各月の金額を格納する構造体
type ResultOfMonth struct {
	// 費目に対応する文字列
	Name string
	// 各月の金額
	Summary []int
}

// まとめを表示するためにhtmlTemplateに渡す構造体
type ParamToShowSummary struct {
	// 表示する年
	Year int
	// 集計対象の年のリスト
	YearList []int
	// 各月の合計金額を格納
	SumOfMonth []int
	// 各月の集計対象以外の金額を格納
	OthersOfMonth []int
	// 表示対象の費目に対応した文字列と各月の金額
	Results map[int]*ResultOfMonth
}

var (
	//===
	//=== 家計簿IDカインドのキー
	//===
	myKakeibo     = "myKakeibo"
	familyKakeibo = "familyKakeibo"
	//===
	//=== タイムゾーン
	//===
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	//===
	//=== 費目で選択可能な文字列とまとめの計算有無
	//===
	Categories = map[int]KakeiboCategory{
		NONE:          {"", false},
		DINEOUT:       {"外食", true},
		HOMEEATING:    {"家庭食", true},
		BUSINESS:      {"仕事関連", false},
		CLOTHES:       {"衣類・雑貨", false},
		HOBBY:         {"趣味・娯楽", false},
		TRAFFIC:       {"交通費", false},
		RENT:          {"家賃", true},
		HOUSEHOLD:     {"家具・家電", true},
		GOODS:         {"日用雑貨", true},
		ELECTRICITY:   {"電気代", true},
		GAS:           {"ガス代", true},
		WATER:         {"水道代", true},
		T_CELLPHONE:   {"T携帯", true},
		N_CELLPHONE:   {"N携帯", true},
		INTERNET:      {"ネット", true},
		HEALTHCARE:    {"医療費", true},
		HAIRSALON:     {"美容院", false},
		CHILD:         {"養育", true},
		ENTERTAINMENT: {"交際", false},
		INSURANCE:     {"保険", false},
		ASAHI:         {"Asahi", true},
		MEMBERSHIP:    {"年会費", true},
		OTHER:         {"雑費", false},
	}
)

// 曜日に対応する文字列を返す関数
func getWeekString(w int) string {
	wday := []string{"(日)", "(月)", "(火)", "(水)", "(木)", "(金)", "(土)"}
	if w >= 0 && w < 7 {
		return wday[w]
	} else {
		return ""
	}
}
