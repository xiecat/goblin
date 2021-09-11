package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"html"
	"html/template"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// StrContainAndinList 多字符串匹配
func StrContainAndinList(rawStr string, checkStrList []string) bool {
	for _, checkStr := range checkStrList {
		if !strings.Contains(rawStr, checkStr) {
			return false
		}
	}

	return true
}

// StrPrefixOrinList 多字符串匹配
func StrPrefixOrinList(rawStr string, checkStrList []string) bool {
	for _, checkStr := range checkStrList {
		if strings.HasPrefix(rawStr, checkStr) {
			return true
		}
	}

	return false
}

// StrContainOrInList string in strings
func StrContainOrInList(rawStr string, checkStrList []string) bool {
	for _, checkStr := range checkStrList {
		if strings.Contains(rawStr, checkStr) {
			return true
		}
	}

	return false
}

// StrEqualOrInList string in string list
func StrEqualOrInList(rawStr string, checkStrList []string) bool {
	for _, checkStr := range checkStrList {
		if rawStr == checkStr {
			return true
		}
	}

	return false
}

// RandChar random char
func RandChar(size int) string {
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(char[rand.Int63()%int64(len(char))])
	}
	return s.String()
}

// TrimAll trim string for char
func TrimAll(s, cutset string) string {
	for _, c := range cutset {
		s = strings.ReplaceAll(s, string(c), "")
	}
	return s
}

// ReverseString string reverse
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Toupper string to upper
func Toupper(raw string) (interface{}, error) {
	return strings.ToUpper(raw), nil
}

// Tolower  string to lower
func Tolower(raw string) (interface{}, error) {
	return strings.ToLower(raw), nil
}

// B64Encode  base64 encode
func B64Encode(raw []byte) (string, error) {
	sEnc := base64.StdEncoding.EncodeToString(raw)
	return sEnc, nil
}

// B64Decode base64 decode
func B64Decode(raw string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(raw)
}

// URLEncode url encode
func URLEncode(raw string) (interface{}, error) {
	return url.PathEscape(raw), nil
}

// URLDecode url decode
func URLDecode(raw string) (interface{}, error) {
	return url.PathUnescape(raw)
}

// HexEncode  hex encode
func HexEncode(raw []byte) (interface{}, error) {
	return hex.EncodeToString(raw), nil
}

// HexDecode  hex decode
func HexDecode(raw string) (interface{}, error) {
	hx, _ := hex.DecodeString(raw)
	return string(hx), nil
}

// HTMLEscape html escape
func HTMLEscape(raw string) (interface{}, error) {
	return html.EscapeString(raw), nil
}

// HTMLUnescape html unescape
func HTMLUnescape(raw string) (interface{}, error) {
	return html.UnescapeString(raw), nil
}

// Md5 hashing
func Md5(raw []byte) (interface{}, error) {
	hash := md5.Sum(raw)

	return hex.EncodeToString(hash[:]), nil
}

// Sha256 hashing
func Sha256(raw []byte) (interface{}, error) {
	h := sha256.New()
	_, err := h.Write(raw)

	if err != nil {
		return nil, err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Sha1 hash
func Sha1(raw []byte) (interface{}, error) {
	h := sha1.New()
	_, err := h.Write(raw)

	if err != nil {
		return nil, err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func EleInArray(e string, arr []string) bool {
	for _, v := range arr {
		if strings.EqualFold(e, v) {
			return true
		}
	}
	return false
}

func TempStr(str string, temp interface{}) (string, error) {
	tmpl, err := template.New("tmp").Parse(str)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, temp)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
