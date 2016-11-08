package com

import(
    "gopkg.in/redis.v4"
    "log"
)

type RedisCom struct {
    Conn    *redis.Client
}

type RedisConfig struct {
    Host        string
    Password    string
    DB          int
}

func (r *RedisCom) InitConn(config *RedisConfig) (*redis.Client){
    r.Conn = redis.NewClient(&redis.Options{
         Addr: config.Host,
         Password: config.Password,
         DB: config.DB,
    })
    pong, err := r.Conn.Ping().Result()
    if err != nil {
        log.Println(pong, err)
    }
    return r.Conn
}

func (r *RedisCom) GetConn() (*redis.Client) {
    return r.Conn
}

