package interact

import (
	"fmt"
	"sync"
	"time"

	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/_map"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/monster"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/tower"
)

type Request struct {
	Attack_Channel  chan tower.Attack // l'attaque mis en place par le tour
	PositionMessage monster.Position_Message
}

type Monster_Agent struct {
	Attack_Channel  chan tower.Attack // egale l‘attribut c dans Request
	Request_Channel chan Request      //channel de Request
	PositionMessage monster.Position_Message
}

type Tower_Agent struct {
	Request_Channel chan Request
	Position        _map.Position // position of tower
	Attack          tower.Attack  // info of attack of tower
}

func New_Tower_Agent(c chan Request, position _map.Position, attack tower.Attack) *Tower_Agent {
	return &Tower_Agent{
		Request_Channel: c,
		Position:        position,
		Attack:          attack,
	}
}

func New_Monster_Agent(c chan Request, position_message monster.Position_Message) *Monster_Agent {
	attack_channel := make(chan tower.Attack)
	return &Monster_Agent{
		Attack_Channel:  attack_channel,
		Request_Channel: c,
		PositionMessage: position_message,
	}
}

func Tower_Attack_Monster(tr tower.Tower_Interface, ms *monster.Monster_Server, wg_1 *sync.WaitGroup) {
	var wg sync.WaitGroup
	c := make(chan Request)
	position := tr.Get_Position()
	tower_agent := New_Tower_Agent(c, position, tr.Get_Attack_Info())
	tower_agent.Start()
	for _, monster := range ms.Monsters {
		wg.Add(1)
		position_message := monster.Set_Position_Message()
		monster_agent := New_Monster_Agent(c, position_message)
		monster_agent.Start(ms, &wg)
	}
	// time.Sleep(5 * time.Second)
	wg_1.Done()
	wg.Wait()
}

func (ag *Monster_Agent) Start(ms *monster.Monster_Server, wg *sync.WaitGroup) {
	go func() {
		ag.Request_Channel <- Request{
			Attack_Channel:  ag.Attack_Channel,
			PositionMessage: ag.PositionMessage,
		}
		attack_info := <-ag.Attack_Channel
		// Blood deduction, slowdown, etc
		fmt.Println("Monster", ag.PositionMessage.ID, " receive an attack", attack_info)
		dead_or_not := ms.Monsters[ag.PositionMessage.ID].Get_Attack(attack_info)
		if dead_or_not {
			// ms.Treat_Dead_or_Reach_Monster(ag.PositionMessage.ID)
			ms.Treat_Dead_Monster(ag.PositionMessage.ID)
		}
		wg.Done()
	}()
}

func (ag *Tower_Agent) Start() {
	go func() {
		for {
			req := <-ag.Request_Channel
			go ag.Attack_Monster(req)
		}
	}()
}

func (ag Tower_Agent) Attack_Monster(req Request) {
	border_inf_x := ag.Position.X - ag.Attack.AttackDistance
	border_sup_x := ag.Position.X + ag.Attack.AttackDistance
	border_inf_y := ag.Position.Y - ag.Attack.AttackDistance
	border_sup_y := ag.Position.Y + ag.Attack.AttackDistance
	if req.PositionMessage.Position.X >= border_inf_x && req.PositionMessage.Position.X <= border_sup_x && req.PositionMessage.Position.Y >= border_inf_y && req.PositionMessage.Position.Y <= border_sup_y {
		req.Attack_Channel <- ag.Attack
	}
}

func One_Attack_Command(tf *tower.Tower_Server, ms *monster.Monster_Server) {
	var wg sync.WaitGroup
	for _, tour := range tf.TowerList {
		wg.Add(1)
		go Tower_Attack_Monster(tour, ms, &wg)
		tick := time.Tick(time.Duration(tour.Get_Attack_Info().AttackSpeed) * time.Millisecond)
		<-tick
	}
	wg.Wait()
}
