package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/umekikazuya/qiita-to-csv/src"
)

func main() {
	// コマンドライン引数の処理
	tokenPtr := flag.String("token", "", "Qiita API access token")
	outputPtr := flag.String("output", "", "Output CSV file path")
	flag.Parse()

	token := *tokenPtr

	// 環境変数からトークンを取得（コマンドライン引数がない場合）
	if token == "" {
		token = os.Getenv("QIITA_ACCESS_TOKEN")
		if token == "" {
			fmt.Println("エラー: Qiitaアクセストークンが必要です。-token オプションまたは QIITA_ACCESS_TOKEN 環境変数で指定してください。")
			os.Exit(1)
		}
	}

	// 出力ファイル名が指定されていない場合は自動生成
	output := *outputPtr
	if output == "" {
		output = fmt.Sprintf("qiita_articles_%s.csv", time.Now().Format("20060102_150405"))
	}

	// Qiitaクライアントの作成
	client := src.NewQiitaClient(token)

	// 全ての投稿を取得
	fmt.Println("Qiitaから投稿データを取得中...")
	items, err := client.GetAllUserItems()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d件の投稿を取得しました\n", len(items))

	// CSVに保存
	fmt.Printf("CSVファイルに保存中: %s\n", output)
	if err := src.ItemToCSV(items, output); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("完了しました！")
}
