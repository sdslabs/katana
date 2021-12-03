package flaghandlerservice

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

type Team struct {
	TeamName string
	TeamID   int
}

func server() {
	ln, err := net.Listen(configs.FlagConfig.SubmissionServicePort, "tcp")
	if err != nil {
		log.Fatal("Failed to Start Flag Submission Service")
	}
	defer ln.Close()
	log.Println("Flag Submission Service Started at port", configs.FlagConfig.SubmissionServicePort)
	connectedTeam := Team{}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("Failed to close", err)
			}
			continue
		}

		log.Println("Connected to", conn.RemoteAddr())
		go handleConnection(conn, connectedTeam)
	}
}

func handleConnection(conn net.Conn, connectedTeam Team) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error Closing", err)
		}
	}()
	writeToCient(conn, "Connected to Flag Submission Service\nInitiate your session by `init <teamName> <password>`\n")

	for {
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)

		if n == 0 || err != nil {
			log.Println("Connection Read err", err)
			return
		}

		cmd, param, password := parseCommand(string(cmdLine[0:n]))

		if cmd == "" {
			writeToCient(conn, "Inavlid Command\n")
			continue
		}
		switch cmd {
		case "init":
			if param == "" || password == "" {
				writeToCient(conn, "Invalid Login Parameters\n")
				continue
			} else if (Team{}) != connectedTeam {
				writeToCient(conn, "Team is already Logged in\n")
				continue
			} else {
				if checkTeam(param) {
					connectedTeam.TeamAddress = conn.RemoteAddr().String()
					connectedTeam.TeamID = param
					writeToCient(conn, "Team successfully connected,\n Enter flags to submit them\n")
					continue
				} else {
					writeToCient(conn, "Invalid TeamID\n")
					continue
				}
			}
		case "exit":

		default:
			if status, points := submitFlag(cmd); status {
				writeToCient(conn, "Submitted successfully, points:"+strconv.Itoa(points)+"\n")
			} else {
				writeToCient(conn, "Invalid Flag")
			}
		}
	}
}

func parseCommand(cmdLine string) (cmd, param, password string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) == 3 {
		cmd = strings.TrimSpace(parts[0])
		param = strings.TrimSpace(parts[1])
		password = strings.TrimSpace(parts[2])
		return
	}
	if len(parts) == 2 {
		cmd = strings.TrimSpace(parts[0])
		param = strings.TrimSpace(parts[1])
		password = ""
		return
	}
	if len(parts) == 1 {
		cmd = strings.TrimSpace(parts[0])
		param = ""
		password = ""
		return
	}
	return "", "", ""
}

func checkTeam(teamName string, password string) (bool, types.CTFTeam) {
	team := &types.CTFTeam{}
	if team, err := mongo.FetchSingleTeam(teamName); err == nil {
		if utils.CompareHashWithPassword(team.Password, password) {
			return true, *team
		}
	}
	return false, *team
}

func submitFlag(flag string) (bool, int) {
	return true, 10
}

func writeToCient(conn net.Conn, message string) {
	if _, err := conn.Write([]byte(message)); err != nil {
		log.Println("failed to write", err)
		return
	}
}
