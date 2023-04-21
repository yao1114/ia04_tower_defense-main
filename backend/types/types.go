package types

// create a new player
type Request_New_Player struct {
	ID     string `json:"id"`     // player_name
	Map_ID int    `json:"map_id"` // Selected map number
}

type Request_Map struct {
	ID string `json:"id"`
}
type Response_New_Map struct {
	ID          string  `json:"id"` // player_name
	Wave        int     `json:"wave"`
	Start_Point []int   `json:"start_point"`
	End_Point   []int   `json:"end_point"`
	Max_X       int     `json:"max_x"`
	Max_Y       int     `json:"max_y"`
	Blocks      [][]int `json:"blocks"`
	// Points      [][]int `json:"points"`
}

type Request_New_Monster struct {
	ID   string `json:"id"` // player_name
	Wave int    `json:"wave"`
}

type Response_New_Monster struct {
	Monster_Normal     []int   `json:"monster_normal"`
	Monster_High_Hp    []int   `json:"monster_high_hp"`
	Monster_High_Speed []int   `json:"monster_high_speed"`
	Road               [][]int `json:"road"`
}

type Response_Monster_Situation struct {
	ID        string            `json:"id"`
	Situation Monster_Situation `json:"situation"`
}

type Monster_Situation struct {
	ID           string `json:"id"`
	Hp           int    `json:"hp"`
	Move_Speed   int    `json:"move_speed"`
	Position     []int  `json:"position"`
	Value        int    `json:"value"`
	Monster_Type string `json:"monster_type"`
}

type Request_New_Tower struct {
	ID       string `json:"id"`
	Tower_ID int    `json:"tower_id"`
	Type     string `json:"type"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type Request_Delete_Tower struct {
	ID       string `json:"id"`
	Tower_ID int    `json:"tower_id"`
}

type Request_Level_Up_Tower struct {
	ID       string `json:"id"`
	Tower_ID int    `json:"tower_id"`
}

type Response_Player_Situation struct {
	ID            string   `json:"id"`
	Result        string   `json:"result"`
	Hp            int      `json:"hp"`
	Money         int      `json:"money"`
	Current_Wave  int      `json:"current_wave"`
	Dead_Monsters []string `json:"dead_monsters"`
}
