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

func server() {
	ln, err := net.Listen(configs.FlagConfig.SubmissionServicePort, "tcp")
	if err != nil {
		log.Fatal(ServiceFail)
	}
	defer ln.Close()
	log.Println(ServiceSuccess, configs.FlagConfig.SubmissionServicePort)
	connectedTeam := types.CTFTeam{}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println(ClosingError, err)
			}
			continue
		}

		log.Println(Connected, conn.RemoteAddr())
		go handleConnection(conn, connectedTeam)
	}
}

func handleConnection(conn net.Conn, connectedTeam types.CTFTeam) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(ClosingError, err)
		}
	}()
	writeToCient(conn, InitInstruction)

	for {
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)

		if n == 0 || err != nil {
			log.Println(ReadError, err)
			return
		}

		cmd, param, password := parseCommand(string(cmdLine[0:n]))

		if cmd == "" {
			writeToCient(conn, InvalidCommand)
			continue
		}
		switch cmd {
		case "init":
			if param == "" || password == "" {
				writeToCient(conn, InvalidParams)
				continue
			} else if (types.CTFTeam{}) != connectedTeam {
				writeToCient(conn, TeamAlreadyExists)
				continue
			} else {
				if condition, team := checkTeam(param, password); condition {
					connectedTeam = team
					writeToCient(conn, TeamConnected)
					continue
				} else {
					writeToCient(conn, InvalidCreds)
					continue
				}
			}
		case "exit":

		default:
			if (types.CTFTeam{}) == connectedTeam {
				writeToCient(conn, NoLogin)
			} else if status, points := submitFlag(cmd, connectedTeam); status {
				writeToCient(conn, SubmitSuccess+strconv.Itoa(points)+"\n")
			} else {
				writeToCient(conn, InvalidFlag)
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

func writeToCient(conn net.Conn, message string) {
	if _, err := conn.Write([]byte(message)); err != nil {
		log.Println(WriteError, err)
		return
	}
}
