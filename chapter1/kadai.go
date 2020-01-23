package chapter1

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/apbgo/go-study-group/chapter1/lib"
)

// Calc opには+,-,×,÷の4つが渡ってくることを想定してxとyについて計算して返却(正常時はerrorはnilでよい)
// 想定していないopが渡って来た時には0とerrorを返却
func Calc(op string, x, y int) (int, error) {

	// ヒント：エラーにも色々な生成方法があるが、ここではシンプルにfmtパッケージの
	// fmt.Errorf(“invalid op=%s”, op) などでエラー内容を返却するのがよい
	// https://golang.org/pkg/fmt/#Errorf

	// TODO Q1
	result := 0
	var err error
	switch op {
	case "+":
		result = x + y
	case "-":
		result = x - y
	case "×":
		result = x * y
	case "÷":
		if y == 0 {
			err = fmt.Errorf("integer divide by zero")
		} else {
			result = x / y
		}
	default:
		err = fmt.Errorf("invalid op=%s", op)
	}
	return result, err
}

// StringEncode 引数strの長さが5以下の時キャメルケースにして返却、それ以外であればスネークケースにして返却
func StringEncode(str string) string {
	// ヒント：長さ(バイト長)はlen(str)で取得できる
	// chapter1/libのToCamelとToSnakeを使うこと

	// TODO Q2
	if len(str) <= 5 {
		return lib.ToCamel(str)
	}
	return lib.ToSnake(str)
}

// Sqrt 数値xが与えられたときにz²が最もxに近い数値zを返却
func Sqrt(x float64) float64 {

	// TODO Q3
	z := 1.0
	for !isApproximation(z*z, x) {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

// 差が0.00000000000005以下ならtrueを返却する
func isApproximation(x, y float64) bool {
	difference := x - y
	if difference < 0 {
		difference *= -1
	}
	if difference < 0.00000000000005 {
		// 差が既定値以下
		return true
	}
	return false
}

// Pyramid x段のピラミッドを文字列にして返却
// 期待する戻り値の例：x=5のとき "1\n12\n123\n1234\n12345"
// （x<=0の時は"error"を返却）
func Pyramid(x int) string {
	// ヒント：string <-> intにはstrconvを使う
	// int -> stringはstrconv.Ioa() https://golang.org/pkg/strconv/#Itoa

	// TODO Q4
	if x <= 0 {
		// 0以下はerrorを返却
		return "error"
	}
	var result string = ""
	for i := 1; i <= x; i++ {
		if i != 1 {
			// 2以上のときに改行つける
			result += "\n"
		}
		for j := 1; j <= i; j++ {
			result += strconv.Itoa(j)
		}
	}

	return result
}

// StringSum x,yをintにキャストし合計値を返却 (正常終了時、errorはnilでよい)
// キャスト時にエラーがあれば0とエラーを返却
func StringSum(x, y string) (int, error) {

	// ヒント：string <-> intにはstrconvを使う
	// string -> intはstrconv.Atoi() https://golang.org/pkg/strconv/#Atoi

	// TODO Q5
	iX, err := strconv.Atoi(x)
	if err != nil {
		return 0, err
	}
	iY, err := strconv.Atoi(y)
	if err != nil {
		return 0, err
	}
	return iX + iY, nil
}

// SumFromFileNumber ファイルを開いてそこに記載のある数字の和を返却
func SumFromFileNumber(filePath string) (int, error) {
	// ヒント：ファイルの扱い：os.Open()/os.Close()
	// bufio.Scannerなどで１行ずつ読み込むと良い

	// TODO Q6 オプション
	file, err := os.Open(filePath)
	if err != nil {
		// 開けない
		return 0, err
	}
	defer file.Close()
	result := 0
	sc := bufio.NewScanner(file)
	for i := 1; sc.Scan(); i++ {
		if err := sc.Err(); err != nil {
			// エラー処理
			return 0, err
		}
		rowVal, err := strconv.Atoi(sc.Text())
		if err != nil {
			// 数字じゃなかったら飛ばす
			continue
		}
		result += rowVal
	}
	return result, nil
}
