set -e

pushd ~/go/src/github.com/foggerty/n2t

echo Building:
go build -v ./...

echo Vet:
go tool vet -shadow ./

echo Test:
go test -race ./...

echo Install:
go install ./...

popd



		
