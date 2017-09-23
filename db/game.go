package cc_fb

import (
	"encoding/json"
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/firebase-db"
	"math/rand"
)

type FBGameState struct {
	ID    string `json:"id"`
	Exist bool   `json:"exist"`
}

func randomString(n int) string {
	r := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = r[rand.Intn(len(r))]
	}
	return string(b)
}

type PostID struct {
	Name string `json:"name"`
}

func NewGameState() (FBGameState, error) {
	g := FBGameState{
		Exist: true,
	}
	g_json, err := json.Marshal(&g)
	if err != nil {
		return FBGameState{}, err
	}

	p := PostID{}

	_, _, err = firebase_db.Post(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, g.ID),
		g_json,
		false,
		map[string]string{},
		&p,
	)
	if err != nil {
		return FBGameState{}, err
	}

	g.ID = p.Name
	return g, nil
}

func LoadGameState(id string) (FBGameState, error) {
	if id == "" {
		return NewGameState()
	}

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
	if g.Exist {
		g.ID = id
	}
	return g, nil
}
