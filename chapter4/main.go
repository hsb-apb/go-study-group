package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var delimiter = flag.String("d", ",", "区切り文字を指定してください")
var fields = flag.Int("f", 1, "フィールドの何番目を取り出すか指定してください")

// go-cutコマンドを実装しよう
func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "ファイルパスを指定してください。")
		os.Exit(1)
	}
	fp, err := os.Open(flag.Args()[0])
	if err != nil {
		// Openエラー処理
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	writer := bufio.NewWriter(os.Stdout)
	for scanner.Scan() {
		slice := strings.Split(scanner.Text(), *delimiter)
		for i, str := range slice {
			if i == (*fields - 1) {
				str = strings.Trim(str, " ")
				writer.WriteString(str)
				writer.WriteRune('\n')
				break
			}
		}
	}
	writer.Flush()
}
