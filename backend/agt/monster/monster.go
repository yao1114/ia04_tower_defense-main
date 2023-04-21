package monster

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/_map"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/tower"
)

type Monster_Factory_Interface interface {
	Create_Monster(typ string) Monster_Interface
}

// factory Interface
type Monster_Server struct {
	sync.RWMutex
	Monsters map[string]Monster_Interface
	//Dead_Monsters         []string
	Send_Out_Monster       chan Request_New_Monster_With_Road
	Blood_Channel          chan int
	Money_Channel          chan int
	Treat_Monster_Channel  chan string
	Monster_Killed_Channel chan string
}

type Position_Message struct {
	Position _map.Position
	ID       string // id of monster
}

// monster interface
type Monster_Interface interface {
	Set_Position(_map.Position)             // change the pos of monster for mobiling
	Set_Position_Message() Position_Message // Package Location Information，for tower_attack_monster
	Get_Attack(tower.Attack) bool           // Handling of monsters after receiving an attack
	Move(*Monster_Server)                   // monster moves along the Astar path
	Reach_Destination() bool
	Get_Hp() int
	Get_Move_Speed() int
	Get_Value() int
	Get_Monster_Type() string
}

type Monster_Normal struct {
	ID          string
	MonsterType string
	// fois/millisecond
	Hp                 int
	MoveSpeed          int
	Position           _map.Position
	Value              int
	Road               []_map.Position // path followed by monster
	Current_Road_Index int
}

type Monster_High_Hp struct {
	ID          string
	MonsterType string
	// fois/millisecond
	Hp                 int
	MoveSpeed          int
	Position           _map.Position
	Value              int
	Road               []_map.Position
	Current_Road_Index int
}

type Monster_High_Speed struct {
	ID          string
	MonsterType string
	// fois/millisecond
	Hp                 int
	MoveSpeed          int
	Position           _map.Position
	Value              int
	Road               []_map.Position
	Current_Road_Index int
}

type Request_New_Monster_With_Road struct {
	Wave int
	Road []_map.Position
}

func New_Monster_Server() *Monster_Server {
	monster_map := make(map[string]Monster_Interface)
	send_out_chan := make(chan Request_New_Monster_With_Road)
	blood_channel := make(chan int)
	money_channel := make(chan int)
	treat_monster_channel := make(chan string)
	monster_killed_channel := make(chan string)
	return &Monster_Server{
		Monsters: monster_map,
		//Dead_Monsters:         dead_monsters,
		Send_Out_Monster:       send_out_chan,
		Blood_Channel:          blood_channel,
		Money_Channel:          money_channel,
		Treat_Monster_Channel:  treat_monster_channel,
		Monster_Killed_Channel: monster_killed_channel,
	}
}

func New_Monster_Normal(id string, position_x int, position_y int, road []_map.Position) *Monster_Normal {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	return &Monster_Normal{
		ID:                 id,
		MonsterType:        "NORMAL",
		Hp:                 300,
		MoveSpeed:          2,
		Position:           position,
		Value:              30,
		Road:               road,
		Current_Road_Index: 0,
	}
}

func New_Monster_High_Hp(id string, position_x int, position_y int, road []_map.Position) *Monster_High_Hp {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	return &Monster_High_Hp{
		ID:                 id,
		MonsterType:        "HIGH_HP",
		Hp:                 500,
		MoveSpeed:          2,
		Position:           position,
		Value:              50,
		Road:               road,
		Current_Road_Index: 0,
	}
}

func New_Monster_High_Speed(id string, position_x int, position_y int, road []_map.Position) *Monster_High_Speed {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	return &Monster_High_Speed{
		ID:                 id,
		MonsterType:        "HIGH_SPEED",
		Hp:                 200,
		MoveSpeed:          1,
		Position:           position,
		Value:              50,
		Road:               road,
		Current_Road_Index: 0,
	}
}

func (monster *Monster_Normal) Set_Position_Message() Position_Message {
	return Position_Message{
		Position: monster.Position,
		ID:       monster.ID,
	}
}

func (monster *Monster_High_Hp) Set_Position_Message() Position_Message {
	return Position_Message{
		Position: monster.Position,
		ID:       monster.ID,
	}
}

func (monster *Monster_High_Speed) Set_Position_Message() Position_Message {
	return Position_Message{
		Position: monster.Position,
		ID:       monster.ID,
	}
}

func (monster *Monster_Normal) Get_Attack(attack tower.Attack) bool {
	// monster dead
	if monster.Hp-attack.Damage <= 0 {
		return true
	} else {
		monster.Hp -= attack.Damage
		return false
	}
}

func (monster *Monster_High_Hp) Get_Attack(attack tower.Attack) bool {
	// monster dead
	if monster.Hp-attack.Damage <= 0 {
		return true
	} else {
		monster.Hp -= attack.Damage
		return false
	}
}

func (monster *Monster_High_Speed) Get_Attack(attack tower.Attack) bool {
	// monster dead
	if monster.Hp-attack.Damage <= 0 {
		return true
	} else {
		monster.Hp -= attack.Damage
		return false
	}
}

func (monster *Monster_Normal) Set_Position(position _map.Position) {
	monster.Position.X = position.X
	monster.Position.Y = position.Y
	monster.Current_Road_Index += 1
	fmt.Println("monster", monster.ID, ": I'm here:", position.X, position.Y, " and my current index: ", monster.Current_Road_Index)
}

func (monster *Monster_High_Hp) Set_Position(position _map.Position) {
	monster.Position.X = position.X
	monster.Position.Y = position.Y
	monster.Current_Road_Index += 1
	fmt.Println("monster", monster.ID, ": I'm here:", position.X, position.Y, " and my current index: ", monster.Current_Road_Index)
}

func (monster *Monster_High_Speed) Set_Position(position _map.Position) {
	monster.Position.X = position.X
	monster.Position.Y = position.Y
	monster.Current_Road_Index += 1
	fmt.Println("monster", monster.ID, ": I'm here:", position.X, position.Y, " and my current index: ", monster.Current_Road_Index)
}

func (monster *Monster_Normal) Move(ms *Monster_Server) {
	tick := time.Tick(time.Duration(monster.MoveSpeed) * time.Second)
	<-tick
	// move monster
	if monster.Current_Road_Index < len(monster.Road)-1 {
		monster.Set_Position(monster.Road[monster.Current_Road_Index+1])
	}
	// If the end point is reached, send the ID to the channel and let it perform the blood deduction
	if monster.Reach_Destination() {
		ms.Treat_Monster_Channel <- monster.ID // inform server: I reached the end of the line
	}
}

func (monster *Monster_High_Hp) Move(ms *Monster_Server) {
	tick := time.Tick(time.Duration(monster.MoveSpeed) * time.Second)
	<-tick
	if monster.Current_Road_Index < len(monster.Road)-1 {
		monster.Set_Position(monster.Road[monster.Current_Road_Index+1])
	}
	if monster.Reach_Destination() {
		ms.Treat_Monster_Channel <- monster.ID // inform server: reach ending point
	}
}

func (monster *Monster_High_Speed) Move(ms *Monster_Server) {
	tick := time.Tick(time.Duration(monster.MoveSpeed) * time.Second)
	<-tick
	if monster.Current_Road_Index < len(monster.Road)-1 {
		monster.Set_Position(monster.Road[monster.Current_Road_Index+1])
	}
	if monster.Reach_Destination() {
		ms.Treat_Monster_Channel <- monster.ID // inform server: reach ending point
	}
}

func (monster *Monster_Normal) Reach_Destination() bool {
	len := len(monster.Road)
	destination := monster.Road[len-1]
	position := monster.Position
	if position.X == destination.X && position.Y == destination.Y {
		return true
	}
	return false
}

func (monster *Monster_High_Hp) Reach_Destination() bool {
	len := len(monster.Road)
	destination := monster.Road[len-1]
	position := monster.Position
	if position.X == destination.X && position.Y == destination.Y {
		return true
	}
	return false
}

func (monster *Monster_High_Speed) Reach_Destination() bool {
	len := len(monster.Road)
	destination := monster.Road[len-1]
	position := monster.Position
	if position.X == destination.X && position.Y == destination.Y {
		return true
	}
	return false
}

func (monster *Monster_Normal) Get_Hp() int {
	return monster.Hp
}

func (monster *Monster_High_Hp) Get_Hp() int {
	return monster.Hp
}

func (monster *Monster_High_Speed) Get_Hp() int {
	return monster.Hp
}

func (monster *Monster_Normal) Get_Move_Speed() int {
	return monster.MoveSpeed
}

func (monster *Monster_High_Hp) Get_Move_Speed() int {
	return monster.MoveSpeed
}

func (monster *Monster_High_Speed) Get_Move_Speed() int {
	return monster.MoveSpeed
}

func (monster *Monster_Normal) Get_Monster_Type() string {
	return monster.MonsterType
}

func (monster *Monster_High_Hp) Get_Monster_Type() string {
	return monster.MonsterType
}

func (monster *Monster_High_Speed) Get_Monster_Type() string {
	return monster.MonsterType
}

func (monster *Monster_Normal) Get_Value() int {
	return monster.Value
}

func (monster *Monster_High_Hp) Get_Value() int {
	return monster.Value
}

func (monster *Monster_High_Speed) Get_Value() int {
	return monster.Value
}

func (ms *Monster_Server) Create_Monster(typ string, position_x int, position_y int, road []_map.Position, wave int) Monster_Interface {
	switch typ {
	case "Monster_Normal":
		id := "wave" + strconv.Itoa(wave) + " " + strconv.Itoa(len(ms.Monsters)+1)
		monster := New_Monster_Normal(id, position_x, position_y, road)
		ms.Monsters[id] = monster
		return monster
	case "Monster_High_Hp":
		id := "wave" + strconv.Itoa(wave) + " " + strconv.Itoa(len(ms.Monsters)+1)
		monster := New_Monster_High_Hp(id, position_x, position_y, road)
		ms.Monsters[id] = monster
		return monster
	case "Monster_High_Speed":
		// Monster's number+1
		id := "wave" + strconv.Itoa(wave) + " " + strconv.Itoa(len(ms.Monsters)+1)
		monster := New_Monster_High_Speed(id, position_x, position_y, road)
		// add new monster into map
		ms.Monsters[id] = monster
		return monster
	default:
		return nil
	}
}

func (ms *Monster_Server) Treat_Dead_Monster(id string) {
	monster, ok := ms.Monsters[id]
	if ok {
		fmt.Println(monster.Set_Position_Message().ID, " : I'm dead !")
		ms.Money_Channel <- ms.Monsters[id].Get_Value() // send Money earned from killing monsters to restserveragent
		ms.Lock()
		ms.Monster_Killed_Channel <- id
		delete(ms.Monsters, id)
		ms.Unlock()
	}
}

func (ms *Monster_Server) Treat_Reach_Monster(id string) {
	monster, ok := ms.Monsters[id]
	if ok {
		fmt.Println(monster.Set_Position_Message().ID, " : arrives at the destination !")
		ms.Blood_Channel <- 1
		ms.Lock()
		delete(ms.Monsters, id)
		ms.Unlock()
	}
}

func (ms *Monster_Server) Move_Monster() {
	for _, monster := range ms.Monsters {
		go monster.Move(ms)
	}
}

// according to wave, add monster into Monster_Server, round++
func (ms *Monster_Server) handlerRequest(req Request_New_Monster_With_Road) {
	MonsterNb := req.Wave*2 + 50
	start := req.Road[0]
	for i := 1; i <= MonsterNb; i++ {
		MonsterType := rand.Intn(3)
		switch MonsterType {
		case 0:
			monster := ms.Create_Monster("Monster_Normal", start.X, start.Y, req.Road, req.Wave)
			ms.Monsters[monster.Set_Position_Message().ID] = monster
		case 1:
			monster := ms.Create_Monster("Monster_High_Hp", start.X, start.Y, req.Road, req.Wave)
			ms.Monsters[monster.Set_Position_Message().ID] = monster
		case 2:
			monster := ms.Create_Monster("Monster_High_Speed", start.X, start.Y, req.Road, req.Wave)
			ms.Monsters[monster.Set_Position_Message().ID] = monster
		default:
			return
		}
	}
	fmt.Println(ms.Monsters)
}

func (ms *Monster_Server) Run() {
	log.Println("Monster_Server start working")
	for {
		select {
		case req := <-ms.Send_Out_Monster:
			ms.handlerRequest(req)
		case req := <-ms.Treat_Monster_Channel:
			ms.Treat_Reach_Monster(req)
		default:
			time.Sleep(time.Second) // can be ajusted
			ms.Move_Monster()       // One step at a time
		}
	}
}
