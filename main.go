package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	config "github.com/nutactl/config"
	constant "github.com/nutactl/constant"
	command "github.com/nutactl/entity"
	host "github.com/nutactl/entity"
	client "github.com/nutactl/http"
	table "github.com/nutactl/table"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const VERSION string = "0.0.1"

func main() {
	cfg := config.MakeConfig()

	fmt.Println("Please insert command")

	for {
		input := prompt.Input("> ", Completer)
		fmt.Println("Your input: " + input)
		cmd := ConvertToCommand(input)

		RunCommand(cfg, cmd)
	}
}
func Completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "hosts <keyword>", Description: "Search host"},
		{Text: "version", Description: "Version"},
		{Text: "exit", Description: "Exit Program"},
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

func RunCommand(cfg config.Config, c command.Command) {
	if c.Cmd == constant.Hosts {
		RunHosts(cfg, c)
	} else if c.Cmd == constant.Exit {
		RunExit()
	} else if c.Cmd == constant.Version {
		fmt.Println("%s", VERSION)
	} else {
		fmt.Println("not implemented.")
	}
}

func RunHosts(cfg config.Config, c command.Command) {
	keyword := ""
	tableUI := table.MakeTable()
	if len(c.Args) > 0 {
		keyword = c.Args[0]
	}

	var hostList = GetAllHostsByKeyword(cfg, keyword)

	var i = 0
	for h := hostList.Front(); h != nil; h = h.Next() {
		var hostEntity = h.Value.(host.Host)
		table.InsertHostData(
			*tableUI,
			i,
			hostEntity.HostName,
			hostEntity.IP)

		i += 1
	}

	fmt.Println(tableUI.String())
}

func RunExit() {

	fmt.Println("Bye, Bye")
	os.Exit(3)
}

func GetAllHostsByKeyword(cfg config.Config, keyword string) list.List {
	filter := fmt.Sprintf("vm_name==.*%s.*", keyword)
	requestPayload := client.MakeVmsListRequestPayload(999, filter)
	resp := client.GetVmsLists(cfg.NutanixUrl+"/vms/list", cfg.UserName, cfg.Password, requestPayload)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Fail nutanix request : %d", resp.StatusCode)
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
				if endpoint.(map[string]interface{})["type"] == constant.Learned {
					ip = endpoint.(map[string]interface{})["ip"].(string)
				}
			}
			hostList.PushBack(host.Host{HostName: hostName, IP: ip})
		}
	}

	return hostList

}
