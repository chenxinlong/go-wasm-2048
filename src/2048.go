package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

type PanelMatrix [4][4]int

var (
	panel = &PanelMatrix{}
)

func main()  {
	rand.Seed(time.Now().UnixNano())

	d := make(chan struct{}, 0)
	registerInvokes()
	<-d
}

func registerInvokes()  {
	js.Global().Set("wasmStart", js.FuncOf(_start))
	js.Global().Set("wasmMove", js.FuncOf(_move))
}

func _start(this js.Value, args []js.Value) interface{} {
	panel = &PanelMatrix{}
	panel.genRandomCell()
	panel.genRandomCell()

	return panel.toJs()
}

func _move(this js.Value, args []js.Value) interface{} {
	before, _ := json.Marshal(panel)
	panel.move(args[0].Int())
	after, _ := json.Marshal(panel)

	// 如果移动前后 panel 无任何变化，或者移动后 panel 已无空位则不产生新值
	if string(before) != string(after) && panel.hasAvailPos() {
		panel.genRandomCell()
	}
	return panel.toJs()
}

// 移动
func (pl *PanelMatrix) move(keyCode int) {
	for i:=0;i<4;i++ {
		// move horizon
		if keyCode == 37 || keyCode == 39 {
			panel[i][0],panel[i][1],panel[i][2],panel[i][3] = merge(panel[i][0], panel[i][1], panel[i][2], panel[i][3], keyCode)

		// move vertical
		} else {
			panel[0][i], panel[1][i], panel[2][i], panel[3][i] = merge(panel[0][i],panel[1][i],panel[2][i], panel[3][i], keyCode)
		}
	}
}

// 场面上产生一个新值
func (pl *PanelMatrix) genRandomCell()  {
	x, y := genRandPosition()
	val := genRandValue()
	fmt.Printf("New cell, pos:[%v, %v] val:%v\n",x,y, val)
	pl[x][y] = val
}

// js 只接收 go 的 []interface{}
func (pl *PanelMatrix) toJs() js.Value {
	ret := []interface{}{}
	for _, row := range pl {
		fmt.Println(row)
		for _, col := range row {
			ret = append(ret, col)
		}
	}
	fmt.Println("=======")
	return js.ValueOf(ret)
}

// {2,4} 中随机出一个值 (加权随机，提高2的权重)
func genRandValue() int {
	pool := []int{2,4}
	totalWeight := 100
	randNum := rand.Intn(totalWeight)

	// Weight 4 = 35
	// Weight 2 = 65
	idx := 0
	if randNum < 35 {
		idx = 1
	}
	return pool[idx]
}

// 场上剩余可用 position 中随机出一个
func genRandPosition() (x, y int) {
	if panel.hasAvailPos() == false {
		return 0, 0
	}

	x, y = rand.Intn(4), rand.Intn(4)
	for panel[x][y] != 0 {
		x, y = rand.Intn(4), rand.Intn(4)
	}

	return
}

// panel 中是否有未填的格子
func (pl *PanelMatrix)hasAvailPos() bool {
	for _, row := range pl {
		if v1,_,_,_ := squash(row[0],row[1],row[2],row[3], 39); v1==0 {
			return true
		}
	}
	return false
}

// 把该行(or该列)的所有0 squash 到居左(or居尾)
func squash(v1,v2,v3,v4, keycode int) (int, int, int, int) {
	a := []int{v1,v2,v3,v4}
	n := []int{}
	for _, v := range a {
		if v != 0 {
			n = append(n, v)
		}
	}

	lack := 4 - len(n)
	for i:=0;i<lack;i++{
		if keycode <= 38 {
			// squash to head
			n = append(n, 0)
		} else {
			// squash to tail
			n = append([]int{0}, n...)
		}
	}

	return n[0], n[1], n[2], n[3]
}

// 合并计算
func merge(v1,v2,v3,v4, keycode int) (nv1,nv2,nv3,nv4 int) {
	v1,v2,v3,v4 = squash(v1,v2,v3,v4, keycode)
	a := []int{v1,v2,v3,v4}

	if keycode == 39 || keycode == 40 {
		for i := 3; i > 0; i-- {
			if a[i] == a[i-1] {
				a[i], a[i-1] = a[i]+a[i-1], 0
			}
		}
	}

	if keycode == 37 || keycode == 38 {
		for i:=0; i<3; i++ {
			if a[i] == a[i+1] {
				a[i], a[i+1] = a[i]+a[i+1], 0
			}
		}
	}

	return squash(a[0],a[1],a[2],a[3], keycode)
}
