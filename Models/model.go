package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Response struct {
	Status string      `json:"status"`
	Code   int64       `json:"code"`
	Data   interface{} `json:"data"`
}
type PlayerData struct {
	Player    Player
	Career    Career
	Teams     TeamList
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Player struct {
	P_ID   string `json:"playerId" gorm:"default:uuid_generate_v4();primaryKey"`
	P_Name string `json:"playerName" validate:"required"`

	P_Age     int64  `json:"playerAge" validate:"required,gt=18"`
	JerseyNo  int64  `json:"jerseyNo" validate:"required"`
	PhoneNo   string `json:"phoneNo" validate:"required"`
	Country   string `json:"country" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TeamList struct {
	P_ID      string `json:"playerId"`
	T_ID      string `json:"teamId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Career struct {
	P_ID        string  `json:"playerId"`
	MPlayed     int64   `json:"matchesPlayed"`
	RunScored   int64   `json:"runScored"`
	HScored     int64   `json:"highestScore"` //high score
	AvgScore    float64 `json:"averageScore"`
	BallsFaced  int64   `json:"ballsFaced"`
	Fifites     int64   `json:"fifties"`
	Hundreds    int64   `json:"hundreds"`
	TwoHundreds int64   `json:"twoHundreds"`
	Fours       int64   `json:"fours"`
	Sixes       int64   `json:"sixes"`
	BBowl       int64   `json:"ballsBowled"` //Balls Bowled
	RConced     int64   `json:"runsConced"`  //Runs Conceded
	Wickets     int64   `json:"wickets"`
	BowlAvg     float64 `json:"bowlingAverage"`
	Economy     float64 `json:"economy"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type MatchRecord struct {
	M_ID string
	S_ID string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Match struct {
	U_ID   string `json:"userId"`
	M_ID   string `json:"matchId" gorm:"default:uuid_generate_v4();primaryKey"`
	S_ID   string `json:"scorecardId" gorm:"default:uuid_generate_v4()"` //scorecard related to it
	T1_ID  string `json:"team1Id" validate:"required"`
	T2_ID  string `json:"team2Id" validate:"required"`
	Date   string `json:"date"`
	Venue  string `json:"venue" validate:"required"`
	Text   string `json:"text"` //who won the match/
	Status string `json:"status" gorm:"default:active"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Team struct {
	U_ID      string `json:"userId" `
	T_ID      string `json:"teamId" gorm:"default:uuid_generate_v4()"`
	T_Name    string `json:"teamName" validate:"required"`
	T_Captain string `json:"teamCaptain" validate:"required" `
	T_Type    string `json:"teamType" validate:"required"`
	P_ID      string `json:"playerId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ScoreCard struct {
	S_ID      string  `json:"scorecardId"`
	P_ID      string  `json:"playerId"`
	PType     string  `json:"playerType" validate:"required"`
	RunScored int64   `json:"runScored"`
	Fours     int64   `json:"fours"`
	Sixes     int64   `json:"sixes"`
	SR        float64 `json:"strikeRate"`
	BPlayed   int64   `json:"ballsPlayed"`
	OBowled   int64   `json:"oversBowled"`
	MOvers    int64   `json:"maidenOvers"`
	RunGiven  int64   `json:"runsGiven"`
	Wickets   int64   `json:"wickets"`
	NB        int64   `json:"noBalls"`
	WD        int64   `json:"wideBalls"`
	Eco       float64 `json:"economy"`
	IsOut     string  `json:"isOut" gorm:"default:not_out"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Balls struct {
	B_ID      string  `json:"ballId" gorm:"default:uuid_generate_v4()"`
	M_ID      string  `json:"matchId" `
	P_ID      string  `json:"playerId"  `
	BallType  string  `json:"ballType"  `
	Runs      int64   `json:"runs"  ` //runs on that particular ball
	IsValid   string  `json:"isValid"  `
	Over      float64 `json:"over"`
	BallCount int64   `json:"ballCount"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CardData struct {
	M_ID      string `json:"matchId" validate:"required"`
	Batsmen   string `json:"batsmen" validate:"required"`
	Baller    string `json:"baller" validate:"required"`
	Runs      int64  `json:"runs" validate:"required oneof=1 2 3 4 5 6 7"`
	Ball_Type string `json:"ballType" validate:"required"`
	PrevRuns  int64  `json:"prevRuns"`
}
type Credential struct {
	User_ID   string `json:"user_id" gorm:"default:uuid_generate_v4()"`
	Username  string `json:"userName"  validate:"required"`
	Email     string `json:"email" gorm:"unique"  validate:"required,email"`
	Role      string `json:"role"`
	Password  string `json:"password" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Claims struct {
	UserID string
	Role   string
	jwt.RegisteredClaims
}

type Inning struct {
	M_ID   string
	T_ID   string
	TScore int64
}

type Toss struct {
	Toss_ID  string `json:"toss_id" gorm:"default:uuid_generate_v4()"`
	M_ID     string `json:"match_id" validate:"required"`
	T1_ID    string `json:"head_team"  validate:"required"`
	T2_ID    string `json:"tail_team"  validate:"required"`
	Decision string `json:"decision"`
	TossWon  string `json:"toss_won"`
}
