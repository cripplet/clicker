package cc_fb

import (
	"encoding/json"
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/firebase-db"
	"math/rand"
)

type FBGameState struct {
	ID       string `json:"id"`
	Upgrades string `json:"upgrades"`
}

func randomString(n int) string {
	r := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = r[rand.Intn(len(r))]
	}
	return string(b)
}

func SaveGameState(g FBGameState, new_game bool) (FBGameState, error) {
	if new_game {
		g.ID = randomString(32)
	}

	g_json, err := json.Marshal(g)
	if err != nil {
		return FBGameState{}, err
	}

	_, _, err = firebase_db.Put(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, g.ID),
		g_json,
		false,
		"",
		map[string]string{},
		&g,
	)
	if err != nil {
		return FBGameState{}, err
	}

	return g, err
}

func LoadGameState(id string) (FBGameState, error) {
	g := FBGameState{}

	_, _, err := firebase_db.Get(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, id),
		false,
		map[string]string{},
		&g,
	)
	if err != nil {
		return FBGameState{}, err
	}

	return SaveGameState(FBGameState(g), true)
}
