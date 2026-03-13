package main

import (
	"fmt"
)

type Server struct {
	Hostname string
	IP       string
	Port     int
	OS       string
}

// 設定を一行で表示する
func (s Server) Display() {
	fmt.Printf("Hostname: %s | IP: %s | Port: %d | OS: %s\n", s.Hostname, s.IP, s.Port, s.OS)
}

func (s *Server) UpdateOS(newOS string) {
	s.OS = newOS
}

func main() {
	myServer := Server{
		Hostname: "db-master",
		IP:       "10.0.0.5",
		Port:     5432,
		OS:       "AlmaLinux9",
	}

	fmt.Print("--- 更新前 ---\n")
	myServer.Display()

	// データの書き換えを実行
	myServer.UpdateOS("RockyLinux9")

	fmt.Print("--- 更新後 ---\n")
	myServer.Display()
}
