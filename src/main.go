package main

import (
    "fmt"
)

func main() {
    state := InitState(5, 5, 2, 4)
    recvChan := make(chan *State)
    sendChans := make([]chan *State, 0)
    for i := 0; i < 2; i++ {
        sendChans = append(sendChans, make(chan *State))
    }
    go Loop(0, sendChans[0], recvChan)
    go Loop(1, sendChans[1], recvChan)
    for {
        fmt.Println("A", state)
        for i := 0; i < 2; i++ {
            sendChans[i] <- state
        }
        fmt.Println("B", state)
        if GameOver(state) {
            fmt.Println("D", state)
            fmt.Println(GetLines(state))
            break
        }
        state = <- recvChan
        fmt.Println("C", state)
    }
}
