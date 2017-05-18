# This is how we want to name the binary output
OUTPUT=netnode_probe_exe
# These are the values we want to pass for Version and BuildTime
GITTAG=$(tag)
BUILD_TIME=`date +%FT%T%z`
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X netnode_probe.GitTag=${GITTAG} -X netnode_probe.BuildTime=${BUILD_TIME}"
all:
	go build ${LDFLAGS} -o ${BINARY}  netnode_probe_exe

