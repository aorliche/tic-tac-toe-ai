package main

import (
    "fmt"
    ttt "github.com/aorliche/tttai/ttt"
)

func main() {
    nplay := 2
    state := ttt.InitState(3, 3, nplay, 3)
    recvChan := make(chan *ttt.State)
    sendChans := make([]chan *ttt.State, 0)
    for i := 0; i < nplay; i++ {
        sendChans = append(sendChans, make(chan *ttt.State))
    }
    go ttt.Loop(0, sendChans[0], recvChan, 10, 1000)
    go ttt.Loop(1, sendChans[1], recvChan, 10, 1000)
    for {
        fmt.Println("A", state)
        for i := 0; i < 2; i++ {
            sendChans[i] <- state
        }
        fmt.Println("B", state)
        if ttt.GameOver(state) {
            fmt.Println("D", state)
            fmt.Println(ttt.GetLines(state))
            break
        }
        state = <- recvChan
        fmt.Println("C", state)
        fmt.Println("0", ttt.GetLineBonus(state, ttt.GetLines(state), 0), "1", ttt.GetLineBonus(state, ttt.GetLines(state), 1))
    }
}
