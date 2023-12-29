package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    ttt "github.com/aorliche/tttai/ttt"

    "github.com/gorilla/websocket"
)

type Game struct {
    Key int
    Human int
    State *ttt.State
    sendChans []chan *ttt.State
    recvChan chan *ttt.State
    over bool
}

type Request struct {
    Key int
    Action string
    Payload string
    Depth int
    TimeMillis int
    NTop int
}

var games = make(map[int]*Game)
var upgrader = websocket.Upgrader{} // Default options

func NextGameIdx() int {
    max := -1
    for key := range games {
        if key > max {
            max = key
        }
    }
    return max+1
}

func GameLoop(game *Game, conn *websocket.Conn) {
    for {
        for i := 0; i < game.State.NPlayers; i++ {
            if game.sendChans[i] != nil {
                game.sendChans[i] <- game.State
            }
        }
        if ttt.GameOver(game.State) {
            SendGame(game, conn)
            CloseGame(game)
            break
        }
        game.State = <- game.recvChan
        SendGame(game, conn)
        time.Sleep(500 * time.Millisecond)
        if game.over {
            break
        }
    }
}

// Need to return result of fn() for side effects
func TryMove(game *Game, cand *ttt.State) *ttt.State {
    fns := ttt.GetCandidates(game.State, game.Human)
    for _, fn := range fns {
        checkState := fn()
        equal := true
        outer:
        for i := 0; i < game.State.Rows; i++ {
            for j := 0; j < game.State.Cols; j++ {
                if checkState.Board[i][j] != cand.Board[i][j] {
                    equal = false;
                    break outer;
                }
            }
        }
        if equal {
            return checkState
        }
    }
    return nil
}

func SendGame(game *Game, conn *websocket.Conn) {
    jsn, err := json.Marshal(game)
    if err != nil {
        log.Println(err)
        return
    }
    conn.WriteMessage(websocket.TextMessage, jsn)
}

func CloseGame(game *Game) {
    log.Println("game closed")
    game.over = true
    for _, chn := range game.sendChans {
        chn <- nil
        close(chn)
    }
    close(game.recvChan)
}

func Socket(w http.ResponseWriter, r *http.Request) {
    var game *Game
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            if game != nil {
                CloseGame(game)
                game = nil
            }
            return  
        }
        // Do we ever get any other types of messages?
        if msgType != websocket.TextMessage {
            log.Println("Not a text message")
            return
        }
        var req Request
        json.NewDecoder(bytes.NewBuffer(msg)).Decode(&req)
        switch req.Action {
            case "New":  
                log.Println(req.Payload)
                err := json.NewDecoder(bytes.NewBuffer([]byte(req.Payload))).Decode(&game)
                if err != nil {
                    log.Println(err)
                    continue
                }
                game.Key = NextGameIdx()
                game.State = ttt.InitState(game.State.Rows, game.State.Cols, game.State.NPlayers, game.State.WinThresh)
                game.sendChans = make([]chan *ttt.State, game.State.NPlayers)
                game.recvChan = make(chan *ttt.State)
                // Launch AI players
                for i := 0; i<game.State.NPlayers; i++ {
                    if game.Human == -1 || i != game.Human {
                        game.sendChans[i] = make(chan *ttt.State)
                        go ttt.Loop(i, game.sendChans[i], game.recvChan, req.Depth, req.TimeMillis, req.NTop)
                    }
                }
                go GameLoop(game, conn)
                games[game.Key] = game
                // Must send first if human player goes first
                SendGame(game, conn)
            case "Move":
                game := games[req.Key]
                move := make([]int, 2)
                err := json.NewDecoder(bytes.NewBuffer([]byte(req.Payload))).Decode(&move)
                if err != nil {
                    log.Println(err)
                    continue
                }
                state := game.State.Clone()
                state.Board[move[0]][move[1]] = game.Human
                goodState := TryMove(game, state)
                if goodState == nil {
                    log.Println("Move not allowed")
                    continue
                }
                game.recvChan <- goodState
        }
    }
}

type HFunc func (http.ResponseWriter, *http.Request)

func Headers(fn HFunc) HFunc {
    return func (w http.ResponseWriter, req *http.Request) {
        //fmt.Println(req.Method)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        fn(w, req)
    }
}
func ServeStatic(w http.ResponseWriter, req *http.Request, file string) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    http.ServeFile(w, req, file)
}

func ServeLocalFiles(dirs []string) {
    for _, dirName := range dirs {
        fsDir := "../../../static/" + dirName
        dir, err := os.Open(fsDir)
        if err != nil {
            log.Fatal(err)
        }
        files, err := dir.Readdir(0)
        if err != nil {
            log.Fatal(err)
        }
        for _, v := range files {
            log.Println(v.Name(), v.IsDir())
            if v.IsDir() {
                continue
            }
            reqFile := dirName + "/" + v.Name()
            file := fsDir + "/" + v.Name()
            if reqFile == "/index.html" {
                reqFile = "/"
            }
            http.HandleFunc(reqFile, Headers(func (w http.ResponseWriter, req *http.Request) {ServeStatic(w, req, file)}))
        }
    }
}

func main() {
    log.SetFlags(0)
    ServeLocalFiles([]string{"", "/js", "/css"})
    http.HandleFunc("/ws", Socket)
    log.Fatal(http.ListenAndServe(":8002", nil))
}

