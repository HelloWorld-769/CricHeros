package models

// type ShowPlayer struct {
// 	P_Name   string   `json:"p_name"`
// 	P_Age    int64    `json:"p_age"`
// 	JerseyNo int64    `json:"jersey_no"`
// 	PhoneNo  string   `json:"phone_no"`
// 	Country  string   `json:"country"`
// 	MPlayed  int64    `json:"m_played"`
// 	Team     []string `json:"team" gorm:"type:text[]"`
// }

type Player struct {
	P_ID       string `json:"player_id" gorm:"default:uuid_generate_v4();primaryKey"`
	P_Name     string `json:"player_name"`
	MPlayed    int64  `json:"matches_played"`
	P_Age      int64  `json:"player_age"`
	JerseyNo   int64  `json:"jersey_no"`
	PhoneNo    string `json:"phone_no"`
	Country    string `json:"country"`
	T_ID       string `json:"t_id"`
	Is_Captain bool   `json:"is_captain" gorm:"default:false"`
}
type BattingCareer struct {
	P_ID       string  `json:"player_id"`
	IngBat     int64   `json:"inning_bat"`
	RunScored  int64   `json:"run_scored"`
	HScored    int64   `json:"highest_score"` //high score
	AvgScore   float64 `json:"average_score"`
	BallsFaced int64   `json:"balls_faced"`
	BatSR      float64 `json:"bat_sr"` //batting strike rate
	Fifites    int64   `json:"fifties"`
	Hundreds   int64   `json:"hundreds"`
	Fours      int64   `json:"fours"`
	Sixes      int64   `json:"sixes"`
}
type BowlingCareer struct {
	P_ID    string  `json:"player_id"`
	IngBowl int64   `json:"inning_bowl"`
	BBowl   int64   `json:"balls_bowled"` //Balls Bowled
	RConced int64   `json:"runs_conced"`  //Runs Conceded
	Wickets int64   `json:"wickets"`
	BowlAvg int64   `json:"bowling_average"`
	BowlSR  float64 `json:"bowl_sr"` //Bowling strike rate
	Economy float64 `json:"economy"`
}
type Match struct {
	M_ID  string `json:"match_id" gorm:"default:uuid_generate_v4();primaryKey"`
	S_ID  string `json:"scorecard_id"` //scorecard related to it
	T1_ID string
	T2_ID string
	Text  string `json:"text"` //who won the match
}

type Team struct {
	T_ID      string `json:"team_id" gorm:"default:uuid_generate_v4()"`
	T_Name    string `json:"team_name"`
	T_Captain string `json:"team_captain"`
	T_Type    string `json:"team_type"`
	P_ID      string `json:"player_id"`
}

type Response struct {
	Player  Player
	Batting BattingCareer
	Bowling BowlingCareer
}
