// ========================================================================================
package main

import "fmt"

// mp "gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/_map"

// "gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/interact"
// "gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/monster"
// "gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/tower"

// func test_map() {
// 	presetMap := []string{
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		"X . X X X X X X X X X X X X X X X X X X X X X X X X X",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		"X X X X X X X X X X X X X X X X X X X X X X X X . X X",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
// 	}
// 	start := mp.NewPosition(0, 0)
// 	end := mp.NewPosition(18, 10)
// 	m := mp.NewMap(presetMap, start, end)
// 	ms := monster.New_Monster_Server(&m)
// 	tf := tower.New_Tower_Factory()
// 	ms.Create_Monster("Monster_Normal", 0, 0)
// 	tf.Create("AVERAGE", 3, 0)
// 	tf.Create("AVERAGE", 11, 21)
// 	// fmt.Println(m)
// 	// fmt.Println(ms)
// 	// fmt.Println(tf.TowerList[0].Get_Attack_Info())
// 	searchRoad := mp.NewSearchRoad(&m)
// 	if searchRoad.FindoutRoad() {
// 		// fmt.Println(searchRoad)
// 		go ms.Move_Monster(searchRoad.Trace)
// 		go interact.Attack_Command(tf, ms)
// 	}
// 	defer fmt.Println(ms.Map_Server.Hp)
// 	time.Sleep(time.Minute)
// 	// if searchRoad.FindoutRoad() {
// 	// 	fmt.Println("found！")
// 	// 	fmt.Println(searchRoad.Trace)
// 	// 	//m.PrintMap(searchRoad)
// 	// } else {
// 	// 	fmt.Println("not found")
// 	// }
// }

// func test_tower() {
// 	var tf = &tower.Tower_Factory{
// 		NumberTower: 0,
// 		TowerList:   make(map[int]tower.Tower_Interface, 0),
// 	}

// 	//create
// 	fmt.Println("-----------------------------create-----------------------------")
// 	tower1 := tf.Create("AVERAGE", 10, 18)
// 	tower2 := tf.Create("HIGH_DAMAGE", 10, 18)
// 	tower3 := tf.Create("SLOW_DOWN", 10, 18)
// 	fmt.Println(tower1)
// 	fmt.Println(tower2)
// 	fmt.Println(tower3)

// 	//level up
// 	println("-----------------------------LevelUp-----------------------------")
// 	tower1.Level_up()
// 	tower2.Level_up()
// 	tower3.Level_up()
// 	fmt.Println(tower1)
// 	fmt.Println(tower2)
// 	fmt.Println(tower3)

// 	//destory
// 	println("-----------------------------destory-----------------------------")
// 	Refunds := tf.Destory_tower(tower3)
// 	// println(tower)
// 	fmt.Println("oo......Refunds: ", Refunds)
// 	fmt.Println("oo......tf.TowerList[0]", tf.TowerList[0])
// 	fmt.Println("oo......tf.TowerList[1]", tf.TowerList[1])
// 	fmt.Println("oo......tf.TowerList[2]", tf.TowerList[2])
// 	fmt.Println("oo......tf.NumberTower", tf.NumberTower)
// 	fmt.Println("oo......Value of tower destoryed", tower3)

// }

// func test_monster() {
// 	ms := monster.New_Monster_Server()
// 	tf := tower.New_Tower_Factory()
// 	monster_1 := ms.Create_Monster("Monster_Normal", 10, 10)
// 	monster_2 := ms.Create_Monster("Monster_Normal", 10, 9)
// 	monster_3 := ms.Create_Monster("Monster_Normal", 10, 9)
// 	monster_4 := ms.Create_Monster("Monster_Normal", 10, 9)
// 	monster_5 := ms.Create_Monster("Monster_Normal", 10, 9)
// 	monster_6 := ms.Create_Monster("Monster_Normal", 10, 9)
// 	monster_7 := ms.Create_Monster("Monster_Normal", 10, 9)
// 	monster_8 := ms.Create_Monster("Monster_High_Hp", 10, 9)
// 	fmt.Println(
// 		monster_1,
// 		monster_2,
// 		monster_3,
// 		monster_4,
// 		monster_5,
// 		monster_6,
// 		monster_7,
// 		monster_8,
// 	)
// 	fmt.Println(ms.Monsters)
// 	tower := tf.Create("AVERAGE", 10, 10)
// 	fmt.Println(tower)
// 	//monster.Tower_Attack_Monster(tower, *ms)
// 	fmt.Println(
// 		ms.Monsters,
// 	)
// 	interact.Attack_Command(tf, ms)
// 	time.Sleep(20 * time.Second)
// }

// func test_round() {
// 	RoundTotal := monster.Create_Monster_Round(3)
// 	fmt.Println(RoundTotal.NumberRd)
// 	fmt.Println(RoundTotal.Roundlist[1].Monsters[1])
// 	fmt.Println(RoundTotal.Roundlist[1].Monsters[2])
// 	fmt.Println(RoundTotal.Roundlist[1].Monsters[3])
// 	fmt.Println(RoundTotal.Roundlist[1].Monsters[4])
// 	fmt.Println(RoundTotal.Roundlist[1].Monsters[5])
// 	fmt.Println(RoundTotal.Roundlist[1].Monsters[6])
// 	fmt.Println(RoundTotal.Roundlist[1].Monsters[7])
// 	fmt.Println(RoundTotal.Roundlist[2].Monsters)
// 	fmt.Println(RoundTotal.Roundlist[3].Monsters)

// }

func main() {
	//test_tower()
	//test_monster()
	//test_map()
	//test_round()
	var slice []int
	slice = append(slice, 1, 3, 5)
	fmt.Println(slice)
	var slice_1 []string
	slice_1 = append(slice_1, "1")
	fmt.Println(slice_1)
}
