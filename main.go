package main

import(
    "github.com/codegangsta/cli"
    "github.com/gin-gonic/gin"
    "digger.io/com"
    "os"
    "log"
    "io"
)
var (
    config      Config
    segmentCom  com.SegmentCom
    httpServer  HttpServer
    kbs         map[string]*KnowledgeBase
    sessions    map[string]*TriggerMessage
)

type HttpServer struct {
    engine  *gin.Engine
}

func printLog(l *LogConfig) {
    if (*l).Mode != "print" {
        f, err := os.OpenFile((*l).Path, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
        if err != nil {
            log.Printf("main.PrintLog->", "error opening file: %v", err)
        }
        defer f.Close()
        log.SetOutput(io.MultiWriter(f, os.Stdout))
    }
}

func Init(filename string) {
    // load config
    m := config.Load(filename)
    // log
    printLog(&m.Log)
    // init segment com
    segmentCom.InitConn(&m.Segment)
    // init list
    sessions = make(map[string]*TriggerMessage)
    // init knowledge base
    kbs = make(map[string]*KnowledgeBase)
    for _, v := range m.Knowledge {
        kbs[v] = (new(KnowledgeBase)).Init(v)
    }
    // run http server
    httpServer.Start(&m.Http)
}

func (h *HttpServer) Start(c *HttpConfig) {
    // choose mode
    if (*c).Mode == "debug" {
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }
    // init engine
    h.engine = gin.Default()
    // regsiter router
    h.router()
    // run
    h.engine.Run((*c).Host)
}

func (h *HttpServer) router() {
    // enlation router
    h.engine.POST("/engine/:uid/:kid", h.enginHandle)
}

func (h *HttpServer) enginHandle(c *gin.Context) {
    // header
    c.Header("Access-Control-Allow-Origin", "*")
    // action
    uid := c.Param("uid")
    kid := c.Param("kid")
    msg:= c.PostForm("text")
    if _, ok := sessions[uid]; !ok {
        if _, ok := kbs[kid]; ok {
            sessions[uid] = (new (TriggerMessage)).Init(kbs[kid])
        }
    }
    if _, ok := sessions[uid]; ok {
        q, s, _, _ := (sessions[uid]).Trigger(msg)
        c.JSON(200, gin.H{"ok": true,  "question": q.Text, "story": s.Text})
    } else {
        c.JSON(404, gin.H{"ok": false, "error": "no found the model"})
    }
}

func main() {
    app := cli.NewApp()
    app.Name = "Digger"
    app.Usage = "http"
    app.Version = "0.50"
    app.Action = func(c *cli.Context) error {
        Init("setting.yaml")
        return nil
    }
    app.Run(os.Args)
}
