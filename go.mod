module sensebox.de/wifi-firmware-updater

go 1.16

replace go.bug.st/serial => github.com/cmaglie/go-serial v0.0.0-20200923162623-b214c147e37e

require (
	fyne.io/fyne/v2 v2.0.2
	github.com/arduino/FirmwareUploader v0.0.0-20210427135349-644924404dc0
	go.bug.st/serial v1.1.3
)
