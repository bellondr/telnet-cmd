package main

import (
	"fmt"
	telnet "github.com/wNee/telnet-cmd"
)

func main() {
	client, err := telnet.NewClient("10.169.231.19:23")
	if err != nil {
		fmt.Println("new client err: ", err.Error())
		return
	}
	err = client.Login("root", "")
	if err != nil {
		fmt.Println("login err: ", err.Error())
		return
	}

	out, err := client.RunCmdWithOutput("touch tl-cmd.txt")
	if err != nil {
		fmt.Println("run cmd err: ", err.Error())
		return
	}
	fmt.Println(out)
}
