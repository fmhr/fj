package cmd

type LanguageSet struct {
	Language   string
	FileName   string
	CompileCmd string
	BinaryPath string
	ExeCmd     string
}

var LanguageSets = map[string]LanguageSet{
	"Go": {
		Language:   "Go",
		FileName:   "main.go",
		CompileCmd: "go build -o a.out main.go",
		BinaryPath: "a.out",
		ExeCmd:     "./a.out",
	},
	"rust": {
		Language:   "rust",
		FileName:   "src/main.rs",
		CompileCmd: "cargo build --release --quiet --offline",
		BinaryPath: "target/release/main",
		ExeCmd:     "./target/release/main",
	},
	"C++": {
		Language:   "C++",
		FileName:   "main.cpp",
		CompileCmd: "g++-12 -std=gnu++20 -O2 -DONLINE_JUDGE -DATCODER -Wall -Wextra -mtune=native -march=native -fconstexpr-depth=2147483647 -fconstexpr-loop-limit=2147483647 -fconstexpr-ops-limit=2147483647 -I/opt/ac-library -I/opt/boost/gcc/include -L/opt/boost/gcc/lib -o a.out Main.cpp -lgmpxx -lgmp -L/usr/include/eigen3",
		BinaryPath: "a.out",
		ExeCmd:     "./a.out",
	},
	"java": {
		Language:   "java",
		FileName:   "Main.java",
		CompileCmd: "javac -encoding UTF-8 Main.java",
		BinaryPath: "Main.class",
		ExeCmd:     "java -Xss1024M -DONLINE_JUDGE=true Main",
	},
	"C#": {
		Language:   "C#",
		FileName:   "Main.cs",
		CompileCmd: "sh -c export DOTNET_EnableWriteXorExecute=0 && dotnet publish -c Release -o publish -v q --nologo 1>&2",
		BinaryPath: "publish/Main",
		ExeCmd:     "./publish/Main",
	},
}

func languageList() (langs []string) {
	for _, lang := range LanguageSets {
		langs = append(langs, lang.Language)
	}
	return
}
