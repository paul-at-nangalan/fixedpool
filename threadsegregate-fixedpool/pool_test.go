package threadsegregate_fixedpool

import (
	"sync"
	"testing"
)

type PoolItem struct{
	counter int
	poolid int
}

func (p *PoolItem) SetPoolId(poolid int) {
	p.poolid = poolid
}

func (p *PoolItem) GetPoolId() int {
	return p.poolid
}

func fillPool(pool *ThreadSegregatedFixedPool, numpools, numperpool int){
	for i := 0; i < numpools; i++{
		for x := 0; x < numperpool; x++{
			item := PoolItem{
				counter: x,
			}
			pool.PutById(&item, i)
		}
	}
}

func TestThreadSegregatedFixedPool_PutById(t *testing.T) {
	numpools := 4
	pool := NewThreadSegregatedFixedPool(numpools)
	numperpool := 10
	fillPool(pool, numpools, numperpool)

	///for each pool make sure we can get exactly that amount back
	for i := 0; i < numpools; i++ {
		for x := 0; x < numperpool; x++ {
			item := pool.Pop(i)
			if item == nil{
				t.Error("Failed to get all items, pool ", i, " count ", x)
			}
			if item.GetPoolId() != i{
				t.Error("Item pulled from wrong pool ", item.GetPoolId(), i)
			}
		}
		///try to get 1 extra
		item := pool.Pop(i)
		if item != nil{
			t.Error("Managed to pop more than ", numperpool, " items")
		}
	}

}

func TestThreadSegregatedFixedPool_Put(t *testing.T) {
	numpools := 4
	pool := NewThreadSegregatedFixedPool(numpools)
	numperpool := 100
	fillPool(pool, numpools, numperpool)

	///this should keep running til the end
	for i := 0; i < 1000 * numpools; i++{
		poolid := i % numpools
		item := pool.Pop(poolid)
		if item.GetPoolId() != poolid{
			t.Error("Popped item from pool ", poolid, " has different pool id ", item.GetPoolId())
		}
		pool.Put(item)
	}
	popped := make([]SegmentedPoolItem, 21)
	store := popped[:0]
	for i := 0; i < 1000 * numpools; i++{
		poolid := i % numpools
		item := pool.Pop(poolid)
		if item == nil{
			t.Error("Popped item is nil")
			t.FailNow()
		}
		if item.GetPoolId() != poolid{
			t.Error("Popped item from pool ", poolid, " has different pool id ", item.GetPoolId())
		}
		if len(store) < len(popped){
			store = append(store, item)
		}else{
			///put them all back
			for _, item := range store{
				pool.Put(item)
			}
			store = popped[:0]
		}
	}
}

var wait sync.WaitGroup

func runThread(pool *ThreadSegregatedFixedPool, strategy, numpools int64, numruns int64, t *testing.T){
	defer wait.Add(-1)
	for i := int64(0); i < 1000 * numruns; i++ {
		poolid := (i * strategy) % numpools
		item := pool.Pop(int(poolid))

		if item.GetPoolId() != int(poolid){
			t.Error("Incorrect pool id ", poolid, item.GetPoolId())
			t.FailNow()
		}
		pool.Put(item)
	}
}

func TestThreadSegregatedFixedPool_Threading(t *testing.T) {

	numpools := 4
	pool := NewThreadSegregatedFixedPool(numpools)
	numperpool := 100
	fillPool(pool, numpools, numperpool)

	wait.Add(4)
	go runThread(pool, 1, int64(numpools), 100000, t)
	go runThread(pool, 2, int64(numpools), 100000, t)
	go runThread(pool, 3, int64(numpools), 100000, t)
	go runThread(pool, 4, int64(numpools), 100000, t)
	wait.Wait()
}