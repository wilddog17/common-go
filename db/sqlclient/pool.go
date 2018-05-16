package sqlclient

const POOL_SIZE = 256
const BUF_SIZE = 256

type Record = map[string]interface{}
type Records = []*Record
type PoolChan struct {
	buffers chan Records
}

func (p *PoolChan) Get() Records {
	select {
		case b := <-p.buffers:
			return b
			default:
	}

	return make(Records, 0, BUF_SIZE)
}

func (p *PoolChan) Release(b Records) {
	if(len(p.buffers) == POOL_SIZE || cap(b) > BUF_SIZE) { // gc
		return
	}

	select {
		case p.buffers <- b[:0]:
		default:
	}
}

var Pool = &PoolChan{
	buffers: make(chan Records, POOL_SIZE+1),
}