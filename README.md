# fixedpool
Simple fixed pool of objects

Note: fixedpool is _not_ thread safe, this is so it can be used 
in a non-threaded environment without extra overhead

Thread segregated fixed pool is thread safe and will lock at the pool level
for push and pop. This should mean that if only a single thread uses a given pool ID, there's no contention

Or if each popping thread uses it's own pool ID, but the releasing thread is different, the contention should 
only be between popping thread and releasing thread.

But, it can also be used to reduce contention by assigning about 4 times more
pools than there are threads. Each thread then randomly (or via a scheme) selects which pool ID
to pop from. 

## Usage

### Create a pool item that obeys the interface
For example:

```
type PoolItem struct{
	counter int
	poolid int
}

/// Implement the two required methods SetPoolId and GetPoolId
func (p *PoolItem) SetPoolId(poolid int) {
	p.poolid = poolid
}

func (p *PoolItem) GetPoolId() int {
	return p.poolid
}

```

### Fill the pool at startup

```
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

```

### Pop an item from the pool

```
	pool := NewThreadSegregatedFixedPool(numpools)

	//// fill the pool

	item := pool.Pop(poolid)
```

### Put the item back

```
	pool.Put(item)
```
