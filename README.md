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
