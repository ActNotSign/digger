package com

import(
    "github.com/iwanbk/gobeanstalk"
)

type QueueCom struct {
    Conn    *gobeanstalk.Conn
}

type QueueConfig struct {
    Host    string
    Tube    string
}

func (q *QueueCom) InitConn(config *QueueConfig) (*gobeanstalk.Conn) {
    var err error
    q.Conn, err = gobeanstalk.Dial(config.Host)
    if err != nil {
        panic(err)
    }
    q.Conn.Watch(config.Tube)
    return q.Conn
}
