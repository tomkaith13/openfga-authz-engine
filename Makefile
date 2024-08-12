clean:
	docker-compose down --rmi all -v
run:
	docker-compose up --build -d
restart:
	make clean && make run
stop:
	docker-compose down


