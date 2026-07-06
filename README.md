# fj

AtCoder Heuristic Contest (AHC) 向けのコマンドラインツール

## インストール

```bash
go install github.com/fmhr/fj@latest
```

## クイックスタート

```bash
# 1. 公式ローカルテスターをダウンロードして展開
fj download https://img.atcoder.jp/ahcXXX/tools.zip

# 2. 公式テスターをコンパイル (Rust製の場合の例)
cd tools
cargo build --release
cd ..

# 3. 設定ファイル用のディレクトリを作成
mkdir fj
# 後述の「設定ファイル」の内容で fj/config.toml を手動で作成してください

# 4. 自身のプログラムをコンパイル
g++ -O2 main.cpp -o a.out

# 5. テスト実行（実行するコマンドは必ず "" で囲んで指定します）
fj test "./a.out"
```

## コマンド

### `fj test` (エイリアス: `fj t`)

テストケースを実行します。実行コマンドは必ず `""` で囲んで指定してください。

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
| `--seed` | `-s` | シード値（デフォルト: 0） |
| `--count` | `-n` | テストケース数（デフォルト: 1） |
| `--parallel` | `-p` | 並列実行数（デフォルト: 1） |

### `fj download` (エイリアス: `fj d`)

テスターのZIPファイルをダウンロードして展開します。

```bash
fj download https://img.atcoder.jp/ahcXXX/tools.zip
```

### `fj info`

現在の設定情報を表示します。

```bash
fj info
```

### `fj version` (エイリアス: `fj v`)

バージョン情報を表示します。

## グローバルオプション

| オプション | 説明 |
|-----------|------|
| `--debug` | デバッグモード（ログにファイル名:行番号を表示） |
| `--best` | ベストスコアを更新 |
| `--stderr` | 実行プログラムのstderr出力を表示 |
| `--cloud` | クラウドモードを有効化（※Google Cloud環境の設定が必要です） |
| `--json` | JSON形式で出力 |
| `--csv` | CSV形式で result.csv に出力 |

## 設定ファイル

プロジェクトルートに `fj` ディレクトリを作成し、 `fj/config.toml` というファイル名で以下の設定例を作成してください。

### config.toml のデフォルト内容例
```toml
Language = 'Go'
ExecuteCmd = ''
Args = []
Contest = ''
Reactive = false
TimeLimitMS = 7000
TesterPath = 'tools/target/release/tester'
VisPath = 'tools/target/release/vis'
GenPath = 'tools/target/release/gen'
InfilePath = 'tools/in/'
OutfilePath = 'out/'
Jobs = 0
BinaryPath = ''
Cloud = false
CompilerURL = ''
Source = 'main.go'
tmpBinary = ''
Bucket = 'binary-buckets'
WorkerURL = ''
ConcurrentRequests = 10
```

| 設定項目 | 説明 |
|---------|------|
| `Language` | 使用言語（Go, rust, C++, java, C# 等） |
| `ExecuteCmd` | 実行コマンド |
| `Args` | 実行コマンドの追加オプション引数 |
| `Reactive` | リアクティブ問題かどうか |
| `TesterPath` | リアクティブ問題のテスターパス |
| `VisPath` | ビジュアライザ/スコア計算プログラムのパス |
| `GenPath` | テストケース生成プログラムのパス |
| `InfilePath` | 入力ファイルのディレクトリ |
| `OutfilePath` | 出力ファイルのディレクトリ |
| `TimeLimitMS` | タイムアウト時間（ミリ秒単位。制限時間を超過すると `TimeOut` と判定されます） |
| `Jobs` | 並列実行ジョブ数 |
| `Cloud` | クラウド実行モードをデフォルトで有効にするか |

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