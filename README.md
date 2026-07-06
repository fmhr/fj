# fj

AtCoder Heuristic Contest (AHC) 向けのコマンドラインツール

## インストール
- Go言語がインストールされている環境でインストールが可能です
```bash
go install github.com/fmhr/fj@latest
```

## クイックスタート
- 公式ローカル版テスターが同一ディレクトリ内にあることが前提

# 4. 自身のプログラムをコンパイル
g++ -O2 main.cpp -o a.out

# 5. テスト実行（実行するコマンドは必ず "" で囲んで指定します）
```
fj test "./a.out"
```

## コマンド

```bash
# 単一のテストケース実行（seed=0）
fj test "./a.out"

# シード指定
fj test "./a.out" -s 5

# 複数ケース実行
fj test "./a.out" -n 100

# 並列実行
fj test "./a.out" -n 100 -p 4

# ベストスコアを更新
fj test "./a.out" -n 100 --best
```

| オプション | 短縮形 | 説明 |
|-----------|-------|------|
| `--seed` | `-s` | シード値を指定して実行（デフォルト: 0） |
| `--count` | `-n` | テストケース数を指定して実行（デフォルト: 1） |
| `--parallel` | `-p` | 並列実行数を指定して実行（デフォルト: 1） |


| オプション | 説明 |
|-----------|------|
| `--debug` | デバッグモード（ログにファイル名:行番号を表示） |
| `--best` | ベストスコアを更新 |
| `--stderr` | 実行プログラムのstderr出力を表示 |
| `--cloud` | クラウドモードを有効化（※Google Cloud環境の設定が必要です） |
| `--json` | JSON形式で出力 |
| `--csv` | CSV形式で result.csv に出力 |

## ディレクトリ構成

```
project/
├── fj/
│   ├── config.toml      # 設定ファイル
│   └── best_score.json  # ベストスコア記録（--best使用時）
├── tools/
│   ├── in/              # 入力ファイル
│   └── target/release/
│       ├── tester       # リアクティブ用テスター
│       └── vis          # ビジュアライザ
├── out/                 # 出力ファイル
└── a.out                # 実行ファイル
```

## ライセンス

MIT License