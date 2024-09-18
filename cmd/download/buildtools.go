package download

import (
	"fmt"
	"log"
	"os/exec"
)

// buildTools is build AtCoder tools
// 共通: cargo run -r --bin gen seeds.txt
// Reactive: cargo build --release --bin tester
// no Reactive: cargo build --release --bin vis
func buildTools() error {
	// make inputfiles
	// AHC035のREADME.htmlを参考
	// 他のコンテストでコマンドが違っていたら修正
	// オプションは考慮しない
	// rustup update はしない
	generate := "cd tools && cargo run -r --bin gen seeds.txt"
	runCmd(generate)
	if IsReactive() {
		buildTester := "cd tools && cargo build --release --bin tester"
		return runCmd(buildTester)
	}
	buildVis := "cd tools && cargo build --release --bin vis"
	return runCmd(buildVis)
}

func runCmd(cmds string) error {
	cmd := exec.Command("sh", "-c", cmds)
	output, err := cmd.CombinedOutput()
	fmt.Println("Cmd:", "sh -c", cmds)
	if err != nil {
		log.Println("CMD:", "sh -c", cmds)
		log.Println("Output:", string(output))
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Println("[SUCCESS]", cmds)
	return nil
}
