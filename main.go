package main

import (
	"encoding/json"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const key string = "api_key=RGAPI-2acba61b-513f-4408-b47f-8e05b466e4a4"

var numberOfPulls int = 0

//		*********	 SUMMONER SEARCH STRUCTS	*********		//
type summoner struct {
	ProfileIconId 		int
	Name				string
	Puuid				string
	SummonerLevel		int
	AccountId			string
	Id					string
	RevisionDate 		int
}

//		*********	  MATCH HISTORY STRUCTS		*********		//
type matchHistory struct {
	Matches			[]match
	EndIndex		int
	StartIndex		int
	TotalGames		int
}

type match struct {
	Lane			string
	GameId			int
	Champion		int
	PlatformId		string
	Timestamp		int
	Queue			int
	Role			string
	Season			int
}

//		*********	   MATCH STATS STRUCTS		*********		//
type matchStats struct {
	SeasonId				int
	QueueId					int
	GameId					int
	ParticipantIdentities	[]participantIdentities
	GameVersion				string
	PlatformId				string
	GameMode				string
	MapId					int
	GameType				string
}

type participantIdentities struct {
	Player					player
	ParticipantId			int
}

type player struct {
	CurrentPlatformId		string
	SummonerName			string
	MatchHistoryUri			string
	PlatformId				string
	CurrentAccountId		string
	ProfileIcon				int
	SummonerId				string
	AccountId				string
}

//		*********	  	TIMELINE STRUCTS		*********		//
type timeline struct {
	Frames				[]frame
	FrameInterval		int
}

type frame struct {
	Timestamp			int
	ParticipantFrames	participantFrames
	Events				[]timelineEvents
}

type participantFrames struct {
	ParticipantFrames1	singularParticipant
	ParticipantFrames2	singularParticipant
	ParticipantFrames3	singularParticipant
	ParticipantFrames4	singularParticipant
	ParticipantFrames5	singularParticipant
	ParticipantFrames6	singularParticipant
	ParticipantFrames7	singularParticipant
	ParticipantFrames8	singularParticipant
	ParticipantFrames9	singularParticipant
	ParticipantFrames10	singularParticipant
}

type singularParticipant struct {
	TotalGold			int
	TeamScore			int
	ParticipantId		int
	Level				int
	CurrentGold			int
	MinionsKilled		int
	DominionScore		int
	Position			position
	Xp					int
	JungleMinionsKilled	int
}

type position struct {
	Y					int
	X					int
}

type timelineEvents struct {
	EventType						string	`json:"eventType"`
	TowerType						string	`json:"towerType"`
	TeamId							int		`json:"teamId"`
	AscendedType					string	`json:"ascendedType"`
	KillerId						int		`json:"killerId"`
	LevelUpType						string	`json:"levelUpType"`
	PointCaptured					string	`json:"pointCaptured"`
	AssistingParticipantIds			[]int	`json:"assistingParticipantIds"`
	WardType						string	`json:"wardType"`
	MonsterType						string	`json:"monsterType"`
	Type							string	`json:"type"`
	SkillSlot						int		`json:"skillSlot"`
	VictimId						int		`json:"victimId"`
	Timestamp						int64	`json:"timestamp"`
	AfterId							int		`json:"afterId"`
	MonsterSubType					string	`json:"monsterSubType"`
	LaneType						string	`json:"laneType"`
	ItemId							int		`json:"itemId"`
	ParticipantId					int		`json:"participantId"`
	BuildingType					string	`json:"buildingType"`
	CreatorId						int		`json:"creatorId"`
	Position						position`json:"position"`
	BeforeId						int		`json:"beforeId"`
}

func urlToStructSummoner(url string)summoner {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var data summoner
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data
}

func urlToStructMatchHistory(url string)matchHistory {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var dataStruct matchHistory
	err = json.Unmarshal(body, &dataStruct)
	if err != nil {
		fmt.Println(err.Error())
	}

	return dataStruct
}

func urlToStructMatchStats(url string)matchStats {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var dataStruct matchStats
	err = json.Unmarshal(body, &dataStruct)
	if err != nil {
		fmt.Println(err.Error())
	}

	return dataStruct
}

func urlToTimeline(url string)timeline {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var dataStruct timeline
	err = json.Unmarshal(body, &dataStruct)
	if err != nil {
		fmt.Println(err.Error())
	}

	return dataStruct
}

func summonerSearch(name string) summoner{
	base := "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/"
	url := base+name+"?"+key
	numberOfPulls++
	return urlToStructSummoner(url)
}

func matchHistorySearch(name string, start int) matchHistory{
	base := "https://euw1.api.riotgames.com/lol/match/v4/matchlists/by-account/"
	filters := "queue=420&beginIndex=" + strconv.Itoa(start) + "&"
	accountId := string(summonerSearch(name).AccountId)
	url := base+accountId+"?"+filters+key
	fmt.Println(url)
	numberOfPulls++
	return urlToStructMatchHistory(url)
}

func matchStatsSearch(matchId int) matchStats{
	gameId  := strconv.Itoa(matchId)
	base := "https://euw1.api.riotgames.com/lol/match/v4/matches/"
	url := base+gameId+"?"+key
	numberOfPulls++
	return urlToStructMatchStats(url)
}

func timelineSearch(matchId int) timeline{
	base := "https://euw1.api.riotgames.com/lol/match/v4/timelines/by-match/"
	url := base+strconv.Itoa(matchId)+"?"+key
	numberOfPulls++
	return urlToTimeline(url)
}

func main() {
	sampleSize := 100
	name := "lnsecure mid"
	m := 0
	totalDeathsGraph := [60]float64{}
	mH0 := matchHistorySearch(name,0)
	mH1 := matchHistorySearch(name,100)
	mH2 := matchHistorySearch(name,200)
	mH3 := matchHistorySearch(name,300)
	mH4 := matchHistorySearch(name,400)
	mH5 := matchHistorySearch(name,500)
	mH6 := matchHistorySearch(name,600)
	mH7 := matchHistorySearch(name,700)
	mH8 := matchHistorySearch(name,800)
	mH9 := matchHistorySearch(name,900)
	mH0.Matches = append(mH0.Matches, mH1.Matches...)
	mH0.Matches = append(mH0.Matches, mH2.Matches...)
	mH0.Matches = append(mH0.Matches, mH3.Matches...)
	mH0.Matches = append(mH0.Matches, mH4.Matches...)
	mH0.Matches = append(mH0.Matches, mH5.Matches...)
	mH0.Matches = append(mH0.Matches, mH6.Matches...)
	mH0.Matches = append(mH0.Matches, mH7.Matches...)
	mH0.Matches = append(mH0.Matches, mH8.Matches...)
	mH0.Matches = append(mH0.Matches, mH9.Matches...)
	fmt.Println(len(mH0.Matches))
	summoner := summonerSearch(name)
	summonerId := summoner.Id
	time.Sleep(2*time.Minute)
	for k := 0; k<sampleSize; k++ {
		playerIdMatch := false
		playerId := 0
		matchStatistics := matchStatsSearch(mH0.Matches[k].GameId)
		for l := 0; playerIdMatch==false; l++{
			if matchStatistics.ParticipantIdentities[l].Player.SummonerId==summonerId {
				playerId = matchStatistics.ParticipantIdentities[l].ParticipantId
				playerIdMatch=true
				}
		}
		timeline := timelineSearch(mH0.Matches[k].GameId)
		if numberOfPulls>95 {
			time.Sleep(2*time.Minute)
			numberOfPulls = 0
		}
		fmt.Println(strconv.Itoa(k) + " games looked at")
		for i := 0; i < len(timeline.Frames); i++ {
			for j := 0; j < len(timeline.Frames[i].Events); j++ {
				if timeline.Frames[i].Events[j].Type == "CHAMPION_KILL" && timeline.Frames[i].Events[j].VictimId == playerId {
					totalDeathsGraph[i]++
					m++
				}
			}
		}
	}

	for i:=1; i< len(totalDeathsGraph); i++{
		totalDeathsGraph[i]=totalDeathsGraph[i]/float64(sampleSize)
	}
	data := make(plotter.Values, 60)
	for i:=0; i<60; i++ {
		data[i]=totalDeathsGraph[i]
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Average deaths per min for: " + name
	p.Y.Label.Text = "Average deaths per min over " + strconv.Itoa(sampleSize) + " ranked games"
	p.X.Label.Text = "Minute"

	w := vg.Points(10)

	barsA, err := plotter.NewBarChart(data, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(2)

	p.Add(barsA)
	p.Legend.Top = true
	p.NominalX("0", "", "2", "", "4", "", "6", "", "8", "", "10", "", "12", "", "14", "", "16", "", "18", "", "20", "", "22", "", "24", "", "26", "", "28", "", "30", "", "32", "", "34", "", "36", "", "38", "", "40", "", "42", "", "44", "", "46", "", "48", "", "50", "", "52", "", "54", "", "56", "", "58", "", "60")

	if err := p.Save(10*vg.Inch, 4*vg.Inch, "barchart.png"); err != nil {
		panic(err)
	}
}
