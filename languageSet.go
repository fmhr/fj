package fj

type LanguageSet struct {
	Language   string
	FileName   string
	CompileCmd string
	ExeCmd     string
}

var LanguageSets = []LanguageSet{
	{
		Language:   "go",
		FileName:   "main.go",
		CompileCmd: "go build -o a.out main.go",
		ExeCmd:     "./a.out",
	},
	{
		Language:   "rust",
		FileName:   "main.rs",
		CompileCmd: "cargo build --release --quiet --offline",
		ExeCmd:     "./a.out",
	},
	{
		Language:   "cpp",
		FileName:   "main.cpp",
		CompileCmd: "g++-12 -std=gnu++20 -O2 -DONLINE_JUDGE -DATCODER -Wall -Wextra -mtune=native -march=native -fconstexpr-depth=2147483647 -fconstexpr-loop-limit=2147483647 -fconstexpr-ops-limit=2147483647 -I/opt/ac-library -I/opt/boost/gcc/include -L/opt/boost/gcc/lib -o a.out Main.cpp -lgmpxx -lgmp -L/usr/include/eigen3",
		ExeCmd:     "./a.out",
	},
	{
		Language:   "java",
		FileName:   "Main.java",
		CompileCmd: "javac -Xss128M -DONLINE_JUDGE=true Main",
		ExeCmd:     "java Main",
	},
	{
		Language:   "cs",
		FileName:   "publish/Main",
		CompileCmd: "export DOTNET_EnableWriteXorExecute=0 && dotnet publish -c Release -o publish -v q --nologo 1>&2",
		ExeCmd:     "./publish/Main",
	},
}
