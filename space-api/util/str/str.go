package str

import (
	"unicode"
)

func ExtractAroundKeyword(str, keyword string, padding int) (string, int) {
	// 将原始字符串转换为 rune 切片，这样能够按字符进行操作
	runes := []rune(str)
	keywordRunes := []rune(keyword)

	// 查找关键词在字符串中的位置
	index := -1
	for i := 0; i <= len(runes)-len(keywordRunes); i++ {
		// 比较当前位置的字符与子串是否匹配
		match := true
		for j := 0; j < len(keywordRunes); j++ {
			if runes[i+j] != keywordRunes[j] {
				match = false
				break
			}
		}
		if match {
			index = i
			break
		}
	}

	if index == -1 {
		return "", -1 // 如果找不到关键词，返回空字符串和 -1
	}

	// 计算前后各10个字符的起始和结束位置
	start := max(index-padding, 0)
	end := min(index+len(keywordRunes)+padding, len(runes))

	// 截取包含关键词的子串
	subRunes := runes[start:end]

	// 获取关键词在新子串中的位置
	keywordIndexInSubStr := -1
	for i := 0; i <= len(subRunes)-len(keywordRunes); i++ {
		match := true
		for j := 0; j < len(keywordRunes); j++ {
			if subRunes[i+j] != keywordRunes[j] {
				match = false
				break
			}
		}
		if match {
			keywordIndexInSubStr = i
			break
		}
	}

	// 将 rune 切片转换回字符串
	return string(subRunes), keywordIndexInSubStr
}

// FindKeywordPositionRune 查找关键词在字符串中的起始位置，返回基于 rune 的位置
func FindKeywordPositionRune(s, keyword string) int {
	// 将字符串和关键词都转换为 rune 切片
	runes := []rune(s)
	keywordRunes := []rune(keyword)

	// 遍历 rune 切片，查找关键词的位置
	for i := 0; i <= len(runes)-len(keywordRunes); i++ {
		match := true
		// 比较当前位置的字符与子串是否匹配
		for j := 0; j < len(keywordRunes); j++ {
			if runes[i+j] != keywordRunes[j] {
				match = false
				break
			}
		}
		if match {
			return i // 返回关键词在 rune 切片中的位置
		}
	}
	return -1 // 如果找不到关键词，返回 -1
}

// 判断字符串是否不是标点符号
func IsNotPunctuation(s string) bool {
	// 如果字符串长度不为1，直接返回 false
	if len(s) != 1 {
		return false
	}

	// 将字符串转换为 rune 类型，判断是否不是标点符号
	r := rune(s[0])
	return !unicode.IsPunct(r)
}
