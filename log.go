package main

import (
	"fmt"
	"strings"
)

func Info(msg string){
	fmt.Println("info  ==> ", msg)
}

func Error(err error) {
	if err  != nil {
		fmt.Println("error ==> " , err)
	}
}

func Debug(msg string) {
	fmt.Println("debug ==> ",strings.TrimSpace(msg))
}
