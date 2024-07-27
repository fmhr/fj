module main

go 1.22.4

require github.com/fmhr/fj v0.0.0-20240725183114-19e7acdb4ec0

require (
	github.com/alecthomas/kingpin/v2 v2.4.0 // indirect
	github.com/alecthomas/units v0.0.0-20231202071711-9a357b53e9c9 // indirect
	github.com/elliotchance/orderedmap/v2 v2.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

// TODO go/srcで開発することでここは消す
replace github.com/fmhr/fj => ./cmd
