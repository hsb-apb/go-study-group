package chapter2

import "fmt"

// 引数のスライスsliceの要素数が
// 0の場合、0とエラー
// 2以下の場合、要素を掛け算
// 3以上の場合、要素を足し算
// を返却。正常終了時、errorはnilでよい
func Calc(slice []int) (int, error) {
	// TODO Q1
	// ヒント：エラーにも色々な生成方法があるが、ここではシンプルにfmtパッケージの
	// fmt.Errorf(“invalid op=%s”, op) などでエラー内容を返却するのがよい
	// https://golang.org/pkg/fmt/#Errorf

	length := len(slice)
	if length == 0 {
		return 0, fmt.Errorf("slice length is zero")
	}
	if length == 1 {
		return slice[0], nil
	}
	if length == 2 {
		return slice[0] * slice[1], nil
	}
	if length >= 3 {
		var ret int
		for _, v := range slice {
			ret += v
		}
		return ret, nil
	}
	return 0, nil
}

type Number struct {
	index int
}

// 構造体Numberを3つの要素数から成るスライスにして返却
// 3つの要素の中身は[{1} {2} {3}]とし、append関数を使用すること
func Numbers() []Number {
	// TODO Q2
	ret := make([]Number, 0, 3)
	ret = append(ret, Number{index: 1})
	ret = append(ret, Number{index: 2})
	ret = append(ret, Number{index: 3})
	return ret
}

// 引数mをforで回し、「値」部分だけの和を返却
// キーに「yon」が含まれる場合は、キー「yon」に関連する値は除外すること
func CalcMap(m map[string]int) int {
	// TODO Q3
	var ret int
	for key, value := range m {
		if key != "yon" {
			ret += value
		}
	}
	return ret
}

type Model struct {
	Value int
}

// 与えられたスライスのModel全てのValueに5を足す破壊的な関数を作成
func Add(models []Model) {
	// TODO  Q4
	for i, _ := range models {
		models[i].Value += 5
	}
}

// 引数のスライスには重複な値が格納されているのでユニークな値のスライスに加工して返却
// 順序はスライスに格納されている順番のまま返却すること
// ex) 引数:[]slice{21,21,4,5} 戻り値:[]int{21,4,5}
func Unique(slice []int) []int {
	// TODO Q5
	ret := make([]int, 0)
	uniqueMap := make(map[int]struct{})
	for _, v := range slice {
		_, ok := uniqueMap[v]
		if ok {
			// uniqueMapに含まれている
			continue
		}
		ret = append(ret, v)
		uniqueMap[v] = struct{}{}
	}
	return ret
}

// 連続するフィボナッチ数(0, 1, 1, 2, 3, 5, ...)を返す関数(クロージャ)を返却
func Fibonacci() func() int {
	// TODO Q6 オプション
	list := make([]int, 0)
	return func() int {
		listLen := len(list)
		if listLen == 0 {
			list = append(list, 0)
			return 0
		}
		if listLen == 1 {
			list = append(list, 1)
			return 1
		}
		ret := list[0] + list[1]
		list = append(list, ret)
		// 直近２つだけ見れたらいい → appendしたら前を詰める
		list = list[1:]
		return ret
	}
}
