package utils

const POOL_SIZE = 256
const BUF_SIZE = 64
type PoolChan struct {
	buffers chan []byte
}

func (p *PoolChan) Get() (b []byte) {
	select {
	case b = <-p.buffers:
	default:
		buff := make([]byte, BUF_SIZE)
		return buff
	}

	return
}

func (p *PoolChan) Release(b []byte) {
	if(len(p.buffers) == POOL_SIZE || cap(b) > BUF_SIZE) { // gc
		return
	}

	select {
	case p.buffers <- b:
	default:
	}
}

var Pool = &PoolChan{
	buffers: make(chan []byte, POOL_SIZE + 1),
}