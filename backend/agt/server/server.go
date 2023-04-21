package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/_map"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/interact"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/monster"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/tower"
	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/types"
)

type RestServerAgent struct {
	sync.Mutex
	addr                  string
	Game_Agents           map[string]*GameAgent
	Create_Player_Channel chan *GameAgent
}

type GameAgent struct {
	sync.Mutex
	ID             string
	Hp             int
	Money          int
	Result         string
	Current_Wave   int
	Monster_killed []string
	Monster_server *monster.Monster_Server
	Tower_server   *tower.Tower_Server
	Map            *_map.Map
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRestServerAgent(addr string) *RestServerAgent {
	game_agents := make(map[string]*GameAgent)
	create_player_channel := make(chan *GameAgent)
	return &RestServerAgent{
		addr:                  addr,
		Game_Agents:           game_agents,
		Create_Player_Channel: create_player_channel,
	}
}

func (rsa *RestServerAgent) NewGameAgent(id string, mp *_map.Map) {
	var monster_killed []string
	monster_server := monster.New_Monster_Server()
	tower_factory := tower.New_Tower_Server()
	game_agent := &GameAgent{
		ID:             id,
		Hp:             10,
		Money:          200,
		Result:         "Playing",
		Current_Wave:   0,
		Monster_server: monster_server,
		Tower_server:   tower_factory,
		Map:            mp,
		Monster_killed: monster_killed,
	}
	rsa.Game_Agents[id] = game_agent
}

func (*RestServerAgent) Decode_New_Player_Request(r *http.Request) (req types.Request_New_Player, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

// func (*RestServerAgent) Decode_Map_Request(r *http.Request) (req types.Request_Map, err error) {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r.Body)
// 	err = json.Unmarshal(buf.Bytes(), &req)
// 	return
// }

func (*RestServerAgent) Decode_New_Tower_Request(r *http.Request) (req types.Request_New_Tower, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) Decode_Delete_Tower_Request(r *http.Request) (req types.Request_Delete_Tower, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) Decode_Level_Up_Tower_Request(r *http.Request) (req types.Request_Level_Up_Tower, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*RestServerAgent) Decode_New_Monster_Request(r *http.Request) (req types.Request_New_Monster, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

/*
// @description: vérifier si le requete est post ou pas
// @param: method string, w http.ResponseWriter, r *http.Request
*/
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(w, "Not Implemented")
		return false
	}
	return true
}

func (rsa *RestServerAgent) Handle_New_Player_Request(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Reçoit une requete de creer un joueur")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")

	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.Decode_New_Player_Request(r)
	if err != nil {
		// 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Check if the map number exists
	_, ok_map := _map.Map_Instances[req.Map_ID]
	if ok_map {
		mp := _map.NewMap(req.Map_ID)
		// check if the player already exists
		_, ok_player := rsa.Game_Agents[req.ID]
		if ok_player {
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprintln(w, "Player Existes")
			log.Println("Player Existes")
			return
		} else {
			rsa.NewGameAgent(req.ID, &mp)
			// Encapsulate return values
			start := mp.Get_Start_Point()
			end := mp.Get_End_Point()
			start_point := []int{start.X, start.Y}
			end_point := []int{end.X, end.Y}
			// block_list
			block_list := [][]int{}
			for _, block := range mp.Blocks {
				temp := []int{block.X, block.Y}
				block_list = append(block_list, temp)
			}
			resp, err_marshal := json.Marshal(types.Response_New_Map{ID: req.ID, Wave: mp.Wave, Start_Point: start_point, End_Point: end_point, Max_X: mp.MaxX, Max_Y: mp.MaxY, Blocks: block_list})
			if err_marshal != nil {
				panic(err_marshal) // Can be ajusted
			}
			w.Write(resp)
			rsa.Create_Player_Channel <- rsa.Game_Agents[req.ID]
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(w, "Map Not Found")
		log.Println("Map Not Found")
		return
	}
}

func (rsa *RestServerAgent) Handle_New_Tower_Request(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Reçoit une requete de creer un tour")

	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.Decode_New_Tower_Request(r)
	if err != nil {
		// 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	player, ok := rsa.Game_Agents[req.ID]
	if ok {
		//send req to Tower_Server's channel
		player.Tower_server.Create_Tower_Request <- req
		resp, _ := json.Marshal(req)
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Player Not Found")
		log.Println("Player Not Found")
	}
}

func (rsa *RestServerAgent) Handle_Delete_Tower_Request(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Reçoit une requete de supprimer un tour")

	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		fmt.Println("Handle_Delete_Tower_Request: Methode is not post")
		return
	}

	// décodage de la requête
	req, err := rsa.Decode_Delete_Tower_Request(r)
	if err != nil {
		// 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		fmt.Println("Handle_Delete_Tower_Request: err")
		return
	}

	player, ok := rsa.Game_Agents[req.ID]
	if ok {
		//send req to Tower_Server's channel
		player.Tower_server.Delete_Tower_Request <- req
		resp, _ := json.Marshal(req)
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Player Not Found")
		log.Println("Player Not Found")
	}
}

func (rsa *RestServerAgent) Handle_Level_Up_Tower_Request(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Reçoit une requete de Tour Level Up")

	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.Decode_Level_Up_Tower_Request(r)
	if err != nil {
		// 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	player, ok := rsa.Game_Agents[req.ID]
	if ok {
		player.Tower_server.Level_Up_Tower_Request <- req
		resp, _ := json.Marshal(req)
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Player Not Found")
		log.Println("Player Not Found")
	}
}

func (rsa *RestServerAgent) Handle_New_Monster_Request(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Reçoit une requete d'envoyer des monstres")

	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.Decode_New_Monster_Request(r)
	if err != nil {
		// 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	player, ok := rsa.Game_Agents[req.ID]
	if ok {
		// wave cannot exceed the wave set by map
		if req.Wave <= player.Map.Wave && req.Wave > 0 {
			// send req to Tower_Server's channel
			search_road := _map.NewSearchRoad(player.Map)
			// Roads of travel is found or not
			if search_road.FindoutRoad() {
				req_chan := monster.Request_New_Monster_With_Road{Wave: player.Current_Wave, Road: search_road.Trace}
				player.Monster_server.Send_Out_Monster <- req_chan
				player.Current_Wave += 1
			} else {
				w.WriteHeader(500)
				fmt.Fprintln(w, "road for monsters Not Found")
			}
		} else {
			w.WriteHeader(404)
			fmt.Fprintln(w, "Wave not allowed")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Player Not Found")
		log.Println("Player Not Found")
	}

}

func Send_Situation(m monster.Monster_Interface, c *websocket.Conn, wg *sync.WaitGroup) {
	position_message := m.Set_Position_Message()
	id := position_message.ID
	x := position_message.Position.X
	y := position_message.Position.Y
	monster_situation := types.Monster_Situation{ID: id, Hp: m.Get_Hp(), Move_Speed: m.Get_Move_Speed(), Position: []int{x, y}, Value: m.Get_Value(), Monster_Type: m.Get_Monster_Type()}
	tosend, err_marshal := json.Marshal(types.Response_Monster_Situation{ID: id, Situation: monster_situation})
	if err_marshal != nil {
		panic(err_marshal) // can be ajusted
	}
	// countdown
	//tick := time.Tick(time.Duration(monster_situation.Move_Speed) * time.Second)
	tick := time.Tick(50 * time.Millisecond)
	<-tick
	//log.Println("tosend:", tosend)
	c.WriteMessage(websocket.BinaryMessage, tosend)
	wg.Done()
}

func (ga *GameAgent) Send_Monster_Situation(c *websocket.Conn) {
	var wg sync.WaitGroup
	for _, monster := range ga.Monster_server.Monsters {
		ga.Lock()
		wg.Add(1)
		Send_Situation(monster, c, &wg)
		ga.Unlock()
	}
	wg.Wait()
}

func (rsa *RestServerAgent) Handle_Monster_Situation(w http.ResponseWriter, r *http.Request) {
	var mu sync.Mutex
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		for _, ga := range rsa.Game_Agents {
			mu.Lock()
			ga.Send_Monster_Situation(c)
			mu.Unlock()
		}
	}

}

func (ga *GameAgent) Send_Player_Situation(c *websocket.Conn) {
	defer ga.Unlock()
	tosend, err_marshal := json.Marshal(types.Response_Player_Situation{ID: ga.ID, Result: ga.Result, Hp: ga.Hp, Money: ga.Money, Current_Wave: ga.Current_Wave, Dead_Monsters: ga.Monster_killed})
	if err_marshal != nil {
		panic(err_marshal) // can be ajusted
	}
	ga.Lock()
	c.WriteMessage(websocket.BinaryMessage, tosend)
}

func (rsa *RestServerAgent) Handle_Player_Situation(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		for _, ga := range rsa.Game_Agents {
			rsa.Lock()
			ga.Send_Player_Situation(c)
			rsa.Unlock()
		}
		time.Sleep(time.Second)
	}
}

func (ga *GameAgent) Attack_Command() {
	// Time interval: the normal monster movement speed
	for {
		interact.One_Attack_Command(ga.Tower_server, ga.Monster_server)
		tick := time.Tick(time.Second)
		<-tick
	}
}

// the victory will directly end the game
func (ga *GameAgent) Game_Result() {
	for {
		if ga.Hp <= 0 {
			ga.Result = "Lose"
		} else if ga.Current_Wave == ga.Map.Wave && len(ga.Monster_server.Monsters) == 0 && ga.Hp > 0 {
			ga.Result = "Win"
		}
	}
}

// Listening for blood loss in the pipe every second
func (ga *GameAgent) Handle_Blood_Channel() {
	for {
		select {
		case blood_lost := <-ga.Monster_server.Blood_Channel:
			ga.Hp -= blood_lost
			log.Println(ga.ID, " has been attacked and has ", ga.Hp, " Hp now :(")
		default:
			time.Sleep(time.Second)
		}
	}
}

// Listen to the pipeline to know if there are any gold coins obtained from killing monsters
func (ga *GameAgent) Handle_Monster_Money_Channel() {
	for {
		select {
		case money := <-ga.Monster_server.Money_Channel:
			ga.Money += money
			log.Println(ga.ID, " has gained ", money, "$", " and You Have ", ga.Money, " now :)")
		default:
			time.Sleep(time.Second)
		}
	}
}

// Listen to the pipeline for changes in gold caused by building/taking down defense towers
func (ga *GameAgent) Handle_Tower_Money_Channel() {
	for {
		select {
		case money := <-ga.Tower_server.Money_Channel:
			ga.Money += money
			if money >= 0 {
				log.Println(ga.ID, " has recovered ", money, "$", " and has ", ga.Money, " now :)")
			} else {
				log.Println(ga.ID, " has payed ", -money, "$", " and has ", ga.Money, " now :(")
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (ga *GameAgent) Handle_Monster_Killed_Channel() {
	for {
		select {
		case id := <-ga.Monster_server.Monster_Killed_Channel:
			ga.Monster_killed = append(ga.Monster_killed, id)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (ga *GameAgent) Start() {
	go ga.Monster_server.Run()
	go ga.Tower_server.Run()
	go ga.Attack_Command()
	go ga.Handle_Blood_Channel()
	go ga.Handle_Monster_Money_Channel()
	go ga.Handle_Tower_Money_Channel()
	go ga.Game_Result()
	go ga.Handle_Monster_Killed_Channel()
}

func (rsa *RestServerAgent) Handle_Player() {
	for {
		select {
		case player := <-rsa.Create_Player_Channel:
			go player.Start()
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

func (rsa *RestServerAgent) Start() {
	go rsa.Handle_Player()
	mux := http.NewServeMux()
	mux.HandleFunc("/new_player", rsa.Handle_New_Player_Request)
	mux.HandleFunc("/monster_situation", rsa.Handle_Monster_Situation)
	mux.HandleFunc("/new_tower", rsa.Handle_New_Tower_Request)
	mux.HandleFunc("/delete_tower", rsa.Handle_Delete_Tower_Request)
	mux.HandleFunc("/level_up_tower", rsa.Handle_Level_Up_Tower_Request)
	mux.HandleFunc("/new_monster_wave", rsa.Handle_New_Monster_Request)
	mux.HandleFunc("/player_situation", rsa.Handle_Player_Situation)
	handler := cors.Default().Handler(mux)
	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	http.ListenAndServe(":8000", handler)
	go log.Fatal(s.ListenAndServe())
}
