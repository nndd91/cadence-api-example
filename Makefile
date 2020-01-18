default: bins

httpserver:
	go build -i -o bins/httpserver app/httpserver/main.go

worker:
	go build -i -o bins/worker app/worker/main.go

bins: httpserver worker

clean:
	rm -rf bins

