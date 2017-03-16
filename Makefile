export GOPATH=${CURDIR}

app:
	go install cmds/api-toolkit

prepare:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

install: assets
		# go install cmds/ssh-execute
		# go install cmds/ssh-http

dev: devassets
		go install cmds/ssh-execute
		go install cmds/ssh-http

assets: src/util/assets/
		PATH=${PATH}:${CURDIR}/bin go-bindata-assetfs -pkg util $<...
		mv bindata_assetfs.go src/util/bindata_assetfs.go

devassets: src/util/assets/
		PATH=${PATH}:${CURDIR}/bin go-bindata-assetfs -debug -pkg util $<...
		mv bindata_assetfs.go src/util/bindata_assetfs.go

install-tools:
	go install github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs
	go install github.com/jteeuwen/go-bindata/go-bindata
