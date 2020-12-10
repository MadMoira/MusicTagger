build_arm:
	@rm -rf build/
	@mkdir build
	@mkdir build/gui
	@CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -v -o build/musictagger
	@cp gui/data.tpl build/gui/data.tpl
