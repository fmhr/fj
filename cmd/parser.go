package main

import (
	"regexp"
	"strconv"
	"strings"
)

type regexStr struct {
	re  *regexp.Regexp
	str string // 数値以外を削除するための文字列
}

var regexStrs = []regexStr{
	{regexp.MustCompile(`N=([0-9]+)`), "N="},
	{regexp.MustCompile(`L=([0-9]+)`), "L="},
	{regexp.MustCompile(`S=([0-9]+)`), "S="},
	{regexp.MustCompile(`Score = ([0-9]+)`), "Score = "},
	{regexp.MustCompile(`Number of wrong answers = ([0-9]+)`), "Number of wrong answers = "},
	{regexp.MustCompile(`Placement cost = ([0-9]+)`), "Placement cost = "},
	{regexp.MustCompile(`Measurement cost = ([0-9]+)`), "Measurement cost = "},
	{regexp.MustCompile(`Measurement count = ([0-9]+)`), "Measurement count = "},
}

func parseInt(src string, re *regexp.Regexp, str string) int {
	match := re.FindString(src)
	if match == "" {
		panic("no match, " + str + " not found")
	} else {
		rtn, err := strconv.Atoi(strings.Replace(match, str, "", -1))
		if err != nil {
			panic(err)
		}
		return rtn
	}
}
