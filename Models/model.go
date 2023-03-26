package models

type Player struct {
	P_ID     string   `json:"player_id" gorm:"default:uuid_generate_v4();primaryKey"`
	P_Name   string   `json:"player_name"`
	M_ID     []string `json:"m_id" gorm:"type:ARRAY"`
	Runs     int64    `json:"runs"`
	Fours    int64    `json:"fours"`
	Sixes    int64    `json:"sixes"`
	Hundreds int64    `json:"hundreds"`
	Role     string   `json:"role"`
}

type Match struct {
	M_ID  string `json:"m_id" gorm:"default:uuid_generate_v4();primaryKey"`
	S_ID  string `json:"s_id"`
	Date  string `json:"date"`
	Venue string `json:"venue"`
	Team1 Team
	Team2 Team
}

type Team struct {
	T_ID   string   `json:"t_id" gorm:"default:uuid_generate_v4();primaryKey"`
	T_Name string   `json:"t_name"`
	Player []string `gorm:"type:ARRAY"`
}

type LastWicket struct {
	P_ID   int    `json:"p_id"`
	P_Name string `json:"p_name"`
	Score  int64
	Overs  int64
	Text   string `json:"text"`
}

type ScoreCard struct {
	S_ID       string `json:"scoreCard_id" gorm:"default:uuid_generate_v4();primaryKey"`
	Runs       int64
	Wickets    int64
	RunRate    float64
	Overs      float64
	Target     int64
	Fours      int64
	Sixes      int64
	LastWicket LastWicket
}
