package _map

import (
	"container/heap"
	"math"
	"strconv"
	"strings"
)

type Position struct {
	X int
	Y int
}

// Used for A* Algorithme
type _AstarPoint struct {
	position Position
	father   *_AstarPoint
	gVal     int
	hVal     int
	fVal     int
}

type OpenList []*_AstarPoint

type SearchRoad struct {
	theMap  *Map
	start   _AstarPoint
	end     _AstarPoint
	closeLi map[string]*_AstarPoint
	openLi  OpenList
	openSet map[string]*_AstarPoint
	TheRoad []*_AstarPoint
	Trace   []Position
}

// var Maps = map[string]int{
// 	"map1": 0,
// 	"map2": 1,
// }

type Map struct {
	// id         string
	StartPoint Position
	EndPoint   Position
	MaxX       int
	MaxY       int
	Blocks     map[string]*Position
	Points     [][]Position
	Wave       int
}

// convert the point into a key
func pointAsKey(x, y int) (key string) {
	key = strconv.Itoa(x) + "," + strconv.Itoa(y)
	return key
}

func NewPosition(x int, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

func (mp *Map) Get_Start_Point() Position {
	return mp.StartPoint
}

func (mp *Map) Get_End_Point() Position {
	return mp.EndPoint
}

// DIY maps
type Map_Management struct {
	Map_String  []string
	Start_Point []int
	End_Point   []int
	Wave        int
}

var Map_Instances = map[int]Map_Management{
	1: Map_Management{
		Map_String: []string{
			"X . . . . . . . . . . . . X X X X X X X X X X X X X X X X X",
			". . . . . . . . . . . . . . . X X X X X X X X X X X X X X X",
			". . . . . . . . . . . . . . . . X X X X X X X X X X X X X X",
			". . . . . . . . . . . . . . . . . X X X X X X X X X X X X X",
			"X X X X X X X X X X . . . . . . . X X X X X X X X X X X X X",
			"X X X X X X X X X X X . . . . . . X X X X X X X X X X X X X",
			"X X X X X X X X X X X X . . . . . X X X X X X X X X X X X X",
			"X X X X X X X X X X X X . . . . X X X X X X X X X X X X X X",
			"X X X X X X X X X X X . . . X X X X X X X X X X X X X X X X",
			"X X X X X . . . . . . . . . X X X X X X X X X X X X X X X X",
			"X X X X . . . . . . . . . X X X X X X X X X X X X X X X X X",
			"X X X X . . . . . . . . X X X X X X X X X X X X X X X X X X",
			"X X X X . . . . . X X X X X X X . . . . X X X X X X X X X X",
			"X X X X . . . . . X X X X X X . . . . . . . . . . . X X X X",
			"X X X X X . . . . . X X X X X . . . . . . . . . . . . . . .",
			"X X X X X X X . . . . . . . . . . . . . . X X X . . . . . .",
			"X X X X X X X X . . . . . . . . X X X X X X X X X X X X X X",
			"X X X X X X X X X X . . . . X X X X X X X X X X X X X X X X",
			"X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X",
			"X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// "X . X X X X X X X X X X X X X X X X X X X X X X X X X",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// "X X X X X X X X X X X X X X X X X X X X X X X X . X X",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
			// ". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		},
		Start_Point: []int{1, 0},
		End_Point:   []int{15, 29},
		Wave:        3,
	},
	2: Map_Management{
		Map_String: []string{
			"X . . . . . . . . . . . . . . . . . . . . . . . . . . X X X",
			". . . . . . . . . . . . . . . . . . . . . . . . . . . . X X",
			". . X X X X X X X X X X X X X X X X X X X X X X X X . . . X",
			"X X X . . . . . . . . . . . . . . . . . . . X X X X X . . .",
			"X . . . . . . . . . . . . . . . . . . . . . . . X X X X . .",
			"X . . . X X X X X X X X X X X X X X X X X . . . . X X X . .",
			"X . . X X X X X X X X X X X X X X X X X X X X . . X X X . .",
			"X . . . X X X X . . . . . . . . . . . X X X X . . X X X . .",
			"X X . . X X X . . . . . . . . . . . . X X X X . . X X X . .",
			"X X . . X X X . . . X X X X X X X X X X X X X . . X X X . .",
			"X X . . X X X . . X X X X X X X X X X X X X X . . X X X . .",
			"X . . . X X X . . X X X X X X X X X X X X X . . . X X X . .",
			". . . X X X . . . X X X X X X X X X X X X . . . X X X X . .",
			". . X X X X . . . X X . . . . . . . . . . . . X X X X X . .",
			". . X X X X . . . . . . . . . . . . . . . . X X X X X X . .",
			". . . X X X X . . . . . . X X X X X X X X X X X X X X X . .",
			"X . . X X X X X X X X X X X X X X X X X X X X X X X X . . .",
			"X . . . X X X X X X X X X X X X X X X X X X X X X X . . . .",
			"X X . . . . . . . . . . . . . . . . . . . . . . . . . . . X",
			"X X X . . . . . . . . . . . . . . . . . . . . . . . . . X X",
		},
		Start_Point: []int{1, 0},
		End_Point:   []int{7, 18},
		Wave:        3,
	},
}

// initialize map
func NewMap(id int) (m Map) {
	charMap := Map_Instances[id].Map_String
	m.Points = make([][]Position, len(charMap))
	m.Blocks = make(map[string]*Position, len(charMap)*2)
	for x, row := range charMap {
		cols := strings.Split(row, " ")
		m.Points[x] = make([]Position, len(cols))
		for y, view := range cols {
			m.Points[x][y] = Position{x, y}
			if view == "X" {
				m.Blocks[pointAsKey(x, y)] = &m.Points[x][y]
			}
		}
	}
	m.MaxX = len(m.Points)
	m.MaxY = len(m.Points[0])
	m.StartPoint = NewPosition(Map_Instances[id].Start_Point[0], Map_Instances[id].Start_Point[1])
	m.EndPoint = NewPosition(Map_Instances[id].End_Point[0], Map_Instances[id].End_Point[1])
	m.Wave = 3
	return m
}

/*
Obj: Find all points around the current point
*/
func (_map *Map) getAdjacentPosition(curPoint *Position) (adjacents []*Position) {
	if x, y := curPoint.X, curPoint.Y-1; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	if x, y := curPoint.X+1, curPoint.Y-1; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	if x, y := curPoint.X+1, curPoint.Y; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	if x, y := curPoint.X+1, curPoint.Y+1; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	if x, y := curPoint.X, curPoint.Y+1; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	if x, y := curPoint.X-1, curPoint.Y+1; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	if x, y := curPoint.X-1, curPoint.Y; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	if x, y := curPoint.X-1, curPoint.Y-1; x >= 0 && x < _map.MaxX && y >= 0 && y < _map.MaxY {
		adjacents = append(adjacents, &_map.Points[x][y])
	}
	return adjacents
}

// func (this *Map) PrintMap(path *SearchRoad) {
// 	fmt.Println("map's border:", this.MaxX, this.MaxY)
// 	for x := 0; x < this.MaxX; x++ {
// 		for y := 0; y < this.MaxY; y++ {
// 			if path != nil {
// 				if x == path.start.x && y == path.start.y {
// 					fmt.Print("S")
// 					goto NEXT
// 				}
// 				if x == path.end.x && y == path.end.y {
// 					fmt.Print("E")
// 					goto NEXT
// 				}
// 				for i := 0; i < len(path.TheRoad); i++ {
// 					if path.TheRoad[i].x == x && path.TheRoad[i].y == y {
// 						fmt.Print("*")
// 						goto NEXT
// 					}
// 				}
// 			}
// 			fmt.Print(this.points[x][y].view)
// 		NEXT:
// 		}
// 		fmt.Println()
// 	}
// }

// ========================================================================================
// Astar Algorithme

func NewAstarPoint(p *Position, father *_AstarPoint, end *_AstarPoint) (ap *_AstarPoint) {
	ap = &_AstarPoint{*p, father, 0, 0, 0}
	if end != nil {
		ap.calcFVal(end)
	}
	return ap
}

func (asp *_AstarPoint) calcGVal() int {
	if asp.father != nil {
		deltaX := math.Abs(float64(asp.father.position.X - asp.position.X))
		deltaY := math.Abs(float64(asp.father.position.Y - asp.position.Y))
		if deltaX == 1 && deltaY == 0 {
			asp.gVal = asp.father.gVal + 10
		} else if deltaX == 0 && deltaY == 1 {
			asp.gVal = asp.father.gVal + 10
		} else if deltaX == 1 && deltaY == 1 {
			asp.gVal = asp.father.gVal + 14
		} else {
			panic("father point is invalid!")
		}
	}
	return asp.gVal
}

func (asp *_AstarPoint) calcHVal(end *_AstarPoint) int {
	asp.hVal = int(math.Abs(float64(end.position.X-asp.position.X)) + math.Abs(float64(end.position.Y-asp.position.Y)))
	return asp.hVal
}

func (asp *_AstarPoint) calcFVal(end *_AstarPoint) int {
	asp.fVal = asp.calcGVal() + asp.calcHVal(end)
	return asp.fVal
}

func (ol OpenList) Len() int           { return len(ol) }
func (ol OpenList) Less(i, j int) bool { return ol[i].fVal < ol[j].fVal }
func (ol OpenList) Swap(i, j int)      { ol[i], ol[j] = ol[j], ol[i] }
func (ol *OpenList) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*ol = append(*ol, x.(*_AstarPoint))
}
func (ol *OpenList) Pop() interface{} {
	old := *ol
	n := len(old)
	x := old[n-1]
	*ol = old[0 : n-1]
	return x
}

func NewSearchRoad(m *Map) *SearchRoad {
	startx := m.StartPoint.X
	starty := m.StartPoint.Y
	endx := m.EndPoint.X
	endy := m.EndPoint.Y
	sr := &SearchRoad{}
	sr.theMap = m
	sr.start = *NewAstarPoint(&Position{startx, starty}, nil, nil)
	sr.end = *NewAstarPoint(&Position{endx, endy}, nil, nil)
	sr.TheRoad = make([]*_AstarPoint, 0)
	sr.openSet = make(map[string]*_AstarPoint, m.MaxX+m.MaxY)
	sr.closeLi = make(map[string]*_AstarPoint, m.MaxX+m.MaxY)
	heap.Init(&sr.openLi)
	heap.Push(&sr.openLi, &sr.start) // add start point into open list
	sr.openSet[pointAsKey(sr.start.position.X, sr.start.position.Y)] = &sr.start
	// add obstacles into close list
	for k, v := range m.Blocks {
		sr.closeLi[k] = NewAstarPoint(v, nil, nil)
	}
	return sr
}

func (search_road *SearchRoad) FindoutRoad() bool {
	for len(search_road.openLi) > 0 {
		// move nodes from openlist to closelist
		x := heap.Pop(&search_road.openLi)
		curPoint := x.(*_AstarPoint)
		delete(search_road.openSet, pointAsKey(curPoint.position.X, curPoint.position.Y))
		search_road.closeLi[pointAsKey(curPoint.position.X, curPoint.position.Y)] = curPoint
		//fmt.Println("curPoint :", curPoint.X, curPoint.Y)
		adjacs := search_road.theMap.getAdjacentPosition(&curPoint.position)
		for _, p := range adjacs {
			//fmt.Println("\t adjact :", p.x, p.y)
			theAP := NewAstarPoint(p, curPoint, &search_road.end)
			if pointAsKey(theAP.position.X, theAP.position.Y) == pointAsKey(search_road.end.position.X, search_road.end.position.Y) {
				// mark pathway after finding it out successfully
				for theAP.father != nil {
					search_road.TheRoad = append(search_road.TheRoad, theAP)
					search_road.Trace = append(search_road.Trace, theAP.position)
					theAP = theAP.father
				}
				search_road.Trace = append(search_road.Trace, search_road.theMap.StartPoint)
				search_road.Trace = reverse(search_road.Trace)
				return true
			}
			_, ok := search_road.closeLi[pointAsKey(p.X, p.Y)]
			if ok {
				continue
			}
			existAP, ok := search_road.openSet[pointAsKey(p.X, p.Y)]
			if !ok {
				heap.Push(&search_road.openLi, theAP)
				search_road.openSet[pointAsKey(theAP.position.X, theAP.position.Y)] = theAP
			} else {
				oldGVal, oldFather := existAP.gVal, existAP.father
				existAP.father = curPoint
				existAP.calcGVal()
				// If the G value of the new node is inferior to the old one, restore the old one
				if existAP.gVal > oldGVal {
					// restore father
					existAP.father = oldFather
					existAP.gVal = oldGVal
				}
			}
		}
	}
	return false
}

func reverse(s []Position) []Position {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
