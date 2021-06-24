package threadsegregate_fixedpool

import (
	"github.com/paul-at-nangalan/fixedpool/fixedpool"
	"sync"
)

type SegmentedPoolItem interface {
	SetPoolId(int)
	GetPoolId()(int)
}

type Pool struct {
	pool *fixedpool.Pool
	mut sync.Mutex
}

type ThreadSegregatedFixedPool struct{
	pools []Pool
}

////each pool needs to be filled at startup ... as we don't know the object,
/// it must be done by the calling thread
func NewThreadSegregatedFixedPool(numthreads int)*ThreadSegregatedFixedPool{
	pools := make([]Pool, numthreads)
	for i := 0; i < numthreads; i++{
		pool := Pool{
			pool: fixedpool.New(),
		}
		pools[i] = pool
	}
	return &ThreadSegregatedFixedPool{
		pools: pools,
	}
}
///Provided each pool is used by only a single thread, there should be no contention
/// if multiple threads us it, there will be contention, but ypu can use it to
/// reduce contention
///

///Use this to initialise the pool with objects
func (p *ThreadSegregatedFixedPool)PutById(obj SegmentedPoolItem, poolid int){
	p.pools[poolid].mut.Lock()
	defer p.pools[poolid].mut.Unlock()
	p.pools[poolid].pool.Put(obj)
}

func (p *ThreadSegregatedFixedPool)Put(obj SegmentedPoolItem){
	id := obj.GetPoolId()
	p.PutById(obj, id)
}

func (p *ThreadSegregatedFixedPool)Pop(poolid int)SegmentedPoolItem{
	p.pools[poolid].mut.Lock()
	defer p.pools[poolid].mut.Unlock()
	popped := p.pools[poolid].pool.Pop()
	if popped == nil{
		return nil
	}
	popped.(SegmentedPoolItem).SetPoolId(poolid)
	return popped.(SegmentedPoolItem)
}
