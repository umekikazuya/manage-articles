package src

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// ItemToCSV はQiitaのアイテム情報をCSV形式に変換します
func ItemToCSV(items []Item, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("CSVファイル作成エラー: %w", err)
	}
	defer file.Close()

	// UTF-8 BOMを書き込む（Excelでの文字化け対策）
	_, err = file.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		return fmt.Errorf("BOM書き込みエラー: %w", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダー行を書き込む
	headers := []string{
		"ID", "タイトル", "URL", "投稿日", "更新日", "いいね数", "タグ",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ヘッダー書き込みエラー: %w", err)
	}

	// 日付フォーマット
	dateFormat := "2006/01/02 15:04:05"

	// 各アイテムをCSVに書き込む
	for _, item := range items {
		// タグをカンマ区切りの文字列に変換
		var tags []string
		for _, tag := range item.Tags {
			tags = append(tags, tag.Name)
		}
		tagString := strings.Join(tags, ", ")

		record := []string{
			item.ID,
			item.Title,
			item.URL,
			item.CreatedAt.Format(dateFormat),
			item.UpdatedAt.Format(dateFormat),
			fmt.Sprintf("%d", item.LikesCount),
			tagString,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("レコード書き込みエラー: %w", err)
		}
	}

	return nil
}
