package commands

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const url = "https://api.opendota.com/api/players/"

type DotaMMR struct {
	TrackedUntil        interface{} `json:"tracked_until"`
	LeaderboardRank     interface{} `json:"leaderboard_rank"`
	Profile             Profile     `json:"profile"`
	MmrEstimate         MmrEstimate `json:"mmr_estimate"`
	SoloCompetitiveRank int64       `json:"solo_competitive_rank"`
	CompetitiveRank     interface{} `json:"competitive_rank"`
	RankTier            interface{} `json:"rank_tier"`
}

type MmrEstimate struct {
	Estimate int64 `json:"estimate"`
}

type Profile struct {
	AccountID      int64       `json:"account_id"`
	Personaname    string      `json:"personaname"`
	Name           interface{} `json:"name"`
	Plus           bool        `json:"plus"`
	Cheese         int64       `json:"cheese"`
	Steamid        string      `json:"steamid"`
	Avatar         string      `json:"avatar"`
	Avatarmedium   string      `json:"avatarmedium"`
	Avatarfull     string      `json:"avatarfull"`
	Profileurl     string      `json:"profileurl"`
	LastLogin      string      `json:"last_login"`
	Loccountrycode string      `json:"loccountrycode"`
	IsContributor  bool        `json:"is_contributor"`
}

func GetDotaMMR(steamID int) (mmr DotaMMR, err error) {
	resp, err := http.Get(url + strconv.Itoa(steamID))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&mmr)
	return
}
