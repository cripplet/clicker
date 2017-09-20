package cc_fb

import (
	"encoding/json"
	"fmt"
	"github.com/cripplet/clicker/db/firebase/config"
	"github.com/cripplet/clicker/firebase-db"
	"math/rand"
)

type FBGameState struct {
	ID       string `json:"id"`
	Upgrades string `json:"upgrades"`
}

type FBUser struct {
	ID     string `json:"id"`
	GameID string `json:"game_id"`
}

type CCAuthenticationError struct {
	ID string
}

func (e *CCAuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error: provided incorrect token for game ID %s", e.ID)
}

func randomString(n int) string {
	r := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = r[rand.Intn(len(r))]
	}
	return string(b)
}

func SaveGameState(g FBGameState, u FBUser, new_game bool) (FBGameState, FBUser, error) {
	if new_game {
		g.ID = randomString(32)
		u.ID = randomString(32)
		u.GameID = g.ID
	}

	g_json, err := json.Marshal(g)
	if err != nil {
		return FBGameState{}, FBUser{}, err
	}

	u_json, err := json.Marshal(u)
	if err != nil {
		return FBGameState{}, FBUser{}, err
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
		return FBGameState{}, FBUser{}, err
	}

	_, _, err = firebase_db.Put(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/user/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, u.ID),
		u_json,
		false,
		"",
		map[string]string{},
		&u,
	)
	if err != nil {
		return FBGameState{}, FBUser{}, err
	}

	return g, u, err
}

func LoadGameState(id string, token string) (FBGameState, FBUser, error) {
	g := FBGameState{}
	u := FBUser{}

	_, _, err := firebase_db.Get(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, id),
		false,
		map[string]string{},
		&g,
	)
	if err != nil {
		return FBGameState{}, FBUser{}, err
	}

	_, _, err = firebase_db.Get(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/user/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, token),
		false,
		map[string]string{},
		&u,
	)
	if err != nil {
		return FBGameState{}, FBUser{}, err
	}
	if u.GameID != id {
		return FBGameState{}, FBUser{}, &CCAuthenticationError{
			ID: id,
		}

	}

	return SaveGameState(FBGameState(g), FBUser(u), true)
}
