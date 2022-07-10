package utils

import (
	"fmt"
	"strings"
)

func AutoUnitConvert(value int64) string {
	if value < 1024 {
		return fmt.Sprintf("%dB", value)
	} else if value < 1024*1024 {
		return fmt.Sprintf("%.2fK", float64(value)/1024)
	} else if value < 1024*1024*1024 {
		return fmt.Sprintf("%.2fM", float64(value)/1024/1024)
	} else if value < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fG", float64(value)/1024/1024/1024)
	} else if value < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fT", float64(value)/1024/1024/1024/1024)
	} else {
		return fmt.Sprintf("%.2fP", float64(value)/1024/1024/1024/1024/1024)
	}
}

func AutoBandwidthConvert(value int64) string {
	value *= 8
	if value < 1024 {
		return fmt.Sprintf("%dbps", value)
	} else if value < 1024*1024 {
		return fmt.Sprintf("%.2fKbps", float64(value)/1024)
	} else if value < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMbps", float64(value)/1024/1024)
	} else if value < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fGbps", float64(value)/1024/1024/1024)
	} else {
		return fmt.Sprintf("%.2fTbps", float64(value)/1024/1024/1024/1024)
	}
}

func ParseForMarkdown(str string) string {
	// TODO: 后续需要优化
	strList := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, v := range strList {
		str = strings.Replace(str, v, "\\"+v, -1)
	}
	return str
}
