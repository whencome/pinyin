package pinyin

import (
	"unicode/utf8"
)

// 获取汉字对应的拼音
// 如果传入的值不是汉字，则原样返回
func getPinyin(hanzi rune, mode Mode) (string, bool) {
	// 根据模式处理
	switch mode {
	case ModeTone:
		return getTone(hanzi)
	case ModeInitialsInCapitals:
        return getInitialsInCapitals(hanzi)
	default:
		return getDefault(hanzi)
	}
}

// 获取包含声调的拼音
// string : 获取的结果
// bool : 标识结果是否是拼音
func getTone(hanzi rune) (string, bool) {
	py, ok := pinyinMap[hanzi]
	if !ok {
		return string(hanzi), false
	}
	return py, true
}

// 获取不包含声调的拼音
func getDefault(hanzi rune) (string, bool) {
	tone, isPinYin := getTone(hanzi)
	if !isPinYin || tone == ""{
		return tone, isPinYin
	}

	output := make([]rune, utf8.RuneCountInString(tone))
	count := 0
	for _, t := range tone {
		neutral, found := tonesMap[t]
		if found {
			output[count] = neutral
		} else {
			output[count] = t
		}
		count++
	}
	return string(output), isPinYin
}

// 获取不包含声调的拼音，且首字母大写
func getInitialsInCapitals(hanzi rune) (string, bool) {
	def, isPinYin := getDefault(hanzi)
	if !isPinYin || def == ""{
		return def, isPinYin
	}
	sr := []rune(def)
	if sr[0] > 32 {
		sr[0] = sr[0] - 32
	}
	return string(sr), isPinYin
}

