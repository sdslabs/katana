package wireguard

import (
	"fmt"
	"log"
	"os"
	"strconv"

	g "github.com/sdslabs/katana/configs"
)

func ApplyFirewall() error {

	//Read challenges folder
	dir, err := os.Open("./challenges")

	if err != nil {
		log.Println("Error in opening challenges folder")
		return err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Println("Error in reading challenges folder")
		return err
	}

	//Store challenge names in a slice
	challengeNames := make([]string, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			challengeNames = append(challengeNames, fileInfo.Name())
		}
	}

	numberOfTeams := g.ClusterConfig.TeamCount
	teamIPs := make([]string, 0)

	baseip := "10.13.13."
	//add team ips to teamIPs
	for i := 0; i < int(numberOfTeams); i++ {
		teamIPs = append(teamIPs, baseip+strconv.Itoa(i+2)+"/32")
	}

	//create a slice of string to store iptables commands
	IpTable := make([]string, 0)

	//add iptables rules to block all internet traffic
	for i := 0; i < int(numberOfTeams); i++ {
		IpTable = append(IpTable, "iptables -I FORWARD -s "+teamIPs[i]+" -o eth+ -j DROP")
	}

	//add iptables rules to allow traffic to all challenges service
	for i := 0; i < int(numberOfTeams); i++ {
		for j := 0; j < len(challengeNames); j++ {
			for k := 0; k < int(numberOfTeams); k++ {
				IpTable = append(IpTable, "iptables -I FORWARD -s "+teamIPs[i]+" -d "+challengeNames[j]+"-svc-"+strconv.Itoa(k)+".katana-team-"+strconv.Itoa(k)+"-ns.svc.cluster.local -j ACCEPT")
			}
		}
	}

	//add iptables rules to allow access to masterpod
	for i := 0; i < int(numberOfTeams); i++ {
		IpTable = append(IpTable, "iptables -I FORWARD -s "+teamIPs[i]+" -d tsuka-svc.katana-team-"+strconv.Itoa(i)+"-ns.svc.cluster.local -j ACCEPT")
	}

	//append all iptables rules to a string
	finalIprules := ""
	for i := 0; i < len(IpTable); i++ {
		finalIprules += IpTable[i] + "; "
	}

	//Overwrite firewall.conf by this string stored in the root directory
	filepath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	filepath = filepath + "/katana-services/Wireguard/root/defaults/firewall.conf"
	err = os.WriteFile(filepath, []byte(finalIprules), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}
