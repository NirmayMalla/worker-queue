package job

import (
	"fmt"
	"time"
)

type EmailHandler struct {}

func (h EmailHandler) Handle(j Job) {
	fmt.Println("Sending email...")
	time.Sleep(2 * time.Second)
}


type FileHandler struct {}

func (h FileHandler) Handle(j Job) {
	fmt.Println("Processing file...")
	time.Sleep(4 * time.Second)
}


type ImageHandler struct {}

func (h ImageHandler) Handle(j Job) {
	fmt.Println("Processing image...")
	time.Sleep(6 * time.Second)
}
