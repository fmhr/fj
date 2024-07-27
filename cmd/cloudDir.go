package fj

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed"../compiler/script/*"
var compilerFiles embed.FS

//go:embed"../worker/script/*""
var workerFiles embed.FS

func mkDirCompilerBase() {
	targetDir := "./fj/compiler/"
	extractEmbeddedFiles(compilerFiles, "cmd/compiler/script", targetDir)
}

func mkDirWorkerBase() {
	targetDir := "./fj/worker/"
	extractEmbeddedFiles(workerFiles, "cmd/worker/script", targetDir)
}

func extractEmbeddedFiles(embeddedFS embed.FS, sourcdDir, targetDir string) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}
	// embeddedFSからファイルを取得
	fs.WalkDir(embeddedFS, sourcdDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// ディレクトリの場合はスキップ
		if d.IsDir() {
			return nil
		}
		// ファイルの内容を取得
		content, err := fs.ReadFile(embeddedFS, path)
		if err != nil {
			return err
		}
		// ターゲットパスを作成
		relativePath, err := filepath.Rel(sourcdDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDir, relativePath)

		// ファイルを作成
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			log.Fatalf("failed to create file: %v", err)
		}
		return nil
	})

}
