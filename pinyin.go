/**
 * 中文转拼音的代码由 https://github.com/chain-zhang/pinyin 修改而来，增加了如下改动：
 * 1. 拼音字典内容直接放到代码中，提高读取效率，原代码从网络读取不确定性较高；
 * 2. 增加保持非中文等其他符号的功能，如果不能转换成拼音则原样保留，比如字母、符号、数字等；
 */
package pinyin

import (
	"bufio"
	"strconv"
	"strings"
	"bytes"
)

var (
	tones = [][]rune{
		{'ā', 'ē', 'ī', 'ō', 'ū', 'ǖ', 'Ā', 'Ē', 'Ī', 'Ō', 'Ū', 'Ǖ'},
		{'á', 'é', 'í', 'ó', 'ú', 'ǘ', 'Á', 'É', 'Í', 'Ó', 'Ú', 'Ǘ'},
		{'ǎ', 'ě', 'ǐ', 'ǒ', 'ǔ', 'ǚ', 'Ǎ', 'Ě', 'Ǐ', 'Ǒ', 'Ǔ', 'Ǚ'},
		{'à', 'è', 'ì', 'ò', 'ù', 'ǜ', 'À', 'È', 'Ì', 'Ò', 'Ù', 'Ǜ'},
	}
	neutrals = []rune{'a', 'e', 'i', 'o', 'u', 'v', 'A', 'E', 'I', 'O', 'U', 'V'}
)

var (
	// 从带声调的声母到对应的英文字符的映射
	tonesMap map[rune]rune
	// 从汉字到声调的映射
	numericTonesMap map[rune]int
	// 从汉字到拼音的映射（带声调）
	pinyinMap map[rune]string
	// 初始状态
	initialized bool
)

type Mode int

const (
	ModeWithoutTone        Mode = iota + 1 // 默认模式，例如：guo
	ModeTone                               // 带声调的拼音 例如：guó
	ModeInitialsInCapitals                 // 首字母大写不带声调，例如：Guo
)

// 初始化
func init() {
	tonesMap = make(map[rune]rune)
	numericTonesMap = make(map[rune]int)
	pinyinMap = make(map[rune]string)
	for i, runes := range tones {
		for j, tone := range runes {
			tonesMap[tone] = neutrals[j]
			numericTonesMap[tone] = i + 1
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(PinYinCodeMaps))
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), "=>")
		if len(strs) < 2 {
			continue
		}
		i, err := strconv.ParseInt(strs[0], 16, 32)
		if err != nil {
			continue
		}
		pinyinMap[rune(i)] = strs[1]
	}
	initialized = true
}

// 拼音结构体
type pinyin struct {
	origin 		string
	separator  	string
	mode   		Mode
}

func New(origin string) *pinyin {
	return &pinyin{
		origin		: origin,
		separator	: " ",
		mode		: ModeWithoutTone,
	}
}

// 设置分隔符
func (py *pinyin) Separate(separator string) *pinyin  {
	py.separator = separator
	return py
}

// 设置转换模式
func (py *pinyin) Mode(mode Mode) *pinyin  {
	py.mode = mode
	return py
}

// 将中文转换为拼音
func (py *pinyin) Convert() (string, error) {
	if !initialized {
		return "", ErrNotInitialized
	}
	sr := []rune(py.origin)
	words := bytes.Buffer{}
	// 空格
	whiteSpace := " "
	whiteSpaceRune := rune(' ')
	// 标记前一个值是否是拼音
	isPrevPinYin := false
	// 标记前一个是否是字母或者数字
	isPrevAlphaDigit := false
	for _, s := range sr{
		isAlphaDigit := (s >= 'a' && s <= 'z') || (s >= 'A' && s <= 'Z') || (s >= '0' && s <= '9')
		// 如果s是单个字符（因为中文不可能是单个字符），则不需要处理
		// 这里需要注意，为了区别字母和拼音，所以如果是单个字母，且前面一个是拼音，则自动插入一个空格
		// 如果是数字或者符合，则原样保留
		if isAlphaDigit && isPrevPinYin && s != whiteSpaceRune {
			words.WriteString(" ")
			words.WriteRune(s)
			isPrevPinYin = false
			isPrevAlphaDigit = isAlphaDigit
			continue
		}
		word, isPinYin := getPinyin(s, py.mode)
		if len(word) > 0 && string(s) != word {
			if isPrevAlphaDigit && isPinYin && word != whiteSpace {
				words.WriteString(" ")
			}
			if isPrevPinYin && isPinYin && word != whiteSpace {
				words.WriteString(py.separator)
			}
			words.WriteString(word)
			isPrevPinYin = isPinYin
		} else {
			words.WriteString(word)
			isPrevPinYin = false
		}
		isPrevAlphaDigit = isAlphaDigit
	}
	return words.String(), nil
}

// 使用默认参数，将中文转换成拼音
func Convert(cnStr string) (string, error) {
	py := New(cnStr)
	return py.Convert()
}