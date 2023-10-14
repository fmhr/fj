package fj

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
)

//go:embed cmd/compiler/Dockerfile
var compilerDockerfile string

//go:embed cmd/compiler/gcloudbuild.sh
var compilerGcloudbuild string

//go:embed cmd/compiler/localbuild.sh
var compilerLocalBuild string

func mkDirCompilerBase() {
	targetDir := "./fj/compiler/"

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	embeddedFiles := map[string]string{
		"Dockerfile":     compilerDockerfile,
		"gcloudbuild.sh": compilerGcloudbuild,
		"localbuild.sh":  compilerLocalBuild,
	}
	for file, content := range embeddedFiles {
		target := filepath.Join(targetDir, file)
		if err := os.WriteFile(target, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to write file: %s: %v", file, err)
		}
	}
}

//go:embed cmd/worker/Dockerfile
var workersDockerfile string

//go:embed cmd/worker/gcloudbuild.sh
var workersGcloudbuild string

//go:embed cmd/worker/localbuild.sh
var workersLocal string

func mkDirWorkerBase() {
	targetDir := "./fj/worker/"

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	embeddedFiles := map[string]string{
		"Dockerfile":     workersDockerfile,
		"gcloudbuild.sh": workersGcloudbuild,
		"localbuild.sh":  workersLocal,
	}
	for file, content := range embeddedFiles {
		target := filepath.Join(targetDir, file)
		if err := os.WriteFile(target, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to write file: %s: %v", file, err)
		}
	}
}
