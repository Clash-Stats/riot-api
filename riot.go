package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type RiotService struct {
	ApiKey string
}

func (rs *RiotService) UrlV4(route string, params ...string) string {
	allParams := append(params, fmt.Sprintf("api_key=%s", rs.ApiKey))
	url := fmt.Sprintf("%s%s?%s", "https://na1.api.riotgames.com", route, strings.Join(allParams, "&"))
	// fmt.Println(url)
	return url
}

func (rs *RiotService) UrlV5(route string, params ...string) string {
	allParams := append(params, fmt.Sprintf("api_key=%s", rs.ApiKey))
	url := fmt.Sprintf("%s%s?%s", "https://americas.api.riotgames.com", route, strings.Join(allParams, "&"))
	// fmt.Println(url)
	return url
}

type SummonerDTO struct {
	AccountId     string `json:"accountId"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	Name          string `json:"name"`
	Id            string `json:"id"`
	Puuid         string `json:"puuid"`
	SummonerLevel int    `json:"summonerLevel"`
}

func (rs *RiotService) SummonerByName(name string) (summoner *SummonerDTO, err error) {
	resp, err := http.Get(rs.UrlV4(fmt.Sprintf("/lol/summoner/v4/summoners/by-name/%s", name)))

	if err != nil {
		log.Println("failed to fetch summoner")
		return nil, err
	}

	if resp.StatusCode >= 300 {
		log.Println("url", resp.Request.URL)
		log.Println("failed to fetch summoner, http error: ", resp.Status)
		return nil, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &summoner)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return summoner, nil
}

func (rs *RiotService) MatchesBySummonerInQueue(summoner *SummonerDTO, queue Queue, count int) (matchIds []string, err error) {

	queueParam := fmt.Sprintf("queue=%d", queue)
	countParam := fmt.Sprintf("count=%d", count)
	uri := rs.UrlV5(fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids", summoner.Puuid), "start=0", countParam, queueParam)

	resp, err := http.Get(uri)

	if err != nil {
		log.Println("failed to fetch match list\n", err)
		return nil, err
	}

	if resp.StatusCode >= 300 {
		log.Println("url", resp.Request.URL)
		log.Println("failed to fetch match list, http error: ", resp.Status)

		return nil, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &matchIds)

	if err != nil {
		log.Println("failed to fetch match list\n", err)
		return nil, err
	}

	return matchIds, nil
}

func (rs *RiotService) MatchDetailsByMatchId(matchId string) (matchDetails *MatchDto, err error) {
	resp, err := http.Get(rs.UrlV5(fmt.Sprintf("/lol/match/v5/matches/%s", matchId)))

	if err != nil {
		log.Println("failed to fetch match details")
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode >= 300 {
		log.Println("url", resp.Request.URL)
		log.Println("failed to fetch match details, http error: ", resp.Status)
		return nil, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &matchDetails)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return matchDetails, nil

}
