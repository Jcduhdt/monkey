package parser

import (
	"fmt"
	"strings"
)

// 查看语句的parser调用情况
var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

// 通过修改这个函数，关闭开启trace打印
func tracePrint(fs string) {
	// 开启
	// fmt.Printf("%s%s\n", identLevel(), fs)
	// 关闭
	fmt.Sprintf("%s%s\n", identLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func unTrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
