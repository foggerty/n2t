function run() {
		if [[ $? -eq 0 ]]
		then
				$1 $2 $3
		else
				echo "Skipping $2"
		fi
}

pushd ~/go/src/github.com/foggerty/n2t

echo Building:
go build -v ./...
echo Vet:
run go vet ./...
echo Test:
run go test ./...
echo Install:
run go install ./...

popd



		
