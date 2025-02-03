package repo

import (
	"bearguard/cm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Member struct {
	Name     string `db:"name" json:"name"`
	MemberID string `db:"member_id" json:"member_id"`
}

func GetDBMembers() ([]Member, error) {
	// プロジェクトルートからJSONファイルのパスを取得
	jsonFilePath := filepath.Join(cm.GetProjectRoot(), "roomId.json")

	// JSONファイルを読み込む
	fileData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	// JSONデータをパース
	var data struct {
		RoomId []struct {
			ID        int    `json:"id"`
			OwnerName string `json:"ownerName"`
		} `json:"roomId"`
	}
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	// Member型に変換
	var members []Member
	for _, room := range data.RoomId {
		members = append(members, Member{
			Name:     room.OwnerName,
			MemberID: strconv.Itoa(room.ID),
		})
	}

	return members, nil
}
