.ONESHELL:
.SILENT:

EXE := triple-star.exe
tidy:
	go mod tidy

run:build
	${EXE}

build:
	go build -o ${EXE} .

clean:
	go clean -x

test:
	go test -v ./...

abigen:
	cd service/contract
	abigen -abi storage.abi --pkg contract --type Storage --out storage.go
	abigen -abi payment.abi --pkg contract --type Payment --out payment.go
