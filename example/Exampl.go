package main

import (
    "github.com/tietang/props"
    "fmt"
)

func main() {

    root := "config/app1/dev"
    address := "127.0.0.1:8500"
    p := props.NewConsulIniConfigSourceByName("consul-ini", address, root)
    fmt.Println(p.Get(""))
}
