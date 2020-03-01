VER = example

docker:
	go build cmd.go
	sudo docker build --tag rbxorkt12/appcd:$(VER) .
	sudo docker push rbxorkt12/appcd:$(VER)
