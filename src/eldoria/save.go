package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SaveData struct {
	CurrentLevel int `json:"current_level"`
	PlayerHP     int `json:"player_hp"`
	PlayerMaxHP  int `json:"player_max_hp"`
}

func getSavePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "eldoria_save.json"
	}

	saveDir := filepath.Join(configDir, "Eldoria")
	err = os.MkdirAll(saveDir, 0755)

	if err != nil {
		return "eldoria_save.json"
	}

	return filepath.Join(saveDir, "save.json")
}

func SaveGame(game *Game) error {
	if game == nil {
		return nil
	}

	if game.Player == nil {
		return nil
	}

	data := SaveData{
		CurrentLevel: game.CurrentLevel,
		PlayerHP:     game.Player.HP,
		PlayerMaxHP:  game.Player.MaxHP,
	}

	savePath := getSavePath()

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(data)
}

func LoadSave() (SaveData, bool) {
	savePath := getSavePath()

	file, err := os.Open(savePath)
	if err != nil {
		return SaveData{}, false
	}

	defer file.Close()

	var data SaveData

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&data)
	if err != nil {
		return SaveData{}, false
	}

	return data, true
}

func HasSave() bool {
	savePath := getSavePath()

	_, err := os.Stat(savePath)

	return err == nil
}

func DeleteSave() error {
	savePath := getSavePath()

	err := os.Remove(savePath)

	if os.IsNotExist(err) {
		return nil
	}

	return err
}
