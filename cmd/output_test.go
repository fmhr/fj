package cmd

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJsonOutput(t *testing.T) {
	// 一時ディレクトリを作成して移動し、テスト後に元に戻る
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory to tempDir: %v", err)
	}
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

	// テストデータ作成
	sm1 := NewSliceMap()
	sm1.Set("seed", "0")
	sm1.Set("score", "1000")

	sm2 := NewSliceMap()
	sm2.Set("seed", "1")
	sm2.Set("score", "2000")

	datas := []SliceMap{sm1, sm2}

	// 実行
	err = writeJsonOutput(datas)
	if err != nil {
		t.Fatalf("JsonOutput failed: %v", err)
	}

	// ファイル出力の検証
	files, err := os.ReadDir(filepath.Join("fj", "data"))
	if err != nil {
		t.Fatalf("failed to read output directory: %v", err)
	}

	var jsonFile os.DirEntry
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "result_") && strings.HasSuffix(file.Name(), ".json") {
			jsonFile = file
			break
		}
	}

	if jsonFile == nil {
		t.Fatal("expected output JSON file not found")
	}

	// 内容の検証
	content, err := os.ReadFile(filepath.Join("fj", "data", jsonFile.Name()))
	if err != nil {
		t.Fatalf("failed to read generated JSON file: %v", err)
	}

	var decoded []SliceMap
	err = json.Unmarshal(content, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal generated JSON: %v", err)
	}

	if len(decoded) != 2 {
		t.Fatalf("expected 2 items, got %d", len(decoded))
	}

	s0, _ := decoded[0].Get("seed")
	s1, _ := decoded[1].Get("seed")
	if s0 != "0" || s1 != "1" {
		t.Errorf("expected seeds to be 0 and 1, got %s and %s", s0, s1)
	}
}

func TestCsvOutput(t *testing.T) {
	tempDir := t.TempDir()
	csvFilePath := filepath.Join(tempDir, "test_result")

	// テストデータ作成
	sm1 := NewSliceMap()
	sm1.Set("seed", "0")
	sm1.Set("score", "1000")
	sm1.Set("extra", "foo")

	sm2 := NewSliceMap()
	sm2.Set("seed", "1")
	sm2.Set("score", "2000")

	datas := []SliceMap{sm1, sm2}

	// 実行 (内部で .csv 拡張子が追加されるため、csvFilePath.csv が作成される)
	err := writeCsvOutput(datas, csvFilePath)
	if err != nil {
		t.Fatalf("CsvOutput failed: %v", err)
	}

	generatedFile := csvFilePath + ".csv"
	file, err := os.Open(generatedFile)
	if err != nil {
		t.Fatalf("failed to open generated CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV records: %v", err)
	}

	// レコード数の検証
	// ヘッダー + 2つのデータ行 = 3行
	if len(records) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(records))
	}

	// ヘッダーの検証
	headers := records[0]
	if len(headers) != 3 {
		t.Fatalf("expected 3 headers, got %d", len(headers))
	}

	// 各行の値の検証
	// ヘッダーが ["seed", "score", "extra"] の並びであることを確認
	keyIndices := make(map[string]int)
	for i, h := range headers {
		keyIndices[h] = i
	}

	// 1行目のデータ確認
	row1 := records[1]
	if row1[keyIndices["seed"]] != "0" || row1[keyIndices["score"]] != "1000" || row1[keyIndices["extra"]] != "foo" {
		t.Errorf("row1 unexpected data: %v", row1)
	}

	// 2行目のデータ確認 (extra は空のはず)
	row2 := records[2]
	if row2[keyIndices["seed"]] != "1" || row2[keyIndices["score"]] != "2000" || row2[keyIndices["extra"]] != "" {
		t.Errorf("row2 unexpected data: %v", row2)
	}
}
