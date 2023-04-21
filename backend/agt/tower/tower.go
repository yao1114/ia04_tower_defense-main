package tower

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/_map"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/types"
)

// Factory interface
type Tower_Factory_Interface interface {
	// public method for building instances
	Create(typ string, x int, y int) Tower_Interface

	// actual method for generating instances
	Create_tower(typ string, x int, y int) Tower_Interface

	// save the generated instance
	Register_tower(tower Tower_Interface)

	Destory_tower(Tower interface{}) int
}

type Tower_Server struct {
	NumberTower            int
	TowerList              map[int]Tower_Interface
	Create_Tower_Request   chan types.Request_New_Tower
	Delete_Tower_Request   chan types.Request_Delete_Tower
	Level_Up_Tower_Request chan types.Request_Level_Up_Tower
	Money_Channel          chan int
}

// tower interface
type Tower_Interface interface {
	GetId() int
	Get_Position() _map.Position
	Get_Attack_Info() Attack
	Free() int
	Level_up()
}

type Attack struct {
	Damage         int
	Deceleration   float64 //Speed reduction effect<1
	AttackSpeed    float64 //ms
	AttackDistance int
}

type Tower_Average struct {
	ID        int
	TowerType string
	Position  _map.Position
	Value     int
	Level     int //1 par default
	Attack    Attack
}

type Tower_Slow_Down struct {
	ID        int
	TowerType string
	Position  _map.Position
	Value     int
	Level     int //1 par default
	Attack    Attack
}

type Tower_High_Damage struct {
	ID        int
	TowerType string
	Position  _map.Position
	Value     int
	Level     int //1 par default
	Attack    Attack
}

type Tower_High_Attack_Speed struct {
	ID        int
	TowerType string
	Position  _map.Position
	Value     int
	Level     int //1 par default
	Attack    Attack
}

type Tower_Long_Attack_Distance struct {
	ID        int
	TowerType string
	Position  _map.Position
	Value     int
	Level     int //1 par default
	Attack    Attack
}

func New_Tower_Server() *Tower_Server {
	tower_list := make(map[int]Tower_Interface)
	create_tower_channel := make(chan types.Request_New_Tower)
	delete_tower_channel := make(chan types.Request_Delete_Tower)
	level_up_channel := make(chan types.Request_Level_Up_Tower)
	money_channel := make(chan int)
	return &Tower_Server{
		NumberTower:            0,
		TowerList:              tower_list,
		Create_Tower_Request:   create_tower_channel,
		Delete_Tower_Request:   delete_tower_channel,
		Level_Up_Tower_Request: level_up_channel,
		Money_Channel:          money_channel,
	}
}

func New_Tower_Average(id int, position_x int, position_y int) *Tower_Average {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	attack := Attack{
		AttackSpeed:    600,
		AttackDistance: 3,
		Damage:         60,
		Deceleration:   0,
	}
	return &Tower_Average{
		ID:        id,
		TowerType: "AVERAGE",
		Position:  position,
		Value:     50,
		Level:     1,
		Attack:    attack,
	}

}

func New_Tower_Slow_Down(id int, position_x int, position_y int) *Tower_Slow_Down {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	attack := Attack{
		AttackSpeed:    600,
		AttackDistance: 3,
		Damage:         60,
		Deceleration:   0.2,
	}
	return &Tower_Slow_Down{
		ID:        id,
		TowerType: "SLOW_DOWN",
		Position:  position,
		Value:     60,
		Level:     1,
		Attack:    attack,
	}
}

func New_Tower_High_Damage(id int, position_x int, position_y int) *Tower_High_Damage {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	attack := Attack{
		AttackSpeed:    600,
		AttackDistance: 3,
		Damage:         100,
		Deceleration:   0,
	}
	return &Tower_High_Damage{
		ID:        id,
		TowerType: "HIGH_DAMAGE",
		Position:  position,
		Value:     60,
		Level:     1,
		Attack:    attack,
	}
}

func New_Tower_High_Attack_Speed(id int, position_x int, position_y int) *Tower_High_Attack_Speed {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	attack := Attack{
		AttackSpeed:    800,
		AttackDistance: 3,
		Damage:         60,
		Deceleration:   0,
	}
	return &Tower_High_Attack_Speed{
		ID:        id,
		TowerType: "HIGH_ATTACK_SPEED",
		Position:  position,
		Value:     60,
		Level:     1,
		Attack:    attack,
	}
}

func New_Tower_Long_Attack_Distance(id int, position_x int, position_y int) *Tower_Long_Attack_Distance {
	position := _map.Position{
		X: position_x,
		Y: position_y,
	}
	attack := Attack{
		AttackSpeed:    600,
		AttackDistance: 6,
		Damage:         60,
		Deceleration:   0.2,
	}
	return &Tower_Long_Attack_Distance{
		ID:        id,
		TowerType: "LONG_ATTACK_DISTANCE",
		Position:  position,
		Value:     60,
		Level:     1,
		Attack:    attack,
	}
}

func (tower *Tower_Average) Free() int {
	value := tower.Value
	tower.ID = -1
	tower.TowerType = "Destroyed"
	tower.Position.X = -1
	tower.Position.Y = -1
	tower.Value = -1
	tower.Level = -1
	tower.Attack.AttackSpeed = -1
	tower.Attack.AttackDistance = -1
	tower.Attack.Damage = -1
	tower.Attack.Deceleration = -1
	return value
}

func (tower *Tower_Slow_Down) Free() int {
	value := tower.Value
	tower.ID = -1
	tower.TowerType = "Destroyed"
	tower.Position.X = -1
	tower.Position.Y = -1
	tower.Value = -1
	tower.Level = -1
	tower.Attack.AttackSpeed = -1
	tower.Attack.AttackDistance = -1
	tower.Attack.Damage = -1
	tower.Attack.Deceleration = -1
	return value
}

func (tower *Tower_High_Damage) Free() int {
	value := tower.Value
	tower.ID = -1
	tower.TowerType = "Destroyed"
	tower.Position.X = -1
	tower.Position.Y = -1
	tower.Value = -1
	tower.Level = -1
	tower.Attack.AttackSpeed = -1
	tower.Attack.AttackDistance = -1
	tower.Attack.Damage = -1
	tower.Attack.Deceleration = -1
	return value
}

func (tower *Tower_High_Attack_Speed) Free() int {
	value := tower.Value
	tower.ID = -1
	tower.TowerType = "Destroyed"
	tower.Position.X = -1
	tower.Position.Y = -1
	tower.Value = -1
	tower.Level = -1
	tower.Attack.AttackSpeed = -1
	tower.Attack.AttackDistance = -1
	tower.Attack.Damage = -1
	tower.Attack.Deceleration = -1
	return value
}

func (tower *Tower_Long_Attack_Distance) Free() int {
	value := tower.Value
	tower.ID = -1
	tower.TowerType = "Destroyed"
	tower.Position.X = -1
	tower.Position.Y = -1
	tower.Value = -1
	tower.Level = -1
	tower.Attack.AttackSpeed = -1
	tower.Attack.AttackDistance = -1
	tower.Attack.Damage = -1
	tower.Attack.Deceleration = -1
	return value
}

func (tower *Tower_Average) Get_Position() _map.Position {
	return tower.Position
}

func (tower *Tower_High_Attack_Speed) Get_Position() _map.Position {
	return tower.Position
}

func (tower *Tower_High_Damage) Get_Position() _map.Position {
	return tower.Position
}

func (tower *Tower_Long_Attack_Distance) Get_Position() _map.Position {
	return tower.Position
}

func (tower *Tower_Slow_Down) Get_Position() _map.Position {
	return tower.Position
}

func (tower *Tower_Average) Get_Attack_Info() Attack {
	return tower.Attack
}

func (tower *Tower_High_Attack_Speed) Get_Attack_Info() Attack {
	return tower.Attack
}

func (tower *Tower_High_Damage) Get_Attack_Info() Attack {
	return tower.Attack
}

func (tower *Tower_Long_Attack_Distance) Get_Attack_Info() Attack {
	return tower.Attack
}

func (tower *Tower_Slow_Down) Get_Attack_Info() Attack {
	return tower.Attack
}

func (tower *Tower_Average) GetId() int {
	return tower.ID
}
func (tower *Tower_Slow_Down) GetId() int {
	return tower.ID
}
func (tower *Tower_High_Attack_Speed) GetId() int {
	return tower.ID
}
func (tower *Tower_High_Damage) GetId() int {
	return tower.ID
}
func (tower *Tower_Long_Attack_Distance) GetId() int {
	return tower.ID
}

func (tower *Tower_Average) Level_up() {
	if tower.Level < 3 {
		tower.Level += 1
	}
	tower.Attack.AttackSpeed += 10
	tower.Attack.AttackDistance += 1
	tower.Attack.Damage += 10
	tower.Value += 50
}

func (tower *Tower_Slow_Down) Level_up() {
	if tower.Level < 3 {
		tower.Level += 1
	}
	tower.Attack.AttackSpeed += 5
	tower.Attack.AttackDistance += 1
	tower.Attack.Damage += 5
	tower.Attack.Deceleration += 0.2
	tower.Value += 50
}

func (tower *Tower_High_Attack_Speed) Level_up() {
	if tower.Level < 3 {
		tower.Level += 1
	}
	tower.Attack.AttackSpeed += 10
	tower.Attack.AttackDistance += 1
	tower.Attack.Damage += 5
	tower.Value += 50
}

func (tower *Tower_High_Damage) Level_up() {
	if tower.Level < 3 {
		tower.Level += 1
	}
	tower.Attack.AttackSpeed += 5
	tower.Attack.AttackDistance += 1
	tower.Attack.Damage += 10
	tower.Value += 50
}

func (tower *Tower_Long_Attack_Distance) Level_up() {
	if tower.Level < 3 {
		tower.Level += 1
	}
	tower.Attack.AttackSpeed += 5
	tower.Attack.AttackDistance += 1
	tower.Attack.Damage += 5
	tower.Value += 50
}

// actual method of generating instances
func (ts *Tower_Server) Create_tower(id int, typ string, position_x int, position_y int) (int, Tower_Interface) {
	switch typ {
	case "AVERAGE":
		New_Tower := New_Tower_Average(id, position_x, position_y)
		return New_Tower.Value, New_Tower
	case "SLOW_DOWN":
		New_Tower := New_Tower_Slow_Down(id, position_x, position_y)
		return New_Tower.Value, New_Tower
	case "HIGH_DAMAGE":
		New_Tower := New_Tower_High_Damage(id, position_x, position_y)
		return New_Tower.Value, New_Tower
	case "HIGH_ATTACK_SPEED":
		New_Tower := New_Tower_High_Attack_Speed(id, position_x, position_y)
		return New_Tower.Value, New_Tower
	case "LONG_ATTACK_DISTANCE":
		New_Tower := New_Tower_High_Attack_Speed(id, position_x, position_y)
		return New_Tower.Value, New_Tower
	default:
		return 0, nil
	}
}

// Save the generated instance
func (ts *Tower_Server) Register_tower(tower Tower_Interface) {
	ts.TowerList[tower.GetId()] = tower
	ts.NumberTower += 1
}

// Destory_tower：Delete Tower and return half of the money (value/2)
// @ input param：Any type of tower
// @ output param：Money to be returned to the player
func (ts *Tower_Server) Destory_tower(Tower interface{}) int {
	Value := reflect.ValueOf(Tower).Elem()
	TowerType := Value.Field(1).String()
	TowerId := Value.Field(0).Int()
	TowerValue := Value.Field(3).Int()
	// fmt.Println("Tower Id", TowerId)
	// fmt.Println("Tower Value", TowerValue)
	ts.TowerList[int(TowerId)].Free()
	delete(ts.TowerList, int(TowerId))
	fmt.Println(">>Tower ", TowerType, " ID -", TowerId, " destoryed")
	// switch TowerType {
	// case "AVERAGE":
	// 	fmt.Println(">>Tower_Average destroyed")
	// case "SLOW_DOWN":
	// 	fmt.Println(">>Tower_Slow_Down destroyed")
	// case "HIGH_DAMAGE":
	// 	fmt.Println(">>Tower_High_Damage destroyed")
	// case "HIGH_ATTACK_SPEED":
	// 	fmt.Println(">>Tower_High_Attack_Speed destroyed")
	// case "LONG_ATTACK_DISTANCE":
	// 	fmt.Println(">>Tower_Long_Attack_Distance destroyed")
	// default:
	// 	fmt.Println(">>Type unkown")
	// }
	ts.NumberTower -= 1
	return int(TowerValue / 2)
}

func (ts *Tower_Server) Run() {
	log.Println("Tower_Server start working")
	for {
		select {
		case req_create := <-ts.Create_Tower_Request:
			ts.handlerCreateRequest(req_create)
		case req_delete := <-ts.Delete_Tower_Request:
			ts.handleDeleteRequest(req_delete)
		case req_level_up := <-ts.Level_Up_Tower_Request:
			ts.handleLevelUpRequest(req_level_up)
		default:
			time.Sleep(time.Second)
		}
	}
}

// according to req, create a tower，and deduct the appropriate amount for the server（pipe）
func (ts *Tower_Server) handlerCreateRequest(req types.Request_New_Tower) {
	_, isOk := ts.TowerList[req.Tower_ID]
	if isOk {
		fmt.Println("Tower existed!")
	} else {
		//<- The price of the tower, Money_Channel can be replaced with unsigned_int
		// Deduct money and pass a negative number to the pipeline
		//-> cannot set the value of the tower to -1 when  switch to uint to destroy it
		//value: value of tower
		value, newTower := ts.Create_tower(req.Tower_ID, req.Type, req.X, req.Y)
		ts.Register_tower(newTower)
		ts.Money_Channel <- -value
	}
}

func (ts *Tower_Server) handleDeleteRequest(req types.Request_Delete_Tower) {
	_, isOk := ts.TowerList[req.Tower_ID]
	if isOk {
		TowerValue := ts.TowerList[req.Tower_ID].Free()
		delete(ts.TowerList, req.Tower_ID)
		ts.NumberTower -= 1
		ts.Money_Channel <- int(TowerValue / 2) // Half of money returned for the destruction of the tower
	} else {
		fmt.Println("Tower not found!")
	}
}

func (ts *Tower_Server) handleLevelUpRequest(req types.Request_Level_Up_Tower) {
	_, isOk := ts.TowerList[req.Tower_ID]
	if isOk {
		fmt.Println("handleLevelUpRequest!")
		ts.TowerList[req.Tower_ID].Level_up()
		ts.Money_Channel <- -30
	} else {
		fmt.Println("Tower not found!")
	}
}
