package main

import (
	"fmt"
	"log"
	"testing"
)

type Factory func() Service

type Service struct {
}

type FactoryV1 interface {
	Create() Service
}

type AbstractFactory struct {
}

func (f *AbstractFactory) Create() {
	f.CreatePart1()
	f.CreatePart2()
}

func (f *AbstractFactory) CreatePart1() {

}

func (f *AbstractFactory) CreatePart2() {
	log.Println("这是 AbstractFactory 的 CreatePart2")
}

type MyFactory struct {
	AbstractFactory
}

func (f *MyFactory) CreatePart2() {
	log.Println("这是 MyFactory 的 CreatePart2")
}

func TestMyFactory(t *testing.T) {
	f := &MyFactory{}
	f.Create()
}

type AbstractFactoryV1 struct {
	CreatePart1 func()
	CreatePart2 func()
}

func (f *AbstractFactoryV1) Create() {
	f.CreatePart1()
	f.CreatePart2()
}

type MyFactoryV1 struct {
	AbstractFactoryV1
}

func NewMyFactoryV1() MyFactoryV1 {
	return MyFactoryV1{
		AbstractFactoryV1{
			CreatePart2: func() {
				fmt.Println("创建第二部分")
			},
			CreatePart1: func() {
				fmt.Println("创建第一部分")
			},
		},
	}
}
