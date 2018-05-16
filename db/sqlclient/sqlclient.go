package sqlclient

import (
    "log"
    "database/sql"

    "github.daumkakao.com/live-core/scaffolding-go.git/logger"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Close() {
    db.Close()
}

//func QueryRaw(q string) (*sql.Rows, error) {
//    rows, err := db.Query(q)
//    return rows, err
//}

func Query(q string) (Records, error) {
    rows, err := db.Query(q)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }

    l := 0
    buf := Pool.Get()
    v := make([]interface{}, len(columns))
    p := make([]interface{}, len(columns)) // clear v
    for rows.Next() {
        for i := range v {
            v[i] = &p[i]
        }

        if err := rows.Scan(v...); err != nil {
            return nil, err
        }

        e := Record{}
        for i, name := range columns {
            b, ok := p[i].([]byte)
            if (ok) {
                e[name] = string(b)
            } else {
                e[name] = p[i]
            }
        }

        l += 1
        buf = append(buf, &e)
        if l == cap(buf) {
            nBuf := make(Records, l, l*2)
            copy(nBuf, buf)
            Pool.Release(buf)
            buf = nBuf
        }
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return buf, nil
}

func Run() {
    var err error
    db, err = sql.Open("mysql", "live:Street@tcp(live-dev.s2.krane.9rum.cc:3306)/admin-tool?charset=utf8")
    if err != nil {
        log.Fatalln(err.Error())
    }

    logger.Debug("[sqlclient] init done")
}
