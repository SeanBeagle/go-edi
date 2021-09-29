all:
	@echo "Try targets: 'build', 'run', 'stop'"

build:
	docker build . --tag go-edi
run:
	docker container run -itd -p 9000:8080 --name go_edi_local go-edi
stop:
	docker stop go_edi_local
