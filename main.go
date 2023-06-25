package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ui-kreinhard/boca-status-readout/query"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("No ip given")
	}
	printerIp := os.Args[1]
	mode := "count"
	if len(os.Args) == 3 {
		mode = os.Args[2]
	}
	printerStatus, err := query.FetchStatus(printerIp)

	if err != nil {
		log.Fatalln(err)
	}

	switch mode {
	case "count":
		fmt.Println(printerStatus.TicketCount)
		return
	case "status":
		fmt.Println(printerStatus.GetIntStatus())
	case "json":
		fmt.Println(printerStatus.ToJson())
	}
}
