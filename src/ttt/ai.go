package ttt

import (
    //"fmt"
    "time"
)

func Loop(me int, inChan chan *State, outChan chan *State, depth int, timeMillis int) {
    var state *State
    for { 
        state = <- inChan 
        // nil state indicates player disconnect
        if state == nil {
            break
        }
        if GameOver(state) {
            break
        }
        //fmt.Println(me, "recvd:", state)
        next := Search(state, me, depth, timeMillis)
        //fmt.Println(me, "after search state:", state)
        if next == nil {
            time.Sleep(100 * time.Millisecond)
            continue
        }
        outChan <- next
    }
}

// Set up iterative deepening
func Search(state *State, me int, depth int, timeMillis int) *State {
    startTime := time.Now()
    var res, resNonZero *State
    for d := 1; d < depth; d++ {
        st, fn, fin := SearchDeep(state, me, d, startTime, timeMillis)
        /*if st != nil {
            if d == 1 {
                fmt.Println("state", state)
            }
            fmt.Println("me", me, "d", d, "see", st, "val", Eval(st, me))
        }*/
        if fn != nil && fin {
            res = fn()
            if Eval(res, me) == 1 {
                return res
            }
            if Eval(st, me) == 0 {
                if resNonZero != nil {
                    return resNonZero
                }
            } else if Eval(st, me) > 0 {
                resNonZero = res

            }
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
func SearchDeep(state *State, me int, d int, startTime time.Time, timeMillis int) (*State, func()*State, bool) {
    if d == 0 {
        return state, nil, true
    }
    if time.Since(startTime).Milliseconds() > int64(timeMillis) {
        return nil, nil, false
    }
    fns := GetCandidates(state, me)
    if len(fns) == 0 {
        return nil, nil, true
    }
    vals := make([]float64, len(fns))
    states := make([]*State, len(fns))
    // Special for consecutive moves
    /*valsCons := make([]float64, len(fns))
    statesCons := make([]*State, len(fns))*/
    for i,fn := range fns {
        next := fn()
        if !GameOver(next) {
            // Check opponents responses 
            // We have the possibility of consecutive moves
            // If I have a good move, assume I go first
            // Don't assume that a valid move can be found
            vals2 := make([]float64, state.NPlayers)
            states2 := make([]*State, state.NPlayers)
            for j := 0; j < state.NPlayers; j++ {
                n, _, fin := SearchDeep(next, j, d-1, startTime, timeMillis)
                if fin {
                    if n != nil {
                        /*if j == me {
                            valsCons[i] = Eval(n, me)
                            statesCons[i] = n
                        } else {*/
                            vals2[j] = Eval(n, me)
                            states2[j] = n
                        //}
                    }
                } else {
                    return nil, nil, false
                }
            }
            // Min of minimax
            best := -1
            for j := 0; j < len(vals2); j++ {
                if states2[j] != nil && (best == -1 || vals2[j] < vals2[best]) {
                    best = j
                }
            }
            if best != -1 {
                vals[i] = vals2[best]
                states[i] = states2[best]
            }
        } else {
            vals[i] = Eval(next, me)
            states[i] = next
        }
    }
    // Combine with consecutive moves
    /*for i := 0; i < len(valsCons); i++ {
        if statesCons[i] != nil && valsCons[i] > vals[i] {
            vals[i] = valsCons[i]
            states[i] = statesCons[i]
        }
    }*/
    best := -1
    for i := 0; i < len(vals); i++ {
        if states[i] != nil && (best == -1 || vals[i] > vals[best]) {
            best = i
        }
    }
    return states[best], fns[best], true
}
