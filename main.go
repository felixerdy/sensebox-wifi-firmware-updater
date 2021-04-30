package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"

	"github.com/arduino/FirmwareUploader/modules/winc"
	"github.com/arduino/FirmwareUploader/utils/context"
)

var ctx = &context.Context{}
var statusText = widget.NewLabel("")
var logBuffer = new(bytes.Buffer)

func main() {
	ctx.FWUploaderBinary = "WINC1500_Updater.ino.sensebox_mcu.bin"
	ctx.Addresses.Set("api.opensensemap.org:443")
	ctx.Addresses.Set("api.testing.opensensemap.org:443")
	ctx.Addresses.Set("api.telegram.org:443")
	ctx.BinaryToRestore = "WINC1500_Updater.ino.sensebox_mcu.bin"
	ctx.ProgrammerPath = "darwin/amd64/bossac"
	ctx.FirmwareFile = "firmwares/m2m_aio_3a0.bin"

	fmt.Println(ctx.Addresses)

	log.SetOutput(logBuffer)
	statusText.SetText(logBuffer.String())

	a := app.New()

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
		statusText.SetText(logBuffer.String())
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
		statusText.SetText(logBuffer.String())
	}

	w := a.NewWindow("senseBox WiFi Firmware Updater")
	w.Resize(fyne.NewSize(400, 300))

	hello := widget.NewLabel("Wähle die senseBox MCU aus")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewSelect(ports, func(value string) {
			ctx.PortName = value
		}),
		statusText,
		container.New(layout.NewCenterLayout(), widget.NewButton("Zertifikate und Firmware Flashen", flash)),
	))

	w.ShowAndRun()
}

func flash() {
	statusText.SetText("Starting Flashing")
	retry := 0
	for {
		ctxCopy := *ctx

		err := winc.Run(&ctxCopy)

		if err == nil {
			log.Println("Operation completed: success! :-)")
			statusText.SetText(logBuffer.String())
			break
		}
		log.Println("Error: " + err.Error())

		if retry >= ctx.Retries {
			log.Fatal("Operation failed. :-(")
			statusText.SetText(logBuffer.String())
		}

		retry++
		log.Println("Waiting 1 second before retrying...")
		statusText.SetText(logBuffer.String())
		time.Sleep(time.Second)
		log.Printf("Retrying upload (%d of %d)", retry, ctx.Retries)
		statusText.SetText(logBuffer.String())
	}
}
