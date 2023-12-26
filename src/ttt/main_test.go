package ttt

import (
    "testing"
)

func TestDiags1x1Hard(t *testing.T) {
    s := InitState(1, 1, 2, 3)
    diags := GetDiags(s)
    if len(diags) != 4 {
        t.Fail()
    }
    if diags[0][0] != -1 || diags[1][0] != -1 || diags[2][0] != -1 || diags[3][0] != -1 {
        t.Fail()
    }
}

func TestDiags3x3(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    diags := GetDiags(s)
    if len(diags) != 12 {
        t.Fail()
    }
    expect := [][]int{
        []int{-1},
        []int{-1, -1},
        []int{-1, -1, -1},
        []int{-1, -1, -1},
        []int{-1, -1},
        []int{-1},
        []int{-1},
        []int{-1, -1},
        []int{-1, -1, -1},
        []int{-1, -1, -1},
        []int{-1, -1},
        []int{-1},
    }
    for i := 0; i < len(expect); i++ {
        if len(diags[i]) != len(expect[i]) {
            t.Fail()
        }
    }
}

func TestDiags3x3Hard(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    s.Board[0][0] = 0
    s.Board[0][1] = 1
    s.Board[0][2] = 2
    s.Board[1][0] = 3
    s.Board[1][1] = 4
    s.Board[1][2] = 5
    s.Board[2][0] = 6
    s.Board[2][1] = 7
    s.Board[2][2] = 8
    diags := GetDiags(s)
    if len(diags) != 12 {
        t.Fail()
    }
    expect := [][]int{
        []int{6},
        []int{3, 7},
        []int{0, 4, 8},
        []int{0, 4, 8},
        []int{1, 5},
        []int{2},
        []int{8},
        []int{5, 7},
        []int{2, 4, 6},
        []int{2, 4, 6},
        []int{1, 3},
        []int{0},
    }
    for i := 0; i < len(expect); i++ {
        if len(diags[i]) != len(expect[i]) {
            t.Fail()
        }
        for j := 0; j < len(expect[i]); j++ {
            if diags[i][j] != expect[i][j] {
                t.Fail()
            }
        }
    }
}

func TestDiags2x3(t *testing.T) {
    s := InitState(2, 3, 2, 3)
    diags := GetDiags(s)
    if len(diags) != 10 {
        t.Fail()
    }
    expect := [][]int{
        []int{-1},
        []int{-1, -1},
        []int{-1, -1},
        []int{-1, -1},
        []int{-1},
        []int{-1},
        []int{-1, -1},
        []int{-1, -1},
        []int{-1, -1},
        []int{-1},
    }
    for i := 0; i < len(expect); i++ {
        if len(diags[i]) != len(expect[i]) {
            t.Errorf("Expected %v, got %v at %d", expect[i], diags[i], i)
        }
    }
}

func TestDiags2x3Hard(t *testing.T) {
    s := InitState(2, 3, 2, 3)
    s.Board[0][0] = 0
    s.Board[0][1] = 1
    s.Board[0][2] = 2
    s.Board[1][0] = 3
    s.Board[1][1] = 4
    s.Board[1][2] = 5
    diags := GetDiags(s)
    if len(diags) != 10 {
        t.Fail()
    }
    expect := [][]int{
        []int{3},
        []int{0, 4},
        []int{0, 4},
        []int{1, 5},
        []int{2},
        []int{5},
        []int{2, 4},
        []int{2, 4},
        []int{1, 3},
        []int{0},
    }
    for i := 0; i < len(expect); i++ {
        if len(diags[i]) != len(expect[i]) {
            t.Errorf("Expected %v, got %v at %d", expect[i], diags[i], i)
        }
        for j := 0; j < len(expect[i]); j++ {
            if diags[i][j] != expect[i][j] {
                t.Errorf("Expected %v, got %v at %d", expect[i], diags[i], i)
            }
        }
    }
}

func TestDiags3x2Hard(t *testing.T) {
    s := InitState(3, 2, 2, 3)
    s.Board[0][0] = 0
    s.Board[0][1] = 1
    s.Board[1][0] = 2
    s.Board[1][1] = 3
    s.Board[2][0] = 4
    s.Board[2][1] = 5
    diags := GetDiags(s)
    if len(diags) != 10 {
        t.Fail()
    }
    expect := [][]int{
        []int{4},
        []int{2, 5},
        []int{0, 3},
        []int{0, 3},
        []int{1},
        []int{5},
        []int{3, 4},
        []int{1, 2},
        []int{1, 2},
        []int{0},
    }
    for i := 0; i < len(expect); i++ {
        if len(diags[i]) != len(expect[i]) {
            t.Errorf("Expected %v, got %v at %d", expect[i], diags[i], i)
        }
        for j := 0; j < len(expect[i]); j++ {
            if diags[i][j] != expect[i][j] {
                t.Errorf("Expected %v, got %v at %d", expect[i], diags[i], i)
            }
        }
    }
}

func TestGetLineWinner(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    line := []int{0, 0, 0}
    if s.GetLineWinner(line) != 0 {
        t.Fail()
    }
    line = []int{1, 1, 0}
    if s.GetLineWinner(line) != -1 {
        t.Fail()
    }
    s = InitState(3, 3, 2, 2)
    line = []int{0, 1, -1}
    if s.GetLineWinner(line) != -1 {
        t.Fail()
    }
    line = []int{0, 1, 1}
    if s.GetLineWinner(line) != 1 {
        t.Fail()
    }
}

func TestGetWinner(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    s.Board[0][0] = 0
    s.Board[0][1] = 1
    s.Board[0][2] = 2
    s.Board[1][0] = 3
    s.Board[1][1] = 4
    s.Board[1][2] = 5
    s.Board[2][0] = 6
    s.Board[2][1] = 7
    s.Board[2][2] = 8
    if w, _ := GetWinner(s); w != -1 {
        t.Fail()
    }
    s.Board[0][0] = 0
    s.Board[0][1] = 0
    s.Board[0][2] = 0
    s.Board[1][0] = 3
    s.Board[1][1] = 4
    s.Board[1][2] = 5
    s.Board[2][0] = 6
    s.Board[2][1] = 7
    s.Board[2][2] = 8
    if w, _ := GetWinner(s); w != 0 {
        t.Fail()
    }
    s.Board[0][0] = 0
    s.Board[0][1] = 1
    s.Board[0][2] = 2
    s.Board[1][0] = 3
    s.Board[1][1] = 1
    s.Board[1][2] = 5
    s.Board[2][0] = 6
    s.Board[2][1] = 1
    s.Board[2][2] = 8
    if w, _ := GetWinner(s); w != 1 {
        t.Fail()
    }
    s.Board[0][0] = 0
    s.Board[0][1] = 1
    s.Board[0][2] = 2
    s.Board[1][0] = 3
    s.Board[1][1] = 0
    s.Board[1][2] = 5
    s.Board[2][0] = 6
    s.Board[2][1] = 7
    s.Board[2][2] = 0
    if w, _ := GetWinner(s); w != 0 {
        t.Fail()
    }
}

func TestCanPlay(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    s.Board[0][0] = 0
    s.Board[0][1] = -1
    s.Board[0][2] = -1
    s.Board[1][0] = 3
    s.Board[1][1] = 4
    s.Board[1][2] = 5
    s.Board[2][0] = 6
    s.Board[2][1] = 7
    s.Board[2][2] = 8
    s.Turn = 1
    if CanPlay(s, 1, 0, 0) {
        t.Fail()
    }
    if !CanPlay(s, 1, 0, 1) {
        t.Fail()
    }
    s.Turn = 0
    if CanPlay(s, 1, 0, 2) {
        t.Fail()
    }
    s.Turn = 1
    if !CanPlay(s, 1, 0, 2) {
        t.Fail()
    }
}

func TestClone(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    s.Board[0][0] = 1
    s2 := s.Clone()
    if s2.Board[0][0] != 1 {
        t.Errorf("Expected %v, got %v", s.Board, s2.Board)
    }
}

func TestGetCandidates(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    fns := GetCandidates(s, 0)
    s2 := fns[0]()
    nZero := 0
    nEmpty := 0
    nOrigEmpty := 0
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if s2.Board[i][j] == 0 {
                nZero++
            } else if s2.Board[i][j] == -1 {
                nEmpty++
            }
            if s.Board[i][j] == -1 {
                nOrigEmpty++
            }
        }
    }
    if nZero != 1 && nEmpty != 8 && nOrigEmpty != 9 {
        t.Fail()
    }
    fns2 := GetCandidates(s2, 0)
    if len(fns2) != 0 {
        t.Fail()
    }
    fns3 := GetCandidates(s2, 1)
    s3 := fns3[0]()
    nZero = 0
    nEmpty = 0
    nOnes := 0
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if s3.Board[i][j] == 0 {
                nZero++
            } else if s3.Board[i][j] == -1 {
                nEmpty++
            } else if s3.Board[i][j] == 1 {
                nOnes++
            }
        }
    }
    if nZero != 1 && nEmpty != 7 && nOnes != 1 {
        t.Fail()
    }
}

func TestClosure(t *testing.T) {
    s := InitState(3, 3, 2, 3)
    fns := GetCandidates(s, 0)
    fns[0]()
    s3 := fns[1]()
    nZero := 0
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if s3.Board[i][j] == 0 {
                nZero++
            }
        }
    }
    if nZero != 1 {
        t.Errorf("%v", s3.Board)
    }
}

func TestGetLineBonus(t *testing.T) {
    s := InitState(4, 3, 2, 3)
    lines := [][]int{
        []int{0,0,1},
    }
    if b := GetLineBonus(s, lines, 0); b > 0 {
        t.Errorf("%v", b)
    }
    lines = [][]int{
        []int{0,0,-1},
    }
    if b := GetLineBonus(s, lines, 0); b == 0 {
        t.Errorf("%v", b)
    }
    lines = [][]int{
        []int{1,0,0},
    }
    if b := GetLineBonus(s, lines, 0); b > 0 {
        t.Errorf("%v", b)
    }
    lines = [][]int{
        []int{-1,0,0},
    }
    if b := GetLineBonus(s, lines, 0); b == 0 {
        t.Errorf("%v", b)
    }
    lines = [][]int{
        []int{1,0,-1},
    }
    if b := GetLineBonus(s, lines, 0); b > 0 {
        t.Errorf("%v", b)
    }
    lines = [][]int{
        []int{-1,0,1},
    }
    if b := GetLineBonus(s, lines, 0); b > 0 {
        t.Errorf("%v", b)
    }
    lines = [][]int{
        []int{-1,0,-1},
    }
    if b := GetLineBonus(s, lines, 0); b > 0 {
        t.Errorf("%v", b)
    }
}
