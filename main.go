package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/c-bata/go-prompt"
	config "github.com/nutactl/config"
	command "github.com/nutactl/entity"
	host "github.com/nutactl/entity"
	client "github.com/nutactl/http"
	table "github.com/nutactl/table"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "hosts", Description: "Search host"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func ConvertToCommand(input string) command.Command {
	s := strings.Split(input, " ")

	return command.Command{
		Cmd:  s[0],
		Args: s[1:],
	}
}
func main() {
	cfg := config.MakeConfig()

	fmt.Println("Please insert command")
	tableUI := table.MakeTable()

	for {
		input := prompt.Input("> ", Completer)
		fmt.Println("Your input: " + input)
		cmd := ConvertToCommand(input)

		RunCommand(cfg, cmd, *tableUI)
	}
}

func RunCommand(cfg config.Config, c command.Command, tableUI simpletable.Table) {
	if c.Cmd == "hosts" {
		RunHosts(cfg, c, tableUI)
	} else if c.Cmd == "exit" {
		RunExit()
	} else {
		fmt.Println("not implemented.")
	}
}

func RunHosts(cfg config.Config, c command.Command, tableUI simpletable.Table) {
	keyword := ""
	if len(c.Args) > 0 {
		keyword = c.Args[0]
	}

	var hostList = GetAllHostsByKeyword(cfg, keyword)

	var i = 0
	for h := hostList.Front(); h != nil; h = h.Next() {
		var hostEntity = h.Value.(host.Host)
		table.InsertHostData(
			tableUI,
			i,
			hostEntity.HostName,
			hostEntity.IP)

		i += 1
	}

	fmt.Println(tableUI.String())
}

func RunExit() {

	fmt.Println("Bye, Bye 👋")
	os.Exit(3)
}

func GetAllHostsByKeyword(cfg config.Config, keyword string) list.List {
	requestPayload := client.MakeVmsListRequestPayload(999)
	resp := client.GetVmsLists(cfg.NutanixUrl+"/vms/lists", cfg.UserName, cfg.Password, requestPayload)
	if resp.StatusCode != http.StatusOK {
		print("Fail nutanix request")
		return list.List{}
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	var hostList = list.List{}
	var entities = result["entities"].([]interface{})
	for _, entity := range entities {
		var status = entity.(map[string]interface{})["status"].(map[string]interface{})
		var hostName = status["name"].(string)
		if keyword != "" && strings.Contains(hostName, keyword) == false {
			continue
		}
		var ip = ""
		var resources = status["resources"].(map[string]interface{})
		var nicList = resources["nic_list"].([]interface{})
		if len(nicList) > 0 {
			var ipEndPointList = nicList[0].(map[string]interface{})["ip_endpoint_list"].([]interface{})
			for _, endpoint := range ipEndPointList {
				if endpoint.(map[string]interface{})["type"] == "LEARNED" {
					ip = endpoint.(map[string]interface{})["ip"].(string)
				}
			}
			hostList.PushBack(host.Host{HostName: hostName, IP: ip})
		}
	}

	return hostList

}
