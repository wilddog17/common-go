package asclient

import (
    "os"
    "log"

    . "github.com/aerospike/aerospike-client-go"
    "github.daumkakao.com/live-core/scaffolding-go.git/logger"
)

type Records = []BinMap
type AsClient struct {
    client *Client
    luaPath string
    luaPkgName string
}

var AS AsClient

func (as *AsClient) Register(filename string) error {
    regTask, err := as.client.RegisterUDFFromFile(nil, as.luaPath+filename, filename, LUA)
    if err != nil {
        return err
    }

    if errReg := <- regTask.OnComplete(); errReg != nil {
        return errReg
    }

    return nil
}

func (as *AsClient) Deregister(filename string) error {
    rmTask, err := as.client.RemoveUDF(nil, filename)
    if err != nil {
        return err
    }

    if errRm := <- rmTask.OnComplete(); errRm != nil {
        return errRm
    }

    return nil
}

func (as *AsClient) PutBins(ns string, set string, k interface{}, bins []*Bin) error {
    key, kErr := NewKey(ns, set, k)
    if kErr != nil {
        return kErr
    }

    err := as.client.PutBins(nil, key, bins...)
    return err
}

func (as *AsClient) Put(ns string, set string, k interface{}, record BinMap) error {
    key, kErr := NewKey(ns, set, k)
    if kErr != nil {
        return kErr
    }

    err := as.client.Put(nil, key, record)
    return err
}

func (as *AsClient) PutTTL(ns string, set string, k interface{}, t uint32, record BinMap) error {
    key, kErr := NewKey(ns, set, k)
    if kErr != nil {
        return kErr
    }

    err := as.client.Put(NewWritePolicy(0, t), key, record)
    return err
}

func (as *AsClient) Remove(ns string, set string, k interface{}) error {
    key, kErr := NewKey(ns, set, k)
    if kErr != nil {
        return kErr
    }

    _, err := as.client.Delete(nil, key)
    return err
}

func (as *AsClient) query(p *QueryPolicy, stm *Statement) (Records, error) {
    recordset, err := as.client.Query(p, stm)
    if err != nil {
        logger.Error("[asclient] " + err.Error())
        return nil, err
    }

    l := 0
    buf := Pool.Get()
    //buf := make(Records, 0, size)
    for res := range recordset.Results() {
        if res.Err != nil {
            logger.Error("[asclient] result " + res.Err.Error())
            return nil, res.Err
        }

        l += 1
        buf = append(buf, res.Record.Bins)
        if l == cap(buf) {
            nBuf := make(Records, l, l*2)
            copy(nBuf, buf)
            Pool.Release(buf)
            buf = nBuf
        }
    }

    return buf, nil
}

func (as *AsClient) QueryEqual(ns string, s string, bn string, q interface{}, binNames ...string) (Records, error) {
    stm := NewStatement(ns, s, binNames...)
    stm.Addfilter(NewEqualFilter(bn, q))

    //p := NewQueryPolicy()
    //p.RecordQueueSize = 128
    //p.Timeout = 500 * time.Millisecond
    //p.SocketTimeout = 500 * time.Millisecond
    //p.MaxRetries = 1

    return as.query(nil, stm)
}

func (as *AsClient) QueryRange(ns string, s string, bn string, b int64, e int64, binNames ...string) (Records, error) {
    stm := NewStatement(ns, s, binNames...)
    stm.Addfilter(NewRangeFilter(bn, b, e))

    //p := NewQueryPolicy()
    //p.RecordQueueSize = 128
    //p.Timeout = 500 * time.Millisecond
    //p.SocketTimeout = 500 * time.Millisecond
    //p.MaxRetries = 1

    return as.query(nil, stm)
}

func GetClient() *AsClient {
    return &AS
}

func Run(config AerospikeConfig) {
    rp, _ := os.Getwd()

    clientPolicy := NewClientPolicy()
    clientPolicy.ConnectionQueueSize = 512 //256
    clientPolicy.LimitConnectionsToQueueSize = true
    //clientPolicy.Timeout = 500 * time.Millisecond

    hosts := make([]*Host, 0)
    for _, h := range config.Hosts {
        hosts = append(hosts, NewHost(h.Host, h.Port))
    }

    client, err := NewClientWithPolicyAndHost(clientPolicy, hosts...)
    if err != nil {
        log.Fatalln(err.Error())
    }

    AS = AsClient{
        client: client,
        luaPath: rp + "/db/asclient/udf/",
        luaPkgName: "default",
    }

    SetLuaPath(AS.luaPath)
    logger.Debug("[asclient] init done")
}