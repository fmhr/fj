package fj

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	// 必要なフォルダの作成
	if _, err := os.Stat("testdata"); err != nil {
		if err := os.Mkdir("testdata", 0755); err != nil {
			t.Fatalf("failed to create testdata directory: %v", err)
		}
	}
	// テストの設定
	cnf := &Config{
		Cmd:         "cat",
		InfilePath:  "testdata/input",
		OutfilePath: "testdata/output",
	}
	seed := 1

	// 入力ファイルの作成
	inputConten := "Hello, World!\n"
	err := os.WriteFile("testdata/input0001.txt", []byte(inputConten), 0644)
	if err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	// テストの実行
	_, err = Run(cnf, seed)
	if err != nil {
		t.Fatalf("failed to run: %v", err)
	}

	// 出力ファイルの確認
	outputContent, err := os.ReadFile("testdata/output0001.out")
	if err != nil || string(outputContent) != inputConten {
		t.Fatalf("Unexpected output. Expected: %s, Actual: %s", inputConten, string(outputContent))
	}

	// 後始末
	defer os.Remove("testdata/input0001.txt")
	defer os.Remove("testdata/output0001.out")
}
