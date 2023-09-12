build:
	docker build -t forum .
run:
	docker run -dp 4000:4000 --rm --name container1 forum
stop:
	docker stop container1