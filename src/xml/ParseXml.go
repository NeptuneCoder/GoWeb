package main

import (
	fmt "fmt"
	"encoding/xml"
)

type Server struct {
	ServerName string
	ServerIP string
}
type Serverslice struct {
	Server []Server

}
func main(){
	var s Serverslice
	s.Server = append(s.Server,Server{ServerName:"Shanghai_VPN",ServerIP:"127.0.0.1"})
	s.Server = append(s.Server,Server{ServerName:"Shanghai_VPN",ServerIP:"162.0.0.1"})

	b,err := xml.Marshal(s)

	if err != nil{
		fmt.Println("json err:", err)
	}

	fmt.Println(string(b))

}