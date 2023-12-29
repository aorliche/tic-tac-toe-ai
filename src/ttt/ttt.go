package ttt

import (
    "math"
)

type State struct {
    Rows int
    Cols int
    Board [][]int
    Turn int
    NPlayers int
    WinThresh int
}

func InitState(rows int, cols int, nplay int, winthresh int) *State {
    board := make([][]int, rows)
    for i := 0; i < rows; i++ {
        board[i] = make([]int, cols)
        for j := 0; j < cols; j++ {
            board[i][j] = -1
        }
    }
    return &State{rows, cols, board, 0, nplay, winthresh}
}

func (state *State) Clone() *State {
    s := InitState(state.Rows, state.Cols, state.NPlayers, state.WinThresh)
    s.Turn = state.Turn
    s.Board = make([][]int, state.Rows)
    for i := 0; i < state.Rows; i++ {
        s.Board[i] = make([]int, state.Cols)
        copy(s.Board[i], state.Board[i])
    }
    return s
}

func GetRows(state *State) [][]int {
    rows := make([][]int, state.Rows)
    for i := 0; i < state.Rows; i++ {
        rows[i] = make([]int, state.Cols)
        copy(rows[i], state.Board[i])
    }
    return rows
}

func GetCols(state *State) [][]int {
    cols := make([][]int, state.Cols)
    for i := 0; i < state.Cols; i++ {
        cols[i] = make([]int, state.Rows)
        for j := 0; j < state.Rows; j++ {
            cols[i][j] = state.Board[j][i]
        }
    }
    return cols
}

func GetDiags(state *State) [][]int {
    maxDiag := state.Rows
    if state.Cols < state.Rows {
       maxDiag = state.Cols 
    }
    // Index
    m := 0
    // Duplicate the two main diagonals
    // Eliminates some edge cases?
    n := 2*(state.Rows+state.Cols)-2
    diags := make([][]int, n)
    // Left to right start at bottom row
    for i := 0; i < state.Rows; i++ {
        dsize := i+1
        if dsize > maxDiag {
            dsize = maxDiag
        }
        diags[m] = make([]int, dsize)
        for j := 0; j < dsize; j++ {
            diags[m][j] = state.Board[state.Rows-i-1+j][j]
        }
        m++
    }
    // Left to right start at left col
    for i := 1; i < state.Cols; i++ {
        dsize := state.Cols - i
        if dsize > maxDiag {
            dsize = maxDiag
        }
        diags[m] = make([]int, dsize)
        for j := 0; j < dsize; j++ {
            diags[m][j] = state.Board[j][i+j]
        }
        m++
    }
    // Right to left start at bottom row
    for i := state.Rows-1; i >= 0; i-- {
        dsize := state.Rows - i
        if dsize > maxDiag {
            dsize = maxDiag
        }
        diags[m] = make([]int, dsize)
        for j := 0; j < dsize; j++ {
            diags[m][j] = state.Board[i+j][state.Cols-j-1]
        }
        m++
    }
    // Right to left start at right col
    for i := state.Cols-2; i >= 0; i-- {
        dsize := i+1
        if dsize > maxDiag {
            dsize = maxDiag
        }
        diags[m] = make([]int, dsize)
        for j := 0; j < dsize; j++ {
            diags[m][j] = state.Board[j][i-j]
        }
        m++
    }
    return diags
}

func GetLines(state *State) [][]int {
    lines := make([][]int, 0)
    lines = append(lines, GetRows(state)...)
    lines = append(lines, GetCols(state)...)
    lines = append(lines, GetDiags(state)...)
    return lines
}

func (state *State) GetLineWinner(line []int) int {
    if len(line) < state.WinThresh {
        return -1
    }
    prev := -1
    seq := 0
    for i := 0; i < len(line); i++ {
        if line[i] != prev {
            prev = line[i]
            seq = 1
        } else if prev != -1 {
            seq++
            if seq == state.WinThresh {
                return prev
            }
        }
    }
    return -1
}

func GetWinner(state *State) (int, [][]int) {
    lines := GetLines(state)
    for _,line := range lines {
        winner := state.GetLineWinner(line)
        if winner != -1 {
            return winner, lines
        }
    }
    return -1, lines
}

func GetCenterBonus(state *State, me int) float64 {
    sum := float64(0)
    r := float64(state.Rows-1)
    c := float64(state.Cols-1)
    mul := 1/((r+1)*(c+1))
    ih := r/2
    jh := c/2
    for i := 0; i < state.Rows; i++ {
        for j := 0; j < state.Cols; j++ {
            if state.Board[i][j] == me {
                sum += ((ih-math.Abs(float64(i)-ih)) + (jh-math.Abs(float64(j)-jh)))*mul
            }
        }
    }
    return sum
}

// 1. Must have potential for making winthresh in a row
// 2. Is close to making winthresh in a row
func GetLineBonus(state *State, lines [][]int, me int) float64 {
    sum := float64(0)
    // Check how many lines can potentially have a win
    npot := 0
    for _,line := range lines {
        if len(line) >= state.WinThresh {
            npot++
        }
    }
    mul := 1/float64(state.WinThresh)/float64(npot)/2
    for _,line := range lines {
        if len(line) < state.WinThresh {
            continue
        }
        n := 0
        m := 0
        for i := 0; i < len(line); i++ {
            if line[i] == -1 || line[i] == me {
                n++
                if line[i] == me {
                    m++
                }
            }
            if (line[i] != -1 && line[i] != me) || i == len(line)-1 {
                if n >= state.WinThresh && m >= 1 {
                    sum += float64(m)*mul + 0.2*float64(n)*mul
                }
                n = 0
                m = 0
            }
        }
    }
    return sum
}

func Eval(state *State, me int) float64 {
    winner, lines := GetWinner(state)
    if winner == me {
        return 10.0
    } else if winner == -1 {
        cBonus := GetCenterBonus(state, me)
        lBonus := GetLineBonus(state, lines, me)
        return 0.5+lBonus+cBonus
    }
    return 0
}

func GameOver(state *State) bool {
    win, _ := GetWinner(state)
    if win != -1 {
        return true
    }
    for i := 0; i < state.Rows; i++ {
        for j := 0; j < state.Cols; j++ {
            if state.Board[i][j] == -1 {
                return false
            }
        }
    }
    return true
}

func CanPlay(state *State, me int, i int, j int) bool {
    if (state.Turn-me) % state.NPlayers != 0 {
        return false
    }
    if state.Board[i][j] != -1 {
        return false
    }
    return true
}

func Play(state *State, me int, i int, j int) *State {
    s := state.Clone()
    s.Board[i][j] = me
    s.Turn++
    return s
}

func CreateCandidate(state *State, me int, i int, j int) func() *State {
    return func() *State {
        return Play(state, me, i, j)
    }
}

func GetCandidates(state *State, me int) []func() *State {
    candidates := make([]func() *State, 0)
    for i := 0; i < state.Rows; i++ {
        for j := 0; j < state.Cols; j++ {
            if CanPlay(state, me, i, j) {
                candidates = append(candidates, CreateCandidate(state, me, i, j))
            }
        }
    }
    return candidates
}
