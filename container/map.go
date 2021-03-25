package container

import (
	"sync"
)

type ConcurrentMap struct {
	sync.RWMutex
	items map[interface{}]interface{}
}

func (p *ConcurrentMap) init() {
	if p.items == nil {
		p.items = make(map[interface{}]interface{})
	}
}

func (p *ConcurrentMap) UnsafeHas(key interface{}) bool {
	_, ok := p.items[key]
	return ok
}

func (p *ConcurrentMap) Has(key interface{}) bool {
	p.RLock()
	defer p.RUnlock()
	_, ok := p.items[key]
	return ok
}

func (p *ConcurrentMap) UnsafeGet(key interface{}) interface{} {
	if p.items == nil {
		return nil
	} else {
		return p.items[key]
	}
}

func (p *ConcurrentMap) Get(key interface{}) interface{} {
	p.RLock()
	defer p.RUnlock()
	return p.UnsafeGet(key)
}

func (p *ConcurrentMap) UnsafeSet(key interface{}, value interface{}) {
	p.init()
	p.items[key] = value
}

func (p *ConcurrentMap) Set(key interface{}, value interface{}) {
	p.Lock()
	defer p.Unlock()
	p.UnsafeSet(key, value)
}

func (p *ConcurrentMap) TestAndSet(key interface{}, value interface{}) interface{} {
	p.Lock()
	defer p.Unlock()

	p.init()

	if v, ok := p.items[key]; ok {
		return v
	} else {
		p.items[key] = value
		return nil
	}
}

func (p *ConcurrentMap) UnsafeDel(key interface{}) {
	p.init()
	delete(p.items, key)
}

func (p *ConcurrentMap) Del(key interface{}) {
	p.Lock()
	defer p.Unlock()
	p.UnsafeDel(key)
}

func (p *ConcurrentMap) UnsafeLen() int {
	if p.items == nil {
		return 0
	} else {
		return len(p.items)
	}
}

func (p *ConcurrentMap) Len() int {
	p.RLock()
	defer p.RUnlock()
	return p.UnsafeLen()
}

func (p *ConcurrentMap) UnsafeRange(f func(interface{}, interface{})) {
	if p.items == nil {
		return
	}
	for k, v := range p.items {
		f(k, v)
	}
}

func (p *ConcurrentMap) RLockRange(f func(interface{}, interface{})) {
	p.RLock()
	defer p.RUnlock()
	p.UnsafeRange(f)
}

func (p *ConcurrentMap) LockRange(f func(interface{}, interface{})) {
	p.Lock()
	defer p.Unlock()
	p.UnsafeRange(f)
}

func (p *ConcurrentMap) UnsafeRangeSucc(f func(interface{}, interface{}) error) {
	if p.items == nil {
		return
	}
	for k, v := range p.items {
		if e:= f(k, v); e==nil{
			break
		}
	}
}
func (p *ConcurrentMap) RLockRangeSucc(f func(interface{}, interface{}) error) {
	p.RLock()
	defer p.RUnlock()
	p.UnsafeRangeSucc(f)
}

func (p *ConcurrentMap) LockRangeSucc(f func(interface{}, interface{}) error) {
	p.Lock()
	defer p.Unlock()
	p.UnsafeRangeSucc(f)
}

func (p *ConcurrentMap) UnsafeRand()(interface{},interface{}){
	for k, v := range p.items {
		return k,v
	}
	return nil,nil
}

func (p *ConcurrentMap) RLockRand()(interface{},interface{}){
	p.RLock()
	defer p.RUnlock()
	return p.UnsafeRand()
}

func (p *ConcurrentMap) LockRand()(interface{},interface{}){
	p.Lock()
	defer p.Unlock()
	return p.UnsafeRand()
}