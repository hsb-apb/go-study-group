package chapter5

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// kadai_baseのほうで定義されているからコッチではコメントアウトしておく
// var delimiter = flag.String("d", ",", "区切り文字を指定してください")
// var fields = flag.Int("f", 1, "フィールドの何番目を取り出すか指定してください")

func Validation(argCount int, fieldNum int) error {

	if argCount == 0 {
		return fmt.Errorf("ファイルパスを指定してください。")
	}
	if fieldNum <= 0 {
		return fmt.Errorf("-f は1以上である必要があります。")
	}
	return nil
}

func Cut(r io.Reader, w io.Writer, delimiterStr string, fieldNum int) error {
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)
	for scanner.Scan() {
		text := scanner.Text()
		sb := strings.Split(text, delimiterStr)
		if len(sb) <= fieldNum-1 {
			return fmt.Errorf("-fの値に該当するデータがありません")
		}
		s := sb[fieldNum-1]
		fmt.Fprintln(writer, s)
	}
	writer.Flush()
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// go-cutコマンドを実装しよう
func GoCut() {
	flag.Parse()

	if err := Validation(flag.NArg(), *fields); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := Cut(file, os.Stdout, *delimiter, *fields); err != nil {
		log.Fatal(err)
	}
}
