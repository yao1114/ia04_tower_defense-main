<template>
  <div id="gameContainer">
    <div id="upperMenu">
      <button class="btn" @click="start" id="start">Start</button>
      <span class="param">Money : {{ money }} </span>
      <span class="param">Wave : {{ currentwave }}/{{ wave }} </span>
      <span class="param">HP : {{ hp }} </span>
    </div>
    <div id="canvasContainer">
      <canvas id="Map" ref="map"></canvas> <!-- layer 1 -->
      <canvas id="Game" ref="game"></canvas> <!-- layer 2 -->
    </div>
  </div>
  <div class="map-menu" v-if="selectingMap" @click="() => (selectingMap = false)">
    <img src="@/components/icons/map1.svg" alt="Map1" @click="() => initMap(1)"/>
    <img src="@/components/icons/map2.svg" alt="Map2" @click="() => initMap(2)"/>
  </div>
</template>

<script setup>
// paper.js  ----------------------> curve fitting [Schneider or RDP]
import { onMounted, ref, getCurrentInstance, reactive } from "vue"
import { fabric } from "fabric"
import { paper } from "paper"
import { MAP1, MAP2 } from "@/assets/maps.js"
import axios from "axios"

const GRID_WIDTH = 35
const CANVAS_LEFT = 80
const CANVAS_TOP = 50
const CURRENT_CONTEXT = getCurrentInstance().ctx

const TOWER_TYPE_DICT = {"HIGH_DAMAGE": "tw1", "HIGH_ATTACK_SPEED" : "tw2", "AVERAGE" : "tw3", "SLOW_DOWN" : "tw4", "LONG_ATTACK_DISTANCE" : "tw5"}
const TOWER_ATTACKSPEED_DICT = {"tw1" : 600, "tw2" : 800, "tw3" : 600, "tw4" : 600, "tw5" : 600}
const TOWER_ATTACKDISTANCE_DICT = {"tw1" : 3*GRID_WIDTH, "tw2" : 3*GRID_WIDTH, "tw3" : 3*GRID_WIDTH, "tw4" : 3, "tw5" : 6*GRID_WIDTH}
const TOWER_DAMAGE_DICT = {"tw1" : 100, "tw2" : 60, "tw3" : 60, "tw4" : 60, "tw5" : 60}
const TOWER_DECELERATION_DICT = {"tw1" : 0, "tw2" : 0, "tw3" : 0, "tw4" : 0.2, "tw5" : 0.2}
const TOWER_VALUE_DICT = {"tw1" : 60, "tw2" : 60, "tw3" : 50, "tw4" : 60, "tw5" : 60}
const TOWER_LEVELUP_ATTACKSPEED_DICT = {"tw1" : 5, "tw2" : 10, "tw3" : 10, "tw4" : 5, "tw5" : 5}
const TOWER_LEVELUP_ATTACKDISTANCE_DICT = {"tw1": GRID_WIDTH, "tw2": GRID_WIDTH, "tw3": GRID_WIDTH, "tw4": GRID_WIDTH, "tw5": GRID_WIDTH}
const TOWER_LEVELUP_DAMAGE_DICT = {"tw1": 10, "tw2": 5, "tw3": 10, "tw4": 5, "tw5": 5}
const TOWER_LEVELUP_VALUE_DICT = {"tw1": 50, "tw2": 50, "tw3": 50, "tw4": 50, "tw5": 50}
const twDistanceMap = new Map  // Map(["1",200])

const MONSTER_HP_DICT = { "NORMAL" : 300, "HIGH_UP" : 500, "HIGH_SPEED" : 200}
const MONSTER_MOVESPEED_DICT = { "NORMAL" : 2, "HIGH_UP" : 2, "HIGH_SPEED" : 1}
const MONSTER_VALUE_DICT = { "NORMAL" : 100, "HIGH_UP" : 150, "HIGH_SPEED" : 150}

let MapHeight = 0
let MapWidth = 0

let timer = null
let money = ref(200)
let currentwave = ref(0)
let wave = ref(3)
let hp = ref(10)
let GameLayer = null
let TowerID = 0
let click_flag = 0 // fabric.js，prevents clicking on an element from being judged as clicking on the canvas again

let dead_index = 0 // for the index of dead monster get through websocket

let gamePaused = ref(false)
let selectingMap = ref(true)

//websocket connection
let monster_conn
let player_conn 

// let id_player
// let id_wave = 0
// let start_point
// let end_point
// let max_x
// let max_y
// let blocks
// let block_length = 0

const drawMonster = (res,canvas) => {
  let MonsterID = res.id
  let pos = res.situation.position
  let type = res.situation.monster_type
  let MonsterHP = res.situation.hp.toString()
  //console.log(res)

  let path
  canvas.remove(canvas.getObjects().find((item) => {return item.id == MonsterID})) // delete previous frame
  canvas.renderAll()

  if (type == "NORMAL"){
    path = "./src/components/icons/monster1.svg"
  }else if (type == "HIGH_HP"){
    path = "./src/components/icons/monster2.svg"
  }else if (type == "HIGH_SPEED"){
    path = "./src/components/icons/monster3.svg"
  }
  fabric.Image.fromURL(path, function (img) {
    img.scale(GRID_WIDTH / img.height).set({
      name:"Monster",
      // name: type,
      // id: MonsterID,
      // top: pos[0] * GRID_WIDTH,
      // left: pos[1] * GRID_WIDTH,
      lockMovementX: true,
      lockMovementY: true,
      selectable: false,
      hasControls: false,
    })
    let textHp = new fabric.Text(MonsterHP, {
      fontFamily: "Comic Sans",
      fontSize: 18,
      fill: "blue",
      fontWeight: "bold",
      top: GRID_WIDTH,
    })
    const monsterWithHp = new fabric.Group([img, textHp], {
      name: type,
      id: MonsterID,
      top: pos[0] * GRID_WIDTH,
      left: pos[1] * GRID_WIDTH,
      lockMovementX: true,
      lockMovementY: true,
      selectable: true,
      hasControls: false,
      hoverCursor: "default"
    })
    canvas.add(monsterWithHp)
  })
  canvas.renderAll()
}

const start = () => {
  console.log("Start")
  if (currentwave.value<wave.value){
    let wave_id
    if (currentwave.value == 0)
      wave_id = 1
    else
      wave_id = currentwave.value
    axios
        .post(
          "http://localhost:8000/new_monster_wave",
          JSON.stringify({
            ID: "owo",
            Wave: wave_id,
          }),
          {
            headers: {
              "Content-Type": "application/json",
            },
          }
        )
        .then((response) => {
          console.log(response)
        })
        .catch((error) => {
          console.log(error.message)
        });
    getMonsterPos(GameLayer)
  }
}

const getMonsterPos = (canvas) => {
  console.log("Starting Websocket monster_conn")
  monster_conn = new WebSocket("ws://localhost:8000/monster_situation")
  monster_conn.onmessage = function(event) {
    //console.log(event.data)
    if (event.data instanceof Blob) {
      var blob = event.data;
      var blobReader = new Response(blob).json()
      blobReader.then(res => {
          //console.log(res.situation.position)
          drawMonster(res,canvas)
      })
    }
  }
  monster_conn.onopen = function(event) {
    console.log(event)
    console.log("Successfully connected to the monster_situation websocket server...")
  }
  monster_conn.onclose = function(event) {
    console.log(event)
    console.log("Socket Closed monster_conn:")
  }
  monster_conn.onerror = function(event) {
    console.log(event)
    console.log("Socket Error: ", event)
  }
}
const getPlayerPos = () => {
  console.log("Starting Websocket player_conn")
  player_conn = new WebSocket("ws://localhost:8000/player_situation")
  player_conn.onmessage = function(event) {
    //console.log(event.data)
    if (event.data instanceof Blob) {
      var blob = event.data;
      var blobReader = new Response(blob).json()
      blobReader.then(res => {
          console.log("player: ",res)
          hp.value = res.hp
          money.value = res.money
          currentwave.value = res.current_wave

          let game_status = res.result
          if (game_status == "Lose" ){
            gameOver(-1)
          }else if (game_status == "Win" ){
            gameOver(1)
          }

          console.log("dead: ",res.dead_monsters)
          if (res.dead_monsters){
            let dead_mst_len = res.dead_monsters.length
            for(let i = 0; i < dead_mst_len ; i++){
              GameLayer.remove(GameLayer.getObjects().find((item) => {return item.id == res.dead_monsters[i]})) //删除死亡的怪兽
            }
            // let dead_mst_len = res.dead_monsters.length
            // for(let i = dead_index; i < dead_mst_len ; i++){
            //   GameLayer.remove(GameLayer.getObjects().find((item) => {return item.id == res.dead_monsters[i]})) //删除死亡的怪兽
            // }
            // dead_index = dead_mst_len
          }
      })
    }
  }
  player_conn.onopen = function(event) {
    console.log(event)
    console.log("Successfully connected to the player_situation websocket server...")
  }
  player_conn.onclose = function(event) {
    console.log(event)
    console.log("Socket Closed player_conn:")
  }
  player_conn.onerror = function(event) {
    console.log(event)
    console.log("Socket Error: ", event)
  }
}
const gameOver = (game_flag) => {
  if (game_flag == 1 ){
    console.log("Win")
    let path = "./src/components/icons/win.svg"
    fabric.Image.fromURL(path, function (img) {
      img.scale(2).set({
        top: (MapHeight-img.height*2)/2,
        left: (MapWidth-img.width*2)/2,
        lockMovementX: true,
        lockMovementY: true,
        selectable: false,
        hoverCursor: "default",
        hasControls: false,
      })
      GameLayer.add(img)
    })
    GameLayer.renderAll()
    player_conn.close()
    monster_conn.close()
  }else if (game_flag == -1 ){
    console.log("Lose")
    let path = "./src/components/icons/lose.svg"
    fabric.Image.fromURL(path, function (img) {
      img.scale(2).set({
        top: (MapHeight-img.height*2)/2,
        left: (MapWidth-img.width*2)/2,
        lockMovementX: true,
        lockMovementY: true,
        selectable: false,
        hoverCursor: "default",
        hasControls: false,
      })
      GameLayer.add(img)
    })
    GameLayer.renderAll()
    player_conn.close()
    monster_conn.close()
  }
  game_flag = 0
}

const pause = () => {
  clearInterval(timer)
  gamePaused.value = true
}
const monstersComing = () => {
  window.requestAnimationFrame(monstersComing)
}
const setMap = (ctx, wd_canvas, ht_canvas) => {
  ctx.width = wd_canvas
  ctx.height = ht_canvas
}
const initMap = (mapNumber) => {
  let road
  switch (mapNumber) {
    case 0:
      road = MAP1
      break
    case 1:
      GameLayer.dispose()
      road = MAP1
      break
    case 2:
      GameLayer.dispose()
      road = MAP2
      break
  }

  if (mapNumber != 0){
    axios
        .post(
          "http://localhost:8000/new_player",
          JSON.stringify({
            ID: "owo",
            Map_ID: mapNumber ,
          }),
          {
            headers: {
              "Content-Type": "application/json",
            },
          }
        )
        .then((response) => {
          console.log(response.data)
          // id_player = response.data.id
          // id_wave = response.data.wave
          // start_point = response.data.start_point
          // end_point = response.data.end_point
          // max_x = response.data.max_x
          // max_y = response.data.max_y
          // blocks = response.data.blocks
          // block_length = response.data.blocks.length
        })
        .catch((error) => {
          console.log(error.message)
        });
    getPlayerPos()
  }
  

  console.log("Initialize! map=", mapNumber)
  const wd_canvas = GRID_WIDTH * road[0].length //get global width
  const ht_canvas = GRID_WIDTH * road.length //get global height
  let Map = document.getElementById("Map") //set canvas width & height
  let Game = document.getElementById("Game")
  setMap(Map, wd_canvas, ht_canvas)
  setMap(Game, wd_canvas, ht_canvas)
  GameLayer = new fabric.Canvas("Game")
  GameLayer.on("mouse:down", (options) => {
    console.log("GameLayer mousedown")
    handleGameLayer(options, road, GameLayer)
  })
  drawRoad(road)
  console.log("draw road")
}

const drawRoad = (arr) => {
  //console.log("Road created")
  let MapLayer = CURRENT_CONTEXT.$refs.map
  MapHeight = MapLayer.height
  MapWidth = MapLayer.width
  paper.setup(MapLayer)

  let mapRect = new paper.Rectangle(new paper.Point(0, 0), new paper.Point(MapWidth, MapHeight))
  let path1 = new paper.Path.Rectangle(mapRect)
  path1.fillColor = "#99cc99"
  /*
  for (let i=0;i<arr.length;i++){ //show obstacles
    for (let j=0;j<arr[i].length;j++){
      if (arr[i][j]==0){
        let roadRect = new paper.Rectangle(new paper.Point(j*wd_grid,i*wd_grid), new paper.Point((j+1)*wd_grid,(i+1)*wd_grid))
        let road = new paper.Path.Rectangle(roadRect)
        road.fillColor = '#99cc99';
      }
    }
  }
  */

  var path2 = getRoadPoints(arr, [0, 1], [1, 1])
  path2.fillColor = "#cccccc"
  //path2.simplify()
  path2.smooth()
}
const getRoadPoints = (arr, firstPt, secondPt) => {
  let path = new paper.Path()
  path.strokeColor = "black"
  path.closed = "true"
  let prevPoint = firstPt
  let nextPoint = secondPt
  let count = 0
  path.add(prevPoint.map((x) => x * GRID_WIDTH))
  path.add(nextPoint.map((x) => x * GRID_WIDTH))
  while (true) {
    count = countNeighbour(arr, nextPoint)
    if (count == 5 || count == 3 || count == 6) {
      if (posPrev(prevPoint, nextPoint) == "D") {
        //D->L
        prevPoint = nextPoint.concat()
        nextPoint[0] -= 1
      } else if (posPrev(prevPoint, nextPoint) == "R") {
        //R->U
        prevPoint = nextPoint.concat()
        nextPoint[1] -= 1
      } else if (posPrev(prevPoint, nextPoint) == "U") {
        //U->R
        prevPoint = nextPoint.concat()
        nextPoint[0] += 1
      } else if (posPrev(prevPoint, nextPoint) == "L") {
        //L->D
        prevPoint = nextPoint.concat()
        nextPoint[1] += 1
      }
    } else if (count == 7 || count == 1 || count == 2) {
      if (posPrev(prevPoint, nextPoint) == "D") {
        //D->R
        prevPoint = nextPoint.concat()
        nextPoint[0] += 1
      } else if (posPrev(prevPoint, nextPoint) == "R") {
        //R->D
        prevPoint = nextPoint.concat()
        nextPoint[1] += 1
      } else if (posPrev(prevPoint, nextPoint) == "U") {
        //U->L
        prevPoint = nextPoint.concat()
        nextPoint[0] -= 1
      } else if (posPrev(prevPoint, nextPoint) == "L") {
        //L->U
        prevPoint = nextPoint.concat()
        nextPoint[1] -= 1
      }
    } else if (count == 4) {
      if (posPrev(prevPoint, nextPoint) == "D") {
        //D->U
        prevPoint = nextPoint.concat()
        nextPoint[1] -= 1
      } else if (posPrev(prevPoint, nextPoint) == "R") {
        //R->L
        prevPoint = nextPoint.concat()
        nextPoint[0] -= 1
      } else if (posPrev(prevPoint, nextPoint) == "U") {
        //U->D
        prevPoint = nextPoint.concat()
        nextPoint[1] += 1
      } else if (posPrev(prevPoint, nextPoint) == "L") {
        //L->R
        prevPoint = nextPoint.concat()
        nextPoint[0] += 1
      }
    }
    if (prevPoint.toString() == firstPt.toString() && nextPoint.toString() == secondPt.toString()) {
      break
    }
    if (count == 4 || count == 5 || count == 7) {
      // choose which point to add
      path.add(nextPoint.map((x) => x * GRID_WIDTH))
    }
  }
  return path
}
const posPrev = (prevPoint, nextPoint) => {
  if (prevPoint[0] - nextPoint[0] == -1) {
    return "L"
  } else if (prevPoint[0] - nextPoint[0] == 1) {
    return "R"
  } else if (prevPoint[1] - nextPoint[1] == -1) {
    return "U"
  } else if (prevPoint[1] - nextPoint[1] == 1) {
    return "D"
  }
}
const countNeighbour = (arr, Point) => {
  let pti = Point[1]
  let ptj = Point[0]
  let sum = 0
  let arrHeight = arr?.length
  let arrWidth = arr[0]?.length
  if (getArrValue(arr, pti - 1, ptj - 1, arrHeight, arrWidth) == 1) {  //L-U
    sum += 1
  }
  if (getArrValue(arr, pti, ptj, arrHeight, arrWidth) == 1) {  //R-D
    sum += 1
  }
  if (getArrValue(arr, pti - 1, ptj, arrHeight, arrWidth) == 1) {  //L-D
    sum += 3
  }
  if (getArrValue(arr, pti, ptj - 1, arrHeight, arrWidth) == 1) {  //R-U
    sum += 3
  }
  return sum
}
const getArrValue = (arr, i, j, arrHeight, arrWidth) => {
  if (i < 0 || j < 0 || i > arrHeight - 1 || j > arrWidth - 1) {
    return null
  } else {
    return arr[i][j]
  }
}
const handleGameLayer = (opt, arr, canvas) => {
  let click_i = Math.floor((opt.e.clientY - CANVAS_TOP) / GRID_WIDTH) //缩放时坐标有问题
  let click_j = Math.floor((opt.e.clientX - CANVAS_LEFT) / GRID_WIDTH)
  //console.log(click_i,click_j)
  if (canvas.getActiveObjects().length == 0) {
    //no object selected --> click on the canvas
    if (canvas.getObjects().find((item) => {return item.name == "atkrange"})) {
      canvas.remove(canvas.getObjects().find((item) => {return item.name == "atkrange"}))
      canvas.remove(canvas.getObjects().find((item) => {return item.name == "delete"}))
      canvas.remove(canvas.getObjects().find((item) => {return item.name == "upgrade"}))
      canvas.renderAll()
    } else if (canvas.getObjects().find((item) => {return item.name == "menugroup"}) == null && arr[click_i][click_j] == 0 && click_flag == 0) {
      // menu doesn't exist && click on the obstacle
      click_flag = 1
      let click_point = new fabric.Rect({
        name: "click_point",
        fill: "orange",
        width: GRID_WIDTH,
        height: GRID_WIDTH,
      })
      let tower_menu = new fabric.Rect({
        name: "tower_menu",
        fill: "green",
        top: (click_i * GRID_WIDTH - 250 > 0) ? -250 : 1.5 * GRID_WIDTH,
        opacity: 0.8,
        width: 340,
        height: 250 - GRID_WIDTH/2
      })
      if (click_j * GRID_WIDTH + 170 - GRID_WIDTH / 2 > MapWidth) // menu position
        tower_menu.left = GRID_WIDTH - 340
      else if (click_j * GRID_WIDTH > 170 - GRID_WIDTH / 2)
        tower_menu.left = GRID_WIDTH / 2 - 170
      else 
        tower_menu.left = 0
      let group = new fabric.Group([click_point, tower_menu], {
        // Consider the position of the delete button, ex: menu near the border and in the range
        name: "menugroup",
        lockMovementX: true,
        lockMovementY: true,
        top: (click_i * GRID_WIDTH - 250 > 0) ? click_i * GRID_WIDTH + tower_menu.top : click_i * GRID_WIDTH,
        //top: click_i * GRID_WIDTH + tower_menu.top,
        left: click_j * GRID_WIDTH + tower_menu.left,
        selectable: false,
        hasControls: false,
        hoverCursor: "default"
      })
      canvas.add(group)
      const path1 = "./src/components/icons/tw1.svg"
      drawChooseMenu(canvas, click_i, click_j, 0, 0, path1, arr)
      const path2 = "./src/components/icons/tw2.svg"
      drawChooseMenu(canvas, click_i, click_j, 0, 1, path2, arr)
      const path3 = "./src/components/icons/tw3.svg"
      drawChooseMenu(canvas, click_i, click_j, 0, 2, path3, arr)
      const path4 = "./src/components/icons/tw4.svg"
      drawChooseMenu(canvas, click_i, click_j, 1, 0, path4, arr)
      const path5 = "./src/components/icons/tw5.svg"
      drawChooseMenu(canvas, click_i, click_j, 1, 1, path5, arr)
      canvas.renderAll()
    } else {
      click_flag = 0
      deleteChooseMenu(canvas)
      canvas.renderAll()
    }
  }
}
const drawChooseMenu = (canvas, i, j, menu_i, menu_j, path, arr) => {
  let re = new RegExp("icons\/\(.*\).svg")
  let Name = "menu" + re.exec(path)[1]
  //console.log(Name)
  const menuItemWidth = 100
  const menuItemMargin = 10
  fabric.Image.fromURL(path, function (img) {
    img.scale(menuItemWidth / img.width).set({
      name: Name,
      // top: i*wd_grid-250+menu_i*(menuItemWidth+menuItemMargin)+menuItemMargin,
      // left: j*wd_grid+wd_grid/2-170+menu_j*(menuItemWidth+menuItemMargin)+menuItemMargin,
      lockMovementX: true,
      lockMovementY: true,
      selectable: true,
      hasControls: false,
      hoverCursor: "pointer",
    })
    var text = TOWER_VALUE_DICT[re.exec(path)[1]].toString()
    let textMoney = new fabric.Text(text, {
      fontFamily: "Comic Sans",
      fontSize: 20,
      fill: "white",
      fontWeight: "bold",
    })
    const towerWithMoney = new fabric.Group([img, textMoney], {
      name: "menu" + re.exec(path)[1],
      top: (i * GRID_WIDTH - 250 > 0) ? i * GRID_WIDTH - 250 + menu_i * (menuItemWidth + menuItemMargin) + menuItemMargin :  (i+1.5) * GRID_WIDTH + menu_i * (menuItemWidth + menuItemMargin) + menuItemMargin,
      //left: j * GRID_WIDTH + GRID_WIDTH / 2 - 170 + menu_j * (menuItemWidth + menuItemMargin) + menuItemMargin,
      lockMovementX: true,
      lockMovementY: true,
      selectable: true,
      hasControls: false,
      hoverCursor: "pointer",
    })
    if (j * GRID_WIDTH + 170 - GRID_WIDTH / 2 > MapWidth) // menu position
      towerWithMoney.left = j * GRID_WIDTH + GRID_WIDTH - 340 + menu_j * (menuItemWidth + menuItemMargin) + menuItemMargin
    else if (j * GRID_WIDTH > 170 - GRID_WIDTH / 2)
      towerWithMoney.left = j * GRID_WIDTH + GRID_WIDTH / 2 - 170 + menu_j * (menuItemWidth + menuItemMargin) + menuItemMargin
    else 
      towerWithMoney.left = j * GRID_WIDTH + menu_j * (menuItemWidth + menuItemMargin) + menuItemMargin

    towerWithMoney.on("selected", (options) => {
      drawTower(canvas, arr, i, j, path)
    })
    canvas.add(towerWithMoney)
  })
}
const deleteChooseMenu = (canvas) => {
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "menugroup"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "menutw1"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "menutw2"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "menutw3"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "menutw4"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "menutw5"}))
}
const drawTower = (canvas, arr, i, j, path) => {
  let re = new RegExp("icons\/\(.*\).svg")
  let twName = re.exec(path)[1]
  if (money.value >= TOWER_VALUE_DICT[twName]){
    TowerID += 1
    let local_id = TowerID
    //console.log(twName)
    fabric.Image.fromURL(path, function (img) {
      img.scale(GRID_WIDTH / img.height).set({
        name: twName,
        id: TowerID,
        top: i * GRID_WIDTH,
        left: j * GRID_WIDTH,
        lockMovementX: true,
        lockMovementY: true,
        selectable: true,
        hoverCursor: "pointer",
        hasControls: false,
      })
      img.on("selected", (options) => {
        showRange(canvas, arr, i, j, local_id, twName)
      })
      canvas.add(img)
    })
    // money.value -= TOWER_VALUE_DICT[twName]
    twDistanceMap.set(TowerID.toString(), TOWER_ATTACKDISTANCE_DICT[twName])
    console.log("twDistanceMap",twDistanceMap)
    deleteChooseMenu(canvas)
    canvas.renderAll()
    arr[i][j] = 2 // 2->tower
    //post request tower coord
    let type_tw
    if (twName == "tw1"){
      type_tw = "HIGH_DAMAGE"
    }else if (twName == "tw2"){
      type_tw = "HIGH_ATTACK_SPEED"
    }else if (twName == "tw3"){
      type_tw = "AVERAGE"
    }else if (twName == "tw4"){
      type_tw = "SLOW_DOWN"
    }else if (twName == "tw5"){
      type_tw = "LONG_ATTACK_DISTANCE"
    }
    axios
        .post(
          "http://localhost:8000/new_tower",
          JSON.stringify({
            ID: "owo",
            Tower_ID: local_id,
            Type: type_tw,
            X: i,
            Y: j,
          }),
          {
            headers: {
              "Content-Type": "application/json",
            },
          }
        )
        .then((response) => {
          console.log(response)
        })
        .catch((error) => {
          console.log(error.message)
        });
  }else{
    console.log("No enough money")
  }
}

const deleteTower = (canvas, arr, i, j, ID, twName) => {
  arr[i][j] = 0 //remove tower from the array
  click_flag = 1
  // money.value += TOWER_VALUE_DICT[twName]
  console.log(twName,typeof(twName))
  canvas.remove(canvas.getObjects().find((item) => {return item.id == ID}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "atkrange"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "delete"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "upgrade"}))
  canvas.renderAll()
  axios
      .post(
        "http://localhost:8000/delete_tower",
        JSON.stringify({
          ID: "owo",
          Tower_ID: ID,
        }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      )
      .then((response) => {
        console.log(response)
      })
      .catch((error) => {
        console.log(error.message)
      });

}


const upgradeTower = (canvas, arr, i, j, ID, twName) => {  
  if (money.value >= TOWER_LEVELUP_VALUE_DICT[twName]){
    if (twDistanceMap.get(ID.toString())<TOWER_ATTACKDISTANCE_DICT[twName]+3*TOWER_LEVELUP_ATTACKDISTANCE_DICT[twName]){ // max level 3
      // money.value -= TOWER_LEVELUP_VALUE_DICT[twName]
      twDistanceMap.set(ID.toString(), twDistanceMap.get(ID.toString())+TOWER_LEVELUP_ATTACKDISTANCE_DICT[twName])
      console.log("twDistanceMap",twDistanceMap)
      axios
        .post(
          "http://localhost:8000/level_up_tower",
          JSON.stringify({
            ID: "owo",
            Tower_ID: ID,
          }),
          {
            headers: {
              "Content-Type": "application/json",
            },
          }
        )
        .then((response) => {
          console.log(response)
        })
        .catch((error) => {
          console.log(error.message)
        });
        // canvas.remove(canvas.getObjects().find((item) => {return item.name == "atkrange"}))
        canvas.remove(canvas.getObjects().find((item) => {return item.name == "delete"}))
        canvas.remove(canvas.getObjects().find((item) => {return item.name == "upgrade"}))
        canvas.renderAll()
    }
  }else{
    console.log("No enough money")
  }

}

const showRange = (canvas, arr, i, j, ID, twName) => {
  //get range with request
  let range = twDistanceMap.get(ID.toString())
  deleteChooseMenu(canvas)
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "atkrange"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "delete"}))
  canvas.remove(canvas.getObjects().find((item) => {return item.name == "upgrade"}))
  let AtkRange = new fabric.Circle({
    name: "atkrange",
    radius: range,
    fill: "red",
    opacity: 0.5,
    top: i * GRID_WIDTH - range + GRID_WIDTH / 2,
    left: j * GRID_WIDTH - range + GRID_WIDTH / 2,
    lockMovementX: true,
    lockMovementY: true,
    selectable: false,
    hasControls: false,
    hoverCursor: "default",
  })
  canvas.add(AtkRange)

  let path = "./src/components/icons/delete.svg"
  fabric.Image.fromURL(path, function (img) {
    img.scale(GRID_WIDTH / img.height).set({
      name: "delete",
      left: j * GRID_WIDTH,
      lockMovementX: true,
      lockMovementY: true,
      selectable: true,
      hoverCursor: "pointer",
      hasControls: false,
    })
    if (i * GRID_WIDTH > range)
      img.top = i * GRID_WIDTH - 0.3 * range
    else
      img.top = 0.3 * range + i * GRID_WIDTH
    img.on("selected", (options) => {
      deleteTower(canvas, arr, i, j, ID, twName)
    })
    canvas.add(img)
  })
  canvas.renderAll()
  let curr_level = ((twDistanceMap.get(ID.toString())-TOWER_ATTACKDISTANCE_DICT[twName])/TOWER_LEVELUP_ATTACKDISTANCE_DICT[twName])
  //console.log("curr_level",curr_level)
  let path_up
  switch (curr_level) {
    case 0:
      path_up = "./src/components/icons/num0.svg"
      break
    case 1:
      path_up = "./src/components/icons/num1.svg"
      break
    case 2:
      path_up = "./src/components/icons/num2.svg"
      break
    case 3:
      path_up = "./src/components/icons/num3.svg"
      break
  }
  fabric.Image.fromURL(path_up, function (img) {
    img.scale(GRID_WIDTH / img.height).set({
      name: "upgrade",
      left: j * GRID_WIDTH,
      lockMovementX: true,
      lockMovementY: true,
      selectable: true,
      hoverCursor: "pointer",
      hasControls: false,
    })
    if (i * GRID_WIDTH > range)
      img.top = i * GRID_WIDTH - 0.6 * range
    else
      img.top = 0.6 * range + i * GRID_WIDTH
    img.on("selected", (options) => {
      upgradeTower(canvas, arr, i, j, ID, twName)
    })
    canvas.add(img)
  })
  canvas.renderAll()
}
onMounted(() => {
  initMap(0)
})
</script>

<style scoped>
#gameContainer {
  position: absolute;
  margin-left: 80px;
  margin-top: 50px;
}
#Map {
  position: absolute;
}
#Game {
  position: absolute;
}
#Result {
  position: absolute;
}
#upperMenu {
  position: absolute;
  top: -50px;
  width: 100%;
  display: flex;
}
.btn {
  background-color: #4caf50;
  text-align: center;
  font-size: 20px;
  font-weight: 800;
  border-radius: 8px;
  margin-right: 15px;
  border: 2px solid #4caf50;
  transition-duration: 0.4s;
  box-shadow: 0 8px 16px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
  cursor: pointer;
  margin: 5px;
  width: 100px;
}
.btn:active {
  background-color: #3e8e41;
  box-shadow: 0 5px #666;
  transform: translateY(4px);
}
.btn span {
  cursor: pointer;
  display: inline-block;
  position: relative;
  transition: 0.5s;
}
.btn span:after {
  content: "»";
  position: absolute;
  opacity: 0;
  top: 0;
  right: -20px;
  transition: 0.5s;
}
.btn:hover span {
  padding-right: 25px;
}
.btn:hover span:after {
  opacity: 1;
  right: 0;
}
.param {
  flex: 1;
  text-align: center;
  font-size: 20px;
  font-weight: 800;
  margin: 5px;
}

.map-menu {
  z-index: 100;
  position: relative;
  top: 0;
  width: 100%;
  height: 100vh;
  font-size: 10rem;
  color: #fff;
  text-align: center;
  line-height: 100vh;

  background-color: rgba(255, 255, 255, 0.5);

  display: flex;
  justify-content: space-around;
  align-items: center;

  cursor: pointer;
}
</style>
