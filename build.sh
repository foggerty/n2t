set -e

pushd ~/go/src/github.com/foggerty/n2t

echo Building:
go build -v ./...

echo Vet:
go vet ./...

echo Test:
go test ./...

echo Install:
go install ./...

popd



		
