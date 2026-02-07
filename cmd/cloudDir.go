package cmd

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed compiler/script/*
var compilerFiles embed.FS

//go:embed worker/script/*
var workerFiles embed.FS

func mkDirCompilerBase() error {
    targetDir := "./fj/compiler/"
    return extractEmbeddedFiles(compilerFiles, "compiler/script", targetDir)
}

func mkDirWorkerBase() error {
    targetDir := "./fj/worker/"
    return extractEmbeddedFiles(workerFiles, "worker/script", targetDir)
}

func extractEmbeddedFiles(embeddedFS embed.FS, sourceDir, targetDir string) error {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			log.Println("Failed to create directory")
			return err
		}
	}
	// embeddedFSからファイルを取得
	err := fs.WalkDir(embeddedFS, sourceDir, func(path string, d fs.DirEntry, err error) error {
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
		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDir, relativePath)

		// ファイルを作成
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			log.Println("Failed to create file")
			return err
		}
		return nil
	})
	return err
}
