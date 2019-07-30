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

echo Running comparison test:
~/go/bin/assembler -in ./Pong.asm -out test.hack
diff -y --suppress-common-lines test.hack Pong-Reference.hack 

popd
