all: deploy


clean:
	rm -rf "$(GOPATH)/pkg/darwin_amd64/fues3"

cross:
	env GOOS=linux GOARCH=arm GOARM=7 godep go build

deploy:
	godep go build

test:
	godep go test -cover ./tests


.PHONY: *
