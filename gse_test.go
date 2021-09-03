package gse

import (
	"testing"

	"github.com/vcaesar/tt"
)

func init() {
	prodSeg.LoadDict()
	prodSeg.LoadStop("zh")
}

func TestLoadDictMap(t *testing.T) {
	m := []map[string]string{
		{
			"text":      "一城山水",
			"frequency": "10",
			"pos":       "n",
		},
		{
			"text":      "山河日月",
			"frequency": "13",
		},
	}

	prodSeg.LoadDictMap(m)

	f, pos, ok := prodSeg.Find("一城山水")
	tt.Bool(t, ok)
	tt.Equal(t, pos, "n")
	tt.Equal(t, 10, f)

	f, _, ok = prodSeg.Find("山河日月")
	tt.Bool(t, ok)
	tt.Equal(t, 13, f)
}

func TestAnalyze(t *testing.T) {
	txt := `城市地标建筑: 纽约帝国大厦, 旧金山湾金门大桥, Seattle Space Needle, Toronto CN Tower, 伦敦大笨钟`

	s := prodSeg.Cut(txt, true)
	tt.Equal(t, 23, len(s))
	tt.Equal(t, "[城市地标 建筑 :  纽约 帝国大厦 ,  旧金山湾 金门大桥 ,  Seattle   Space   Needle ,  Toronto   CN   Tower ,  伦敦 大笨钟]", s)

	a := prodSeg.Analyze(s)
	tt.Equal(t, 23, len(a))
	tt.Equal(t, "[{0 4 0 0  城市地标 3} {4 6 1 0  建筑 14397} {6 8 2 0  :  0} {8 10 3 0  纽约 1758} {10 14 4 0  帝国大厦 3} {14 16 5 0  ,  0} {16 20 6 0  旧金山湾 3} {20 24 7 0  金门大桥 38} {24 26 8 0  ,  0} {26 33 9 0  Seattle 0} {33 34 10 0    0} {34 39 11 0  Space 0} {39 40 12 0    0} {40 46 13 0  Needle 0} {46 48 14 0  ,  0} {48 55 15 0  Toronto 0} {55 56 16 0    0} {56 58 17 0  CN 0} {58 59 18 0    0} {59 64 19 0  Tower 0} {64 66 20 0  ,  0} {66 68 21 0  伦敦 2255} {68 71 22 0  大笨钟 0}]", a)

	tt.Equal(t, 0, a[0].Start)
	tt.Equal(t, 4, a[0].End)
	tt.Equal(t, 0, a[0].Position)
	tt.Equal(t, 0, a[0].Len)
	tt.Equal(t, "城市地标", a[0].Text)
	tt.Equal(t, 3, a[0].Freq)
	tt.Equal(t, "", a[0].Type)
}

func TestHMM(t *testing.T) {
	tt.Equal(t, 587880, len(prodSeg.Dict.Tokens))
	tt.Equal(t, 5.3250742e+07, prodSeg.Dict.totalFrequency)

	hmm := prodSeg.HMMCutMod("纽约时代广场")
	tt.Equal(t, 2, len(hmm))
	tt.Equal(t, "纽约", hmm[0])
	tt.Equal(t, "时代广场", hmm[1])

	// text := "纽约时代广场, 纽约帝国大厦, 旧金山湾金门大桥"
	tx := prodSeg.Cut(text, true)
	tt.Equal(t, 7, len(tx))
	tt.Equal(t, "[纽约时代广场 ,  纽约 帝国大厦 ,  旧金山湾 金门大桥]", tx)

	tx = prodSeg.TrimPunct(tx)
	tt.Equal(t, 5, len(tx))
	tt.Equal(t, "[纽约时代广场 纽约 帝国大厦 旧金山湾 金门大桥]", tx)

	tx = prodSeg.cutDAGNoHMM(text)
	tt.Equal(t, 9, len(tx))
	tt.Equal(t, "[纽约时代广场 ,   纽约 帝国大厦 ,   旧金山湾 金门大桥]", tx)

	tx = append(tx, " 广场")
	tx = append(tx, "ok👌", "的")
	tt.Bool(t, prodSeg.IsStop("的"))

	tx = prodSeg.Trim(tx)
	tt.Equal(t, 7, len(tx))
	tt.Equal(t, "[纽约时代广场 纽约 帝国大厦 旧金山湾 金门大桥 广场 ok]", tx)

	tx1 := prodSeg.CutTrim(text, true)
	tt.Equal(t, 5, len(tx1))
	tt.Equal(t, "[纽约时代广场 纽约 帝国大厦 旧金山湾 金门大桥]", tx1)

	s := prodSeg.CutStr(tx, ", ")
	tt.Equal(t, 80, len(s))
	tt.Equal(t, "纽约时代广场, 纽约, 帝国大厦, 旧金山湾, 金门大桥, 广场, ok", s)

	tx = prodSeg.CutAll(text)
	tt.Equal(t, 21, len(tx))
	tt.Equal(t,
		"[纽约 纽约时代广场 时代 时代广场 广场 ,   纽约 帝国 帝国大厦 国大 大厦 ,   旧金山 旧金山湾 金山 山湾 金门 金门大桥 大桥]",
		tx)

	tx = prodSeg.CutSearch(text, false)
	tt.Equal(t, 20, len(tx))
	tt.Equal(t,
		"[纽约 时代 广场 纽约时代广场 ,   纽约 帝国 国大 大厦 帝国大厦 ,   金山 山湾 旧金山 旧金山湾 金门 大桥 金门大桥]",
		tx)

	tx = prodSeg.CutSearch(text, true)
	tt.Equal(t, 18, len(tx))
	tt.Equal(t,
		"[纽约 时代 广场 纽约时代广场 ,  纽约 帝国 国大 大厦 帝国大厦 ,  金山 山湾 旧金山 旧金山湾 金门 大桥 金门大桥]",
		tx)

	f1 := prodSeg.SuggestFreq("西雅图")
	tt.Equal(t, 79, f1)

	f1 = prodSeg.SuggestFreq("西雅图", "西雅图都会区", "旧金山湾")
	tt.Equal(t, 0, f1)
}

func TestPos(t *testing.T) {
	s := prodSeg.String(text, true)
	tt.Equal(t, 206, len(s))
	tt.Equal(t,
		"纽约/ns 时代/n 广场/n 时代广场/n 纽约时代广场/nt ,/x  /x 纽约/ns 帝国/n 大厦/n 帝国大厦/nr ,/x  /x 金山/nr 旧金山/ns 湾/zg 旧金山湾/ns 金门/n 大桥/ns 金门大桥/nz ", s)

	c := prodSeg.Slice(text, true)
	tt.Equal(t, 20, len(c))

	pos := prodSeg.Pos(text, false)
	tt.Equal(t, 9, len(pos))
	tt.Equal(t,
		"[{纽约时代广场 nt} {, x} {  x} {纽约 ns} {帝国大厦 nr} {, x} {  x} {旧金山湾 ns} {金门大桥 nz}]", pos)

	pos = prodSeg.TrimPosPunct(pos)
	tt.Equal(t, 5, len(pos))
	tt.Equal(t, "[{纽约时代广场 nt} {纽约 ns} {帝国大厦 nr} {旧金山湾 ns} {金门大桥 nz}]", pos)

	pos = prodSeg.Pos(text, true)
	tt.Equal(t, 20, len(pos))
	tt.Equal(t,
		"[{纽约 ns} {时代 n} {广场 n} {时代广场 n} {纽约时代广场 nt} {, x} {  x} {纽约 ns} {帝国 n} {大厦 n} {帝国大厦 nr} {, x} {  x} {金山 nr} {旧金山 ns} {湾 zg} {旧金山湾 ns} {金门 n} {大桥 ns} {金门大桥 nz}]", pos)

	pos1 := prodSeg.PosTrim(text, true, "zg")
	tt.Equal(t, 15, len(pos1))
	tt.Equal(t,
		"[{纽约 ns} {时代 n} {广场 n} {时代广场 n} {纽约时代广场 nt} {纽约 ns} {帝国 n} {大厦 n} {帝国大厦 nr} {金山 nr} {旧金山 ns} {旧金山湾 ns} {金门 n} {大桥 ns} {金门大桥 nz}]", pos1)

	pos = append(pos, SegPos{Text: "👌", Pos: "x"})
	pos = prodSeg.TrimPos(pos)
	tt.Equal(t, 16, len(pos))
	tt.Equal(t,
		"[{纽约 ns} {时代 n} {广场 n} {时代广场 n} {纽约时代广场 nt} {纽约 ns} {帝国 n} {大厦 n} {帝国大厦 nr} {金山 nr} {旧金山 ns} {湾 zg} {旧金山湾 ns} {金门 n} {大桥 ns} {金门大桥 nz}]", pos)

	s = prodSeg.PosStr(pos, ", ")
	tt.Equal(t, 204, len(s))
	tt.Equal(t,
		"纽约/ns, 时代/n, 广场/n, 时代广场/n, 纽约时代广场/nt, 纽约/ns, 帝国/n, 大厦/n, 帝国大厦/nr, 金山/nr, 旧金山/ns, 湾/zg, 旧金山湾/ns, 金门/n, 大桥/ns, 金门大桥/nz", s)

	prodSeg.SkipPos = true
	s = prodSeg.PosStr(pos, ", ")
	tt.Equal(t, 162, len(s))
	tt.Equal(t,
		"纽约, 时代, 广场, 时代广场, 纽约时代广场, 纽约, 帝国, 大厦, 帝国大厦, 金山, 旧金山, 湾, 旧金山湾, 金门, 大桥, 金门大桥", s)

	pos = prodSeg.TrimWithPos(pos, "n", "zg")
	tt.Equal(t, 9, len(pos))
	tt.Equal(t,
		"[{纽约 ns} {纽约时代广场 nt} {纽约 ns} {帝国大厦 nr} {金山 nr} {旧金山 ns} {旧金山湾 ns} {大桥 ns} {金门大桥 nz}]", pos)

	pos2 := prodSeg.PosTrimArr(text, false, "n", "zg")
	tt.Equal(t, 5, len(pos2))
	tt.Equal(t,
		"[纽约时代广场 纽约 帝国大厦 旧金山湾 金门大桥]", pos2)

	pos3 := prodSeg.PosTrimStr(text, false, "n", "zg")
	tt.Equal(t, 64, len(pos3))
	tt.Equal(t,
		"纽约时代广场 纽约 帝国大厦 旧金山湾 金门大桥", pos3)
}

func TestStop(t *testing.T) {
	var seg Segmenter
	err := seg.LoadStop()
	tt.Nil(t, err)
	tt.Equal(t, 88, len(seg.StopWordMap))

	err = seg.LoadStop("testdata/stop.txt")
	tt.Nil(t, err)
	tt.Equal(t, 89, len(seg.StopWordMap))
	tt.Bool(t, seg.IsStop("离开"))

	// err := prodSeg.LoadStop("zh")
	// tt.Nil(t, err)
	tt.Equal(t, 1161, len(prodSeg.StopWordMap))

	b := prodSeg.IsStop("阿")
	tt.True(t, b)

	tt.True(t, prodSeg.IsStop("哎"))
	b = prodSeg.IsStop("的")
	tt.True(t, b)

	prodSeg.AddStop("lol")
	b = prodSeg.IsStop("lol")
	tt.True(t, b)

	prodSeg.RemoveStop("lol")
	b = prodSeg.IsStop("lol")
	tt.False(t, b)

	m := []string{"abc", "123"}
	prodSeg.LoadStopArr(m)
	tt.True(t, prodSeg.IsStop("abc"))
	tt.True(t, prodSeg.IsStop("123"))

	t1 := `hi, bot, 123; 🤖, 机器人; 👌^_^😆`
	s := FilterEmoji(t1)
	tt.Equal(t, "hi, bot, 123; , 机器人; ^_^", s)

	s = FilterSymbol(t1)
	tt.Equal(t, "hibot123机器人", s)

	s = FilterLang(t1, "Han")
	tt.Equal(t, "hibot机器人", s)

	t2 := `<p>test: </p> <div class="bot"> bot 机器人 <<银河系漫游指南>> </div>`
	s = FilterHtml(t2)
	tt.Equal(t, "test:   bot 机器人 <<银河系漫游指南>> ", s)

	s = prodSeg.CutTrimHtmls(t2, true)
	tt.Equal(t, "test bot 机器人 银河系 漫游 指南", s)

	s1 := Range("hibot, 机器人")
	tt.Equal(t, "[h i b o t ,   机 器 人]", s1)
	s = RangeText("hibot, 机器人")
	tt.Equal(t, "h i b o t ,   机 器 人 ", s)
}

func TestNum(t *testing.T) {
	seg, err := New("./testdata/test_dict3.txt")
	tt.Nil(t, err)

	seg.Num = true
	text := "t123test123 num123-1"
	s := seg.Cut(text)
	tt.Equal(t, "[t 1 2 3 test 1 2 3   num 1 2 3 - 1]", s)

	s = seg.CutAll(text)
	tt.Equal(t, "[t 1 2 3 t e s t 1 2 3   n u m 1 2 3 - 1]", s)

	seg.Alpha = true
	s = seg.CutSearch(text)
	tt.Equal(t, "[t 1 2 3 t e s t 1 2 3   n u m 1 2 3 - 1]", s)
}

func TestUrl(t *testing.T) {
	seg, err := New("./testdata/test_dict3.txt")
	tt.Nil(t, err)

	s1 := seg.CutUrls("https://www.g.com/search?q=test%m11.42&ie=UTF-8")
	tt.Equal(t, "https www g com search q test m 11 42 ie utf 8", s1)
}
