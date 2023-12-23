package main 

import (
    "fmt"
    "time"
)

func Loop(me int, inChan chan *State, outChan chan *State) {
    var state *State
    for { 
        state = <- inChan 
        if GameOver(state) {
            break
        }
        fmt.Println(me, "recvd:", state)
        next := Search(state, me)
        fmt.Println(me, "after search state:", state)
        if next == nil {
            time.Sleep(100 * time.Millisecond)
            continue
        }
        outChan <- next
    }
}

// Set up iterative deepening
func Search(state *State, me int) *State {
    startTime := time.Now()
    var res *State
    for d := 1; d < 5; d++ {
        _, fn, fin := SearchDeep(state, me, d, startTime)
        if fn != nil && fin {
            res = fn()
        } else {
            break;
        }
    }
    return res
}

// Iterative deepening worker
// Standard minimax without alpha-beta pruning
// Allow players to make consecutive moves
// If the game rules allow it
func SearchDeep(state *State, me int, d int, startTime time.Time) (*State, func()*State, bool) {
    if d == 0 {
        return state, nil, true
    }
    if time.Since(startTime) > 2*time.Second {
        return nil, nil, false
    }
    fns := GetCandidates(state, me)
    if len(fns) == 0 {
        return nil, nil, true
    }
    vals := make([]float64, len(fns))
    states := make([]*State, len(fns))
    // Special for consecutive moves
    valsCons := make([]float64, len(fns))
    statesCons := make([]*State, len(fns))
    for i,fn := range fns {
        next := fn()
        if !GameOver(next) {
            // Min of minimax
            // Except that we have the possibility of consecutive moves
            // If I have a good move, assume I go first
            // Don't assume that a valid move can be found
            vals2 := make([]float64, len(fns))
            states2 := make([]*State, len(fns))
            for j := 0; j < state.NPlayers; j++ {
                n, _, fin := SearchDeep(next, j, d-1, startTime)
                if fin {
                    if n != nil {
                        if j == me {
                            valsCons[i] = Eval(n, me)
                            statesCons[i] = n
                        // Evaluate with respect to opponent
                        } else {
                            vals[i] = Eval(n, j)
                            states[i] = n
                        }
                    }
                } else {
                    return nil, nil, false
                }
            }
            best := 0
            for j := 1; j < len(vals2); j++ {
                if states2[j] != nil && vals2[j] > vals2[best] {
                    best = j
                }
            }
            // Re-evaluate with respect to me
            if states2[best] != nil {
                vals[i] = Eval(states2[best], me)
                states[i] = states2[best]
            }
        } else {
            vals[i] = Eval(next, me)
            states[i] = next
        }
    }
    // Combine with consecutive moves
    for i := 0; i < len(valsCons); i++ {
        if statesCons[i] != nil && valsCons[i] > vals[i] {
            vals[i] = valsCons[i]
            states[i] = statesCons[i]
        }
    }
    best := 0
    for i := 1; i < len(vals); i++ {
        if states[i] != nil && vals[i] > vals[best] {
            best = i
        }
    }
    return states[best], fns[best], true
}
