BINDIR=${CURDIR}/bin

bindir:
	echo "build walletStatus"
	mkdir -p ${BINDIR}

build: bindir
	go build -o ${BINDIR}/walletStatus ./cmd

run: build
	CONFIG_FILE=./configs/values.yaml ${BINDIR}/walletStatus