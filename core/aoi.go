package core

import (
	"errors"
	"fmt"
	"github.com/viphxin/xingo/logger"
	_ "time"
	"sync"
	"time"
	"github.com/viphxin/xingo/utils"
	"github.com/viphxin/xingo/timer"
)

/*
格子
*/
type Grid struct {
	ID   int32
	minX int32
	minY int32
	maxX int32
	maxY int32
	pids map[int32]bool
	sync.RWMutex
}

func (this *Grid)Add(pid int32){
	this.Lock()
	defer this.Unlock()
	this.pids[pid] = true
}

func (this *Grid)Remove(pid int32){
	this.Lock()
	defer this.Unlock()
	delete(this.pids, pid)
}

func (this *Grid)GetPids() []int32{
	this.RLock()
	defer this.RUnlock()
	pids := make([]int32, 0)
	for pid, _ :=range this.pids{
		pids = append(pids, pid)
	}
	return pids
}

/*
aoi
*/
type AOIMgr struct {
	minX int32
	minY int32
	maxX int32
	maxY int32
	lenX int32
	lenY int32
	grids map[int32]*Grid
}
func debugPrintFunc(params ...interface{}){
	logger.Info("grids safe==================")
	for gridid, grid := range WorldMgrObj.AoiObj1.GetGrids() {
		if len(grid.GetPids()) > 0 {
			logger.Info(fmt.Sprintf("grid: %d. players: %d", gridid, len(grid.GetPids())))
		}
	}

	utils.GlobalObject.GetSafeTimer().CreateTimer(3000, debugPrintFunc, []interface{}{new(timer.ParamNull)})

}

func NewAOIMgr(minX int32, maxX int32, minY int32, maxY int32, lenX int32, lenY int32) *AOIMgr{
	AOIObj := &AOIMgr{
		minX: minX,
		minY: minY,
		maxX: maxX,
		maxY: maxY,
		lenX: lenX,
		lenY: lenY,
		grids: make(map[int32]*Grid, 0),
	}
	AOIObj.InitGrid()
	//debug

	if utils.GlobalObject.GetSafeTimer() != nil{
		utils.GlobalObject.GetSafeTimer().CreateTimer(3000, debugPrintFunc, nil)
	}else{
		go func() {
			for {
				logger.Info("grids==================")
				for gridid, grid := range AOIObj.GetGrids() {
					if len(grid.GetPids()) > 0 {
						logger.Info(fmt.Sprintf("grid: %d. players: %d", gridid, len(grid.GetPids())))
					}
				}
				time.Sleep(3 * time.Second)
			}
		}()
	}
	return AOIObj
}

func (this *AOIMgr)GetGrids() map[int32]*Grid{
	return this.grids
}

func (this *AOIMgr)GetPiceX() int32{
	return (this.maxX - this.minX)/this.lenX
}

func (this * AOIMgr)GetPiceY() int32{
	return (this.maxY - this.minY)/this.lenY
}

func (this *AOIMgr)InitGrid(){
	//x grid长度
	piceX := this.GetPiceX()
	//y grid长度
	piceY := this.GetPiceY()
	for _y := int32(0); _y < this.lenY; _y++{
		for _x := int32(0); _x < this.lenX; _x++{
			grid := &Grid{
				ID:   this.lenX*_y + _x,
				minX: piceX*_x + this.minX,
				minY: piceY*_y + this.minY,
				maxX: piceX*(_x + 1) + this.minX,
				maxY: piceY*(_y + 1) + this.minY,
				pids: make(map[int32]bool, 0),
			}
			this.grids[grid.ID] = grid
		}
	}
	//logger.Info(this.grids)
}

func (this *AOIMgr)GetGridIDByPos(px float32, py float32) int32{
	X := int32(px)
	Y := int32(py)
	_x := (X - this.minX)/this.GetPiceX()
	if (X - this.minX)%this.GetPiceX() > 0{
		_x += 1
	}

	_y := (Y - this.minY)/this.GetPiceY()
	if (Y - this.minY)%this.GetPiceY() > 0{
		_y += 1
	}
	return this.lenX*_y + _x
}

func (this *AOIMgr)GetGridID(p *Player) int32{
	return this.GetGridIDByPos(p.X, p.Z)
}

func (this *AOIMgr)Add2AOI(p *Player) error{
	grid, ok := this.grids[this.GetGridID(p)]
	if ok {
		grid.Add(p.Pid)
	}else{
		return errors.New(fmt.Sprintf("Add2AOI 坐标有误: (%f, %f)", p.X, p.Z))
	}
	return nil
}

func (this *AOIMgr)LeaveAOI(p *Player) error{
	grid, ok := this.grids[this.GetGridID(p)]
	if ok {
		grid.Remove(p.Pid)
	}else{
		return errors.New(fmt.Sprintf("LeaveAOI 坐标有误: (%f, %f)", p.X, p.Z))
	}
	return nil
}

func (this *AOIMgr)LeaveAOIFromGrid(p *Player, gridId int32) error{
	grid, ok := this.grids[gridId]
	if ok {
		grid.Remove(p.Pid)
	}else{
		return errors.New(fmt.Sprintf("LeaveAOI 坐标有误: (%f, %f)", p.X, p.Z))
	}
	return nil
}

/*
通过gridId获取九宫格
*/
func (this *AOIMgr)GetSurroundingByGridId(gridId int32) ([]*Grid, error){
	grids := make([]*Grid, 0)
	grid, ok := this.grids[gridId]
	if ok {
		pos := make([]int32, 0)
		grids = append(grids, grid)
		pos = append(pos, gridId)
		//找中间两个
		if gridId % this.lenX > 0{
			//有左边的grid
			grids = append(grids, this.grids[gridId - 1])
			pos = append(pos, gridId - 1)
		}
		if gridId % this.lenX < this.lenX - 1{
			//有右边的grid
			grids = append(grids, this.grids[gridId + 1])
			pos = append(pos, gridId + 1)
		}

		for _, p := range pos{
			if p / this.lenX > 0 || (p > 0 && p / this.lenX == 0){
				//有下面的格子
				grids = append(grids, this.grids[p - this.lenX])
			}
			if p / this.lenX <= this.lenX - 1{
				//有上面的格子
				grids = append(grids, this.grids[p + this.lenX])
			}
		}
	}else{
		return nil, errors.New(fmt.Sprintf("GetSurroundingByGridId GridID有误: (%d, )", gridId))
	}
	return grids, nil
}

/*
获取九宫格
*/
func (this *AOIMgr)GetSurroundingGrids(px float32, py float32) ([]*Grid, error){
	//当前玩家所在grid
	gridId := this.GetGridIDByPos(px, py)
	grids, err := this.GetSurroundingByGridId(gridId)
	if err != nil{
		logger.Error(err)
		return nil, errors.New(fmt.Sprintf("GetSurroundingGrids 坐标有误: (%f, %f)", px, py))
	}
	return grids, nil
}

func (this *AOIMgr)GetSurroundingPids(p *Player) ([]int32, error){
	pids := make([]int32, 0)
	grids, err := this.GetSurroundingGrids(p.X, p.Z)

	if err == nil{
		for _, grid := range grids{
			pids = append(pids, grid.GetPids()...)
		}
	}else{
		return nil, err
	}

	return pids, nil
}

