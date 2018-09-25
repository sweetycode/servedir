package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/fatih/color"
)

func getIntranetIP() (ipList string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil && !v.IP.IsLoopback() {
					return v.IP.String(), nil
				}
			case *net.IPAddr:
				if v.IP.To4() != nil && !v.IP.IsLoopback() {
					return v.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("not found")
}

func main() {
	webRoot := "."
	listenPort := 8000

	fs := http.FileServer(http.Dir(webRoot))
	http.Handle("/", fs)

	color.Yellow("Listening...")
	color.Yellow("You can access with:")
	ip, _ := getIntranetIP()
	color.Green("    http://%s:%d\n\n", ip, listenPort)

	files, err := ioutil.ReadDir(webRoot)
	if err == nil {
		color.Yellow("List files: ")
		for _, file := range files {
			path := file.Name()
			if file.IsDir() {
				path += "/"
			}
			color.Blue("    http://%s:%d/%s\n", ip, listenPort, path)
		}
	}
	fmt.Println("")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil); err != nil {
		color.Red("error happend: %v", err)
	}
}
