package main

type Source struct {
	Type       string `firestore:"type"`
	URL        string `firestore:"url"`
	AccessedAt string `firestore:"accessedAt"`
	Note       string `firestore:"note,omitempty"`
}

type ZooSeed struct {
	ID          string
	Name        string
	Pref        string
	City        string
	OfficialURL string
	Lat         float64
	Lng         float64
	Sources     []Source
}

type KoalaSeed struct {
	ID           string
	Name         string
	Sex          string
	BirthDate    string
	Status       string // alive / deceased / unknown
	CurrentZooID string
	MotherID     string
	FatherID     string
	PhotoURL     string
	Sources      []Source
}

func (z ZooSeed) toDoc(lastVerifiedAt string) map[string]any {
	return map[string]any{
		"name":           z.Name,
		"pref":           z.Pref,
		"city":           z.City,
		"status":         "open",
		"officialUrl":    z.OfficialURL,
		"lastVerifiedAt": lastVerifiedAt,
		"geo":            map[string]any{"lat": z.Lat, "lng": z.Lng},
		"sources":        toSourceMaps(z.Sources),
	}
}

func (k KoalaSeed) toDoc(lastVerifiedAt string) map[string]any {
	doc := map[string]any{
		"name":           k.Name,
		"sex":            k.Sex,
		"birthDate":      k.BirthDate,
		"status":         k.Status,
		"currentZooId":   k.CurrentZooID,
		"motherId":       k.MotherID,
		"fatherId":       k.FatherID,
		"lastVerifiedAt": lastVerifiedAt,
		"sources":        toSourceMaps(k.Sources),
	}
	if k.PhotoURL != "" {
		doc["photoUrl"] = k.PhotoURL
	}
	return doc
}

func toSourceMaps(ss []Source) []map[string]any {
	out := make([]map[string]any, 0, len(ss))
	for _, s := range ss {
		m := map[string]any{
			"type":       s.Type,
			"url":        s.URL,
			"accessedAt": s.AccessedAt,
		}
		if s.Note != "" {
			m["note"] = s.Note
		}
		out = append(out, m)
	}
	return out
}

func src(today, typ, url string) Source {
	return Source{Type: typ, URL: url, AccessedAt: today}
}

func buildZoos(today string) []ZooSeed {
	return []ZooSeed{
		{
			ID:          "higashiyama-zoo",
			Name:        "東山動植物園",
			Pref:        "愛知県",
			City:        "名古屋市",
			OfficialURL: "https://www.higashiyama.city.nagoya.jp/",
			Lat:         35.1567406,
			Lng:         136.9811452,
			Sources: []Source{
				src(today, "official", "https://www.higashiyama.city.nagoya.jp/"),
			},
		},
		{
			ID:          "hirakawa-zoo",
			Name:        "平川動物公園",
			Pref:        "鹿児島県",
			City:        "鹿児島市",
			OfficialURL: "https://hirakawazoo.jp/",
			Lat:         31.463123,
			Lng:         130.503422,
			Sources: []Source{
				src(today, "official", "https://hirakawazoo.jp/virtual/koara/"),
			},
		},
		{
			ID:          "england-hill",
			Name:        "淡路ファームパーク イングランドの丘",
			Pref:        "兵庫県",
			City:        "南あわじ市",
			OfficialURL: "https://www.england-hill.com/",
			Lat:         34.3071605,
			Lng:         134.8017872,
			Sources: []Source{
				src(today, "official", "https://www.england-hill.com/"),
			},
		},
		{
			ID:          "oji-zoo",
			Name:        "神戸市立王子動物園",
			Pref:        "兵庫県",
			City:        "神戸市",
			OfficialURL: "https://www.kobe-ojizoo.jp/",
			Lat:         34.7099614,
			Lng:         135.2145663,
			Sources: []Source{
				src(today, "official", "https://www.kobe-ojizoo.jp/"),
			},
		},
		{
			ID:          "kanazawa-zoo",
			Name:        "横浜市立金沢動物園",
			Pref:        "神奈川県",
			City:        "横浜市",
			OfficialURL: "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/",
			Lat:         35.349729,
			Lng:         139.597359,
			Sources: []Source{
				src(today, "official", "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/"),
			},
		},
		{
			ID:          "saitama-childrens-zoo",
			Name:        "埼玉県こども動物自然公園",
			Pref:        "埼玉県",
			City:        "東松山市",
			OfficialURL: "https://www.parks.or.jp/sczoo/",
			Lat:         36.000865,
			Lng:         139.374657,
			Sources: []Source{
				src(today, "official", "https://www.parks.or.jp/sczoo/"),
			},
		},
		{
			ID:          "tama-zoo",
			Name:        "多摩動物公園",
			Pref:        "東京都",
			City:        "日野市",
			OfficialURL: "https://www.tokyo-zoo.net/zoo/tama/",
			Lat:         35.649538,
			Lng:         139.402201,
			Sources: []Source{
				src(today, "official", "https://www.tokyo-zoo.net/zoo/tama/"),
			},
		},
	}
}

func buildKoalas(today string) []KoalaSeed {
	return []KoalaSeed{
		// =========================
		// Higashiyama（既存互換）
		// =========================
		{ID: "monaka", Name: "もなか", Sex: "M", BirthDate: "2024-08-??", Status: "alive", CurrentZooID: "higashiyama-zoo", MotherID: "rin", FatherID: "ishin",
			Sources: []Source{src(today, "official", "https://www.higashiyama.city.nagoya.jp/")}},
		{ID: "daifuku", Name: "だいふく", Sex: "M", BirthDate: "2022-03-14", Status: "alive", CurrentZooID: "higashiyama-zoo", MotherID: "rin", FatherID: "taichi",
			Sources: []Source{src(today, "official", "https://www.higashiyama.city.nagoya.jp/")}},
		{ID: "tsukushi", Name: "つくし", Sex: "F", BirthDate: "2020-??-??", Status: "alive", CurrentZooID: "higashiyama-zoo", MotherID: "rin", FatherID: "ishin",
			Sources: []Source{src(today, "official", "https://www.higashiyama.city.nagoya.jp/")}},
		{ID: "rin", Name: "りん", Sex: "F", BirthDate: "2017-??-??", Status: "unknown", CurrentZooID: "", MotherID: "tilly", FatherID: "",
			Sources: []Source{src(today, "official", "https://www.higashiyama.city.nagoya.jp/")}},
		{ID: "ishin", Name: "イシン", Sex: "M", BirthDate: "????-??-??", Status: "unknown", CurrentZooID: "", MotherID: "", FatherID: "",
			Sources: []Source{src(today, "official", "https://www.higashiyama.city.nagoya.jp/")}},
		{ID: "taichi", Name: "タイチ", Sex: "M", BirthDate: "????-??-??", Status: "unknown", CurrentZooID: "", MotherID: "", FatherID: "",
			Sources: []Source{src(today, "official", "https://www.higashiyama.city.nagoya.jp/")}},
		{ID: "tilly", Name: "ティリー", Sex: "F", BirthDate: "????-??-??", Status: "unknown", CurrentZooID: "", MotherID: "", FatherID: "",
			Sources: []Source{src(today, "official", "https://www.higashiyama.city.nagoya.jp/")}},

		// =========================
		// Hirakawa（公式一覧が強い）
		// =========================
		{ID: "hirakawa-yume", Name: "ユメ", Sex: "F", BirthDate: "2015-01-02", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-rio", Name: "リオ", Sex: "F", BirthDate: "2015-12-10", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-himawari", Name: "ヒマワリ", Sex: "F", BirthDate: "2019-06-22", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-kibou", Name: "キボウ", Sex: "F", BirthDate: "2019-10-17", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-indico", Name: "インディコ", Sex: "F", BirthDate: "2019-12-22", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-hinata", Name: "ヒナタ", Sex: "F", BirthDate: "2021-03-27", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-light", Name: "ライト", Sex: "M", BirthDate: "2021-05-14", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-kanae", Name: "カナエ", Sex: "F", BirthDate: "2021-06-13", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-peace", Name: "ピース", Sex: "F", BirthDate: "2021-06-29", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-tsukushi", Name: "つくし", Sex: "F", BirthDate: "2022-09-03", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-archer", Name: "アーチャー", Sex: "M", BirthDate: "2019-04-26", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-nozomu", Name: "ノゾム", Sex: "M", BirthDate: "2022-02-26", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-taiyou", Name: "タイヨウ", Sex: "M", BirthDate: "2022-08-11", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-asahi", Name: "アサヒ", Sex: "M", BirthDate: "2022-12-09", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-tsumugi", Name: "ツムギ", Sex: "F", BirthDate: "2023-05-13", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-arata", Name: "アラタ", Sex: "M", BirthDate: "2023-06-14", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-chabou", Name: "チャーボウ", Sex: "M", BirthDate: "2023-10-21", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-star", Name: "スター", Sex: "M", BirthDate: "2023-11-18", Status: "alive", CurrentZooID: "hirakawa-zoo",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},
		{ID: "hirakawa-himawari-baby-20240601", Name: "ヒマワリの仔", Sex: "F", BirthDate: "2024-06-01", Status: "alive", CurrentZooID: "hirakawa-zoo", MotherID: "hirakawa-himawari",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/zookan/%E3%82%B3%E3%82%A2%E3%83%A9/")}},

		// 平川：来園個体「ししお」（両親まで公式）
		{ID: "hirakawa-shishio", Name: "ししお", Sex: "M", BirthDate: "2022-04-04", Status: "alive", CurrentZooID: "hirakawa-zoo", MotherID: "hirakawa-holly", FatherID: "hirakawa-taichi",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/2025/04/05/%E3%82%B3%E3%82%A2%E3%83%A9%E3%81%AE%E4%B8%80%E8%88%AC%E5%85%AC%E9%96%8B%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6-3/")}},
		{ID: "hirakawa-taichi", Name: "タイチ", Sex: "M", BirthDate: "????-??-??", Status: "unknown",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/2025/04/05/%E3%82%B3%E3%82%A2%E3%83%A9%E3%81%AE%E4%B8%80%E8%88%AC%E5%85%AC%E9%96%8B%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6-3/")}},
		{ID: "hirakawa-holly", Name: "ホリー", Sex: "F", BirthDate: "????-??-??", Status: "unknown",
			Sources: []Source{src(today, "official", "https://hirakawazoo.jp/2025/04/05/%E3%82%B3%E3%82%A2%E3%83%A9%E3%81%AE%E4%B8%80%E8%88%AC%E5%85%AC%E9%96%8B%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6-3/")}},

		// =========================
		// England Hill（公式/県資料が強い）
		// =========================
		{ID: "england-nozomi", Name: "のぞみ", Sex: "F", BirthDate: "2008-03-01", Status: "alive", CurrentZooID: "england-hill",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},
		{ID: "england-daichi", Name: "だいち", Sex: "M", BirthDate: "2013-08-18", Status: "alive", CurrentZooID: "england-hill", MotherID: "england-yume", FatherID: "england-ark",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},
		{ID: "england-yume", Name: "ゆめ", Sex: "F", BirthDate: "????-??-??", Status: "unknown",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},
		{ID: "england-ark", Name: "アーク", Sex: "M", BirthDate: "????-??-??", Status: "unknown",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},
		{ID: "england-umi", Name: "ウミ", Sex: "F", BirthDate: "2014-06-13", Status: "alive", CurrentZooID: "england-hill",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},
		{ID: "england-peter", Name: "ピーター", Sex: "M", BirthDate: "2016-03-28", Status: "alive", CurrentZooID: "england-hill",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},
		{ID: "england-nagi", Name: "ナギ", Sex: "F", BirthDate: "2023-07-31", Status: "alive", CurrentZooID: "england-hill", MotherID: "england-umi", FatherID: "england-peter",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},
		// 借入で来る北方系オス「ノゾム」
		{ID: "england-nozomu", Name: "ノゾム", Sex: "M", BirthDate: "2022-02-26", Status: "alive", CurrentZooID: "england-hill",
			Sources: []Source{src(today, "official", "https://web.pref.hyogo.lg.jp/nk12/press/20250305.html")}},

		// =========================
		// Oji（赤ちゃん＋両親）
		// =========================
		{ID: "oji-ouka", Name: "オウカ(桜花)", Sex: "F", BirthDate: "????-??-??", Status: "alive", CurrentZooID: "oji-zoo",
			Sources: []Source{src(today, "official", "https://www.kobe-ojizoo.jp/info/detail/?id=708")}},
		{ID: "oji-ibuki", Name: "いぶき", Sex: "M", BirthDate: "????-??-??", Status: "alive", CurrentZooID: "oji-zoo",
			Sources: []Source{src(today, "official", "https://www.kobe-ojizoo.jp/info/detail/?id=708")}},
		{ID: "oji-ouki", Name: "おうき(桜希)", Sex: "M", BirthDate: "2024-06-12", Status: "alive", CurrentZooID: "oji-zoo", MotherID: "oji-ouka", FatherID: "oji-ibuki",
			Sources: []Source{src(today, "official", "https://www.kobe-ojizoo.jp/info/detail/?id=708")}},

		// 2019年誕生の赤ちゃん2頭（親と誕生日は別資料が必要なので、まず名前だけ確定で入れる）
		{ID: "oji-mai", Name: "マイ", Sex: "U", BirthDate: "2019-??-??", Status: "unknown", CurrentZooID: "oji-zoo",
			Sources: []Source{src(today, "official", "https://www.kobe-ojizoo.jp/info/detail/?id=363")}},
		{ID: "oji-hana", Name: "ハナ", Sex: "U", BirthDate: "2019-??-??", Status: "unknown", CurrentZooID: "oji-zoo",
			Sources: []Source{src(today, "official", "https://www.kobe-ojizoo.jp/info/detail/?id=363")}},

		// =========================
		// Kanazawa（家系が作れる）
		// =========================
		{ID: "kanazawa-koharu", Name: "コハル", Sex: "F", BirthDate: "????-??-??", Status: "alive", CurrentZooID: "kanazawa-zoo",
			Sources: []Source{src(today, "official", "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/1/2023/01/15286f8ab991421b1be513c4b9a51d9169b2b197.pdf")}},
		{ID: "kanazawa-botan", Name: "ぼたん", Sex: "F", BirthDate: "????-??-??", Status: "unknown", CurrentZooID: "kanazawa-zoo",
			Sources: []Source{src(today, "official", "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/1/20220815092026-a4dd7812b3fadf993e6febb5a4370b959070dd04.pdf")}},
		{ID: "kanazawa-charlie", Name: "チャーリー", Sex: "M", BirthDate: "????-??-??", Status: "alive", CurrentZooID: "kanazawa-zoo",
			Sources: []Source{src(today, "official", "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/1/2023/01/15286f8ab991421b1be513c4b9a51d9169b2b197.pdf")}},
		{ID: "kanazawa-harry", Name: "ハリー", Sex: "M", BirthDate: "2022-04-22", Status: "alive", CurrentZooID: "kanazawa-zoo", MotherID: "kanazawa-koharu", FatherID: "kanazawa-charlie",
			Sources: []Source{src(today, "official", "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/1/2023/01/15286f8ab991421b1be513c4b9a51d9169b2b197.pdf")}},
		{ID: "kanazawa-hinagiku", Name: "ひなぎく", Sex: "F", BirthDate: "2021-12-15", Status: "deceased", CurrentZooID: "kanazawa-zoo", MotherID: "kanazawa-botan", FatherID: "kanazawa-charlie",
			Sources: []Source{
				src(today, "official", "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/1/20220815092026-a4dd7812b3fadf993e6febb5a4370b959070dd04.pdf"),
				src(today, "official", "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/author97ae2/2025/07/423789a2bf077c7959f2cbd4a1758d8d7ec6951e.pdf"),
			},
		},

		// =========================
		// Tama（公式で取れる範囲から）
		// =========================
		{ID: "tama-mirai", Name: "ミライ", Sex: "F", BirthDate: "2008-10-25", Status: "alive", CurrentZooID: "tama-zoo",
			Sources: []Source{src(today, "official", "https://www.tokyo-zoo.net/topic/topics_detail?inst=tama&kind=news&link_num=24266")}},
		{ID: "tama-komachi", Name: "こまち", Sex: "F", BirthDate: "2017-04-27", Status: "alive", CurrentZooID: "tama-zoo",
			Sources: []Source{src(today, "official", "https://www.spt.metro.tokyo.lg.jp/tosei/hodohappyo/press/2020/07/13/03.html")}},
	}
}
