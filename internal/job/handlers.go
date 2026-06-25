package job

import (
	"fmt"
	"time"
)

type EmailHandler struct {}

// In general we return an error because something might fail in one of these jobs
// But since we do not have any actual processing logic for now it always returns nil 

func (h EmailHandler) Handle(j Job) error {
	fmt.Println("Sending email...")
	time.Sleep(2 * time.Second)
	
	return nil
}


type FileHandler struct {}

func (h FileHandler) Handle(j Job) error {
	fmt.Println("Processing file...")
	time.Sleep(4 * time.Second)
	
	return nil
}


type ImageHandler struct {}

func (h ImageHandler) Handle(j Job) error {
	fmt.Println("Processing image...")
	time.Sleep(6 * time.Second)

	return nil
}
