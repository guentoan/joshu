APP_NAME := joshu
VERSION := 1.0.0

# Fyne command
FYNE := fyne package

.PHONY: all clean build-macos build-linux build-windows

all: clean build-macos build-linux build-windows

clean:
	rm -rf $(APP_NAME).app
	rm -rf $(APP_NAME).zip
	rm -rf $(APP_NAME).exe

build-macos: $(MACOS_DIR)
	$(FYNE) -os darwin -icon icon.png -name $(APP_NAME) -appVersion $(VERSION)

build-linux: $(LINUX_DIR)
	$(FYNE) -os linux -icon icon.png -name $(APP_NAME) -appVersion $(VERSION)

build-windows: $(WINDOWS_DIR)
	$(FYNE) -os windows -icon icon.png -name $(APP_NAME) -appVersion $(VERSION)

# Build all versions
build: all

# To package the app for a specific platform, you can run:
# make build-macos
# make build-windows
# make build-linux
