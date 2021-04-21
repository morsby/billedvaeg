watch:
	reflex -s -r ".*\.go|html|ts|css$\" -R "tmp|dist" make run
run:
	go run cmd/main.go
