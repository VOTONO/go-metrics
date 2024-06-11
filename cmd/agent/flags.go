package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

type Interval struct {
	Value int
}

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

func (a Interval) String() string {
	return fmt.Sprintf("%v", a.Value)
}

func (a *Interval) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("interval value, must be an integer")
	}
	a.Value = val
	return nil
}
