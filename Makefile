TAG=`git describe --abbrev=0 --tags`
HEAD=`git log -n 1 --pretty="%h"`
DATE=`date +%FT%T%z`

LDFLAGS = -ldflags "-X main.version=${TAG} -X main.date=${DATE} -X main.head=${HEAD}"

build:
	go build ${LDFLAGS}
