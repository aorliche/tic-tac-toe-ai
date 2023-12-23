package main 

import (
    "time"
)

func Loop(state *State, me int, stateChan chan *State) {
    for {
        next := Search(state, me)
        if next == nil {
            time.Sleep(100 * time.Millisecond)
            continue
        }
        stateChan <- next
    }
}

// Set up iterative deepening
func Search(state *State, me int) *State {
    startTime := time.Now()
    var res *State
    for d := 1; d < 20; d++ {
        state, fin := SearchDeep(state, me, d, startTime)
        if fin {
            res = state
        } else {
            break;
        }
    }
    return res
}

// Iterative deepening worker
func SearchDeep(state *State, me int, d int, startTime time.Time) (*State, bool) {
    if d == 0 {
        return state, true
    }
    if time.Since(startTime) > time.Second {
        return nil, false
    }
    fns := GetCandidates(state, me)
    if len(fns) == 0 {
        return nil, true
    }
    vals := make([]float64, len(fns))
    states := make([]*State, len(fns))
    for _,fn := range fns {
        next := fn()
        if !GameOver(next) {
            n, fin := SearchDeep(next, (me+1)%state.NPlayers, d-1, startTime)
            if fin {
                next = n
            } else {
                return nil, false
            }
        }
        vals = append(vals, Eval(next, me))
        states = append(states, next)
    }
    best := 0
    for i := 1; i < len(vals); i++ {
        if vals[i] > vals[best] {
            best = i
        }
    }
    return states[best], true
}
