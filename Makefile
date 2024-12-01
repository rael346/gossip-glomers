debug:
	sudo ./maelstrom serve

echo:
	go build -o build/echo cmd/echo/main.go
	./maelstrom test \
		-w echo \
		--bin ./build/echo \
		--node-count 1 \
		--time-limit 10 \

idgen:
	go build -o build/idgen cmd/idgen/main.go
	./maelstrom test \
		-w unique-ids \
		--bin ./build/idgen \
		--time-limit 30 \
		--rate 1000 \
		--node-count 3 \
		--availability total \
		--nemesis partition \

broadcast-a:
	go build -o build/broadcast cmd/broadcast/main.go
	./maelstrom test \
		-w broadcast \
		--bin ./build/broadcast \
		--node-count 1 \
		--time-limit 20 \
		--rate 10 \

broadcast-b:
	go build -o build/broadcast cmd/broadcast/main.go
	./maelstrom test \
		-w broadcast \
		--bin ./build/broadcast \
		--node-count 5 \
		--time-limit 20 \
		--rate 10 \

broadcast-c:
	go build -o build/broadcast cmd/broadcast/main.go
	./maelstrom test \
		-w broadcast \
		--bin ./build/broadcast \
		--node-count 5 \
		--time-limit 20 \
		--rate 10 \
		--nemesis partition \

broadcast-d:
	go build -o build/broadcast cmd/broadcast/main.go
	./maelstrom test \
		-w broadcast \
		--bin ./build/broadcast \
		--node-count 25 \
		--time-limit 20 \
		--rate 100 \
		--latency 100 \
		--topology total \

broadcast-e:
	go build -o build/broadcast cmd/broadcast/main.go
	./maelstrom test \
		-w broadcast \
		--bin ./build/broadcast \
		--node-count 25 \
		--time-limit 20 \
		--rate 100 \
		--latency 100 \
		--topology tree4 \
