package flaghandlerservice

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"sync"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func Server(wg sync.WaitGroup) {
	ln, err := net.Listen("tcp", configs.FlagConfig.SubmissionServicePort)
	if err != nil {
		log.Fatal(ServiceFail, err)
	}
	defer func() {
		ln.Close()
		wg.Done()
	}()
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
			return

		default:
			if (types.CTFTeam{}) == connectedTeam {
				writeToCient(conn, NoLogin)
			} else if status, points := submitFlag(cmd, connectedTeam); status {
				connectedTeam.Score = connectedTeam.Score + points
				writeToCient(conn, SubmitSuccess+strconv.Itoa(points)+TotalScore+strconv.Itoa(connectedTeam.Score)+"\n")
			} else {
				writeToCient(conn, InvalidFlag)
			}
		}
	}
}
func parseCommand(cmdLine string) (cmd, param, password string) {
	r, _ := regexp.Compile("(init|exit|[A-Za-z0-9]+)([[:blank:]]([A-Za-z0-9]+)[[:blank:]]([A-Za-z0-9]+))?")
	matched := r.FindStringSubmatchIndex(cmdLine)
	if matched[6] == -1 {
		return cmdLine[matched[2]:matched[3]], "", ""
	}
	return cmdLine[matched[2]:matched[3]], cmdLine[matched[6]:matched[7]], cmdLine[matched[8]:matched[9]]
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
