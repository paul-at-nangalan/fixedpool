package threadsegregate_fixedpool

import "testing"

type PoolItem struct{
	poolid int
}

func (p *PoolItem) SetPoolId(poolid int) {
	p.poolid = poolid
}

func (p *PoolItem) GetPoolId() int {
	return p.poolid
}

func TestThreadSegregatedFixedPool_PutById(t *testing.T) {

}

func TestThreadSegregatedFixedPool_Put(t *testing.T) {

}

func TestThreadSegregatedFixedPool_Pop(t *testing.T) {

}
