APP = digger
DB = github.com/go-sql-driver/mysql
YAML = github.com/go-yaml/yaml
REDISV4 = gopkg.in/redis.v4
GIN = github.com/gin-gonic/gin
CLI = github.com/codegangsta/cli
SET = github.com/deckarep/golang-set

export http_proxy=localhost:9743
export https_proxy=localhost:9743
export GOPATH = ${PWD}

develop:
	go get ${DB}
	go get ${YAML}
	go get ${REDISV4}
	go get ${GIN}
	go get ${CLI}
	go get ${SET}
	go build -o ${APP} -ldflags '-s -w'
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux_${APP} -ldflags '-s -w'
	
	#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows_${APP} -ldflags '-s -w'
build:
	go get ${DB}
	go get ${YAML}
	go get ${REDISV4}
	go get ${GIN}
	go get ${CLI}
	go get ${SET}
	go build -o ${APP} -ldflags '-s -w'
	goupx -s=true -u ${APP}

run:
	@go run *.go

clean:
	@rm bin/windows_${APP}.exe bin/mac_${APP} bin/linux_${APP}
	go clean
