package fixedpool

import (
	"container/list"
)

type Pool struct{
	l *list.List
}

func New()*Pool{
	return &Pool{
		l: list.New(),
	}
}

func (p *Pool)Put(b interface{}){
	p.l.PushBack(b)
}

func (p *Pool)Pop()interface{}{
	//fmt.Println("Buffer pool available ", p.l.Len())
	f := p.l.Front()
	if f == nil{
		return nil
	}
	p.l.Remove(f)
	return f.Value
}
