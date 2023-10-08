// 言語選択をしてそれにあわせたDockerfileからCloud buildsをつかう。

package fj

import (
	_ "embed"
	"log"
)

func builds(srcfile string, language string) {
	if language != "Go" {
		log.Fatalf("%s is no suport language", language)
	}
	// dockerfile

}
