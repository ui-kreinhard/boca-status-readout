package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	ticketCountQuery = "/html/body/table[1]/tbody/tr/td[2]/table/tbody/tr/td[2]/table/tbody/tr[1]/td[2]"
	readyQuery       = "/html/body/table[1]/tbody/tr[1]/td[1]/table/tbody/tr[1]/td[2]"
	paperLowQuery    = "/html/body/table[1]/tbody/tr[1]/td[1]/table/tbody/tr[2]/td[2]"
	paperOutQuery    = "/html/body/table[1]/tbody/tr[1]/td[1]/table/tbody/tr[3]/td[2]"
	paperJamQuery    = "/html/body/table[1]/tbody/tr[1]/td[1]/table/tbody/tr[4]/td[2]"
	cutterJamQuery   = "/html/body/table[1]/tbody/tr[1]/td[1]/table/tbody/tr[5]/td[2]"
)

type PrinterMaintenanceStatus string

type PrinterStatus struct {
	TicketCount int  `json:"ticket_count"`
	Ready       bool `json:"ready"`
	PaperOut    bool `json:"paper_out"`
	PaperLow    bool `json:"paper_low"`
	PaperJam    bool `json:"paper_jam"`
	CutterJam   bool `json:"cutter_jam"`
}

func basicAuth(url string) (io.Reader, error) {
	var username string = "boca"
	var passwd string = "printer"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		return req.Body, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("401")
	}
	return resp.Body, nil
}

func LoadURL(url string) (*html.Node, error) {
	r, err := basicAuth(url)
	if err != nil {
		return nil, err
	}
	return html.Parse(r)
}

func getPrinterVersion(ip string) (string, error) {
	doc, err := LoadURL("http://" + ip)
	if err != nil {
		return "", err
	}
	formTarget, err := readOutByQuery(doc, "//frame/@src")
	if err != nil {
		return "", err
	}
	switch formTarget {
	case "gets_protected.cgi":
		return "R", nil
	default:
		return "N", nil
	}
}

func buildUrlNew(ip string) string {
	return fmt.Sprintf("http://%v/realtime.cgi", ip)
}

func FetchStatus(ip string) (*PrinterStatus, error) {
	version, err := getPrinterVersion(ip)
	if err != nil {
		return nil, err
	}
	url := buildUrl(ip)
	if version == "R" {
		url = buildUrlNew(ip)
	}
	doc, err := LoadURL(url)
	if err != nil {
		return nil, err
	}

	ps := NewPrinterStatus()

	ps.TicketCount, err = readOutTicketCount(doc)
	if err != nil {
		return nil, err
	}
	ps.Ready, err = readOutMaintenanceStatus(doc, readyQuery)
	if err != nil {
		return nil, err
	}
	ps.PaperLow, err = readOutMaintenanceStatus(doc, paperLowQuery)
	if err != nil {
		return nil, err
	}
	ps.PaperOut, err = readOutMaintenanceStatus(doc, paperOutQuery)
	if err != nil {
		return nil, err
	}
	ps.PaperJam, err = readOutMaintenanceStatus(doc, paperJamQuery)
	if err != nil {
		return nil, err
	}
	ps.CutterJam, err = readOutMaintenanceStatus(doc, cutterJamQuery)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func NewPrinterStatus() *PrinterStatus {
	return &PrinterStatus{}
}

func (ps *PrinterStatus) String() string {
	return fmt.Sprintf("Ticket count: %d\nReady: %t", ps.TicketCount, ps.Ready)
}

func (ps *PrinterStatus) ToJson() string {
	ret, _ := json.Marshal(ps)
	return string(ret)
}

func (p *PrinterStatus) GetIntStatus() int {
	ret := 0
	if p.PaperLow {
		ret |= 32768
	}
	if p.PaperOut {
		ret |= 16384
	}
	if p.PaperJam {
		ret |= 1024
	}
	if !p.Ready {
		ret |= 512
	}
	return ret
}

func readOutTicketCount(doc *html.Node) (int, error) {
	ticketCountString, err := readOutByQuery(doc, ticketCountQuery)
	if err != nil {
		return -1, err
	}
	ticketCount, err := strconv.Atoi(ticketCountString)
	if err != nil {
		return -1, err
	}
	return ticketCount, nil
}

func readOutMaintenanceStatus(doc *html.Node, query string) (bool, error) {
	response, err := readOutByQuery(doc, query)
	if err != nil {
		return false, err
	}
	return isTrue(response), nil
}

func readOutByQuery(doc *html.Node, query string) (string, error) {
	nodes, err := htmlquery.QueryAll(doc, query)
	if err != nil {
		return "", err
	}
	if len(nodes) != 1 {
		for _, node := range nodes {
			log.Println(htmlquery.InnerText(node))
		}
		return "", errors.New(fmt.Sprintf("abigous query result len: %d", len(nodes)))
	}
	node := nodes[0]
	return strings.TrimRight(htmlquery.InnerText(node), " "), nil
}

func isTrue(s string) bool {
	return s == "YES"
}

func buildUrl(ip string) string {
	return fmt.Sprintf("http://%v/realtime.htm", ip)
}
