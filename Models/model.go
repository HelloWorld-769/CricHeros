package models

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Response struct {
	Player Player
	Career Career
	Teams  TeamList
	gorm.Model
}
type Player struct {
	gorm.Model
	P_ID   string `json:"player_id" gorm:"default:uuid_generate_v4();primaryKey"`
	P_Name string `json:"player_name"`

	P_Age    int64  `json:"player_age"`
	JerseyNo int64  `json:"jersey_no"`
	PhoneNo  string `json:"phone_no"`
	Country  string `json:"country"`
}

type TeamList struct {
	gorm.Model
	P_ID string `json:"p_id"`
	T_ID string `json:"t_id"`
}
type Career struct {
	P_ID       string  `json:"player_id"`
	MPlayed    int64   `json:"matches_played"`
	RunScored  int64   `json:"run_scored"`
	HScored    int64   `json:"highest_score"` //high score
	AvgScore   float64 `json:"average_score"`
	BallsFaced int64   `json:"balls_faced"`
	//	BatSR      float64 `json:"bat_sr"` //batting strike rate
	Fifites     int64   `json:"fifties"`
	Hundreds    int64   `json:"hundreds"`
	TwoHundreds int64   `json:"two_hundreds"`
	Fours       int64   `json:"fours"`
	Sixes       int64   `json:"sixes"`
	BBowl       int64   `json:"balls_bowled"` //Balls Bowled
	RConced     int64   `json:"runs_conced"`  //Runs Conceded
	Wickets     int64   `json:"wickets"`
	BowlAvg     float64 `json:"bowling_average"`
	//BowlSR     float64 `json:"bowl_sr"` //Bowling strike rate
	Economy float64 `json:"economy"`

	gorm.Model
}
type MatchRecord struct {
	M_ID string
	S_ID string

	gorm.Model
}
type Match struct {
	M_ID   string `json:"match_id" gorm:"default:uuid_generate_v4();primaryKey"`
	S_ID   string `json:"scorecard_id" gorm:"default:uuid_generate_v4()"` //scorecard related to it
	T1_ID  string `json:"team1_id"`
	T2_ID  string `json:"team2_id"`
	Date   string `json:"date"`
	Venue  string `json:"venue"`
	Text   string `json:"text"` //who won the match/
	Status string `json:"status" gorm:"default:active"`

	gorm.Model
}

type Team struct {
	U_ID      string `json:"user_id"`
	T_ID      string `json:"team_id" gorm:"default:uuid_generate_v4()"`
	T_Name    string `json:"team_name"`
	T_Captain string `json:"team_captain"`
	T_Type    string `json:"team_type"`
	P_ID      string `json:"player_id"`

	gorm.Model
}

type ScoreCard struct {
	S_ID      string `json:"scorecard_id"`
	P_ID      string `json:"player_id"`
	PType     string `json:"player_type"`
	RunScored int64
	Fours     int64   `json:"fours"`
	Sixes     int64   `json:"sixes"`
	SR        float64 `json:"strike_rate"`
	BPlayed   int64   `json:"balls_played"`
	OBowled   int64   `json:"overs_bowled"`
	MOvers    int64   `json:"maiden_overs"`
	RunGiven  int64   `json:"runs_given"`
	Wickets   int64   `json:"wickets"`
	NB        int64   `json:"no_balls"`
	WD        int64   `json:"wide_balls"`
	Eco       float64 `json:"economy"`
	IsOut     string  `json:"is_out" gorm:"default:not_out"`

	gorm.Model
}
type Balls struct {
	B_ID      string `json:"ball_id" gorm:"default:uuid_generate_v4()"`
	M_ID      string
	P_ID      string
	BallType  string
	Runs      int64  //runs on that particular ball
	IsValid   string `json:"is_valid"`
	Over      float64
	BallCount int64
	gorm.Model
}
type Credential struct {
	User_ID  string `json:"user_id" gorm:"default:uuid_generate_v4()"`
	Username string `json:"username"`
	PhoneNo  string `json:"phone_no"`
	Password string `json:"password"`
	gorm.Model
}

type Err struct {
	Message    string
	StatusCode int64
}
type Claims struct {
	UserId   string
	Username string `json:"user_id"`
	jwt.RegisteredClaims
}

type Inning struct {
	T_ID   string
	TScore int64
}

type Toss struct {
	Toss_ID  string `json:"toss_id" gorm:"default:uuid_generate_v4()"`
	M_ID     string `json:"match_id"`
	T1_ID    string `json:"head_team"`
	T2_ID    string `json:"tail_team"`
	Decision string `json:"decision"`
	TossWon  string `json:"toss_won"`
}
