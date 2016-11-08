package main

import(
    "log"
    "io/ioutil"
    "github.com/go-yaml/yaml"
    "digger.io/com"
)

type ConfigInformation struct {
    Segment     com.SegmentConfig
    Http        HttpConfig
    Log         LogConfig
    Knowledge   []string
}

type HttpConfig struct {
    Host    string
    Mode    string
}

type LogConfig struct {
    Path    string
    Mode    string
}

type Config struct {
    path    string
    params  *ConfigInformation
}

func (c *Config) Load(filename string) (*ConfigInformation){
    dat, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Panic(err)
    } else {
        c.params = new(ConfigInformation)
        err = yaml.Unmarshal([]byte(dat), c.params)
        if err != nil {
            log.Panic(err)
        }
    }
    return c.params
}
