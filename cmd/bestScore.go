package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

type ScoreEntry struct {
	Seed      int `json:"seed"`
	BestScore int `json:"best_score"`
}

type BestData struct {
	MiniMax int          `json:"minimax"`
	Scores  []ScoreEntry `json:"scores"`
}

const bestScoreFile = "fj/best_score.json"

var mutex sync.Mutex

// createNewBestScorejson は新しいベストスコアjsonを作成する
func createNewBestScorejson(minimax int) error {
	// if already exists, return
	if _, err := os.Stat(bestScoreFile); err == nil {
		log.Printf("%s already exists\n", bestScoreFile)
		return nil
	}

	if minimax != 1 && minimax != -1 {
		return fmt.Errorf("minimax must be 1 or -1")
	}

	jsonData := BestData{
		MiniMax: minimax,
		Scores:  []ScoreEntry{},
	}

	data, err := json.MarshalIndent(jsonData, "", "")
	if err != nil {
		return fmt.Errorf("failed to marshal best score data: %v", err)
	}

	if err := os.WriteFile(bestScoreFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write best score file: %v", err)
	}
	return nil
}

func readBestScore() (*BestData, error) {
	file, err := os.ReadFile(bestScoreFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read best score file: %v", err)
	}

	var data BestData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal best score data: %v", err)
	}

	return &data, nil
}

func writeBestScore(data *BestData) error {
	sort.Slice(data.Scores, func(i, j int) bool {
		return data.Scores[i].Seed < data.Scores[j].Seed
	})
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal best score data: %v", err)
	}

	if err := os.WriteFile(bestScoreFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write best score file: %v", err)
	}
	return nil
}

func GetBestScores() (map[int]int, error) {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := readBestScore()
	if err != nil {
		return nil, err
	}

	bestScores := make(map[int]int)
	for _, score := range data.Scores {
		bestScores[score.Seed] = score.BestScore
	}
	return bestScores, nil
}

// UpdateBestScore は指定されたseedのベストスコアを設定する
func UpdateBestScore(seed, score int) error {
	if score <= 0 {
		return nil
	}
	mutex.Lock()
	defer mutex.Unlock()

	data, err := readBestScore()
	if err != nil {
		return err
	}
	for i, s := range data.Scores {
		if s.Seed == seed {
			if data.MiniMax == 1 {
				if score > data.Scores[i].BestScore || data.Scores[i].BestScore == -1 {
					data.Scores[i].BestScore = score
				}
			} else if data.MiniMax == -1 {
				if score < data.Scores[i].BestScore || data.Scores[i].BestScore == -1 {
					data.Scores[i].BestScore = score
				}
			}
			return writeBestScore(data)
		}
	}
	data.Scores = append(data.Scores, ScoreEntry{Seed: seed, BestScore: score})
	return writeBestScore(data)
}
