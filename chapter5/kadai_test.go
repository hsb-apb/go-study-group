package chapter5

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumTestValidation(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		if err := Validation(1, 1); err != nil {
			t.Fail()
		}
	})
	t.Run("異常系１", func(t *testing.T) {
		t.Parallel()
		err := Validation(0, 1)
		assert.EqualError(t, err, "ファイルパスを指定してください。")
	})
	t.Run("異常系２", func(t *testing.T) {
		t.Parallel()
		err := Validation(1, 0)
		assert.EqualError(t, err, "-f は1以上である必要があります。")
	})
}

func TestCut(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		stdin := bytes.NewBufferString("foo,hogehoge,aaaaa\nfoo2,aabbcc,bbbbb")
		stdout := new(bytes.Buffer)
		err := Cut(stdin, stdout, ",", 1)
		assert.NoError(t, err)
		expected := "hogehoge\naabbcc\n"
		assert.Equal(t, expected, string(stdout.Bytes()))
	})

	t.Run("異常系１", func(t *testing.T) {
		t.Parallel()
		stdin := bytes.NewBufferString("foo,hogehoge,aaaaa\nfoo2,aabbcc,bbbbb")
		stdout := new(bytes.Buffer)
		err := Cut(stdin, stdout, ",", 4)
		assert.EqualError(t, err, "-fの値に該当するデータがありません")
	})
}

func BenchmarkCut(b *testing.B) {
	b.ResetTimer()
	stdin := bytes.NewBufferString("foo,hogehoge,aaaaa\nfoo2,aabbcc,bbbbb")
	stdout := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		Cut(stdin, stdout, ",", 2)
	}
}
