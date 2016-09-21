function run() {
		if [ $? -eq 0 ]
		then
				$1 $2 $3
		else
				echo "Skipping $2"
		fi
}

pushd ~/go/src/github.com/foggerty/n2t

go build -v ./...
run go vet ./...
run go test ./...
run go install ./...

popd



		
