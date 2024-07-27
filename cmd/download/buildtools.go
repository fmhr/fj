package download

import (
	"fmt"
	"os/exec"
)

// buildTools is build AtCoder tools
// 共通: cargo run -r --bin gen seeds.txt
// Reactive: cargo build --release --bin tester
// no Reactive: cargo build --release --bin vis
func buildTools() {
	// make inputfiles
	// AHC035のREADME.htmlを参考
	// 他のコンテストでコマンドが違っていたら修正
	// オプションは考慮しない
	// rustup update はしない
	generate := "cd tools && cargo run -r --bin gen seeds.txt"
	runCmd(generate)
	if IsReactive() {
		buildTester := "cd tools && cargo build --release --bin tester"
		runCmd(buildTester)
	} else {
		buildVis := "cd tools && cargo build --release --bin vis"
		runCmd(buildVis)
	}
}

func runCmd(cmds string) {
	cmd := exec.Command("sh", "-c", cmds)
	output, err := cmd.CombinedOutput()
	fmt.Println("Cmd:", "sh -c", cmds)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
	}
}
