
const $ = q => document.querySelector(q);
const $$ = q => [...document.querySelectorAll(q)];

function drawText(ctx, text, p, color, font, stroke) {
    ctx.save();
    if (font) ctx.font = font;
    const tm = ctx.measureText(text);
    ctx.fillStyle = color;
    if (p.ljust)
        ctx.fillText(text, p.x, p.y);
    else if (p.rjust)
        ctx.fillText(text, p.x-tm.width, p.y);
    else
        ctx.fillText(text, p.x-tm.width/2, p.y);
    if (stroke) {
        ctx.strokeStyle = stroke;
        ctx.lineWidth = 1;
        ctx.strokeText(text, p.x-tm.width/2, p.y);
    }   
    ctx.restore();
    return tm; 
}

class Board {
    constructor(canvas) {
        this.canvas = canvas;
    }

    click(x, y) {
        const dx = this.canvas.width / this.cols;
        const dy = this.canvas.height / this.rows;
        const c = Math.floor(x / dx);
        const r = Math.floor(y / dy);
        if (this.board[r][c] == -1) {
            //this.board[r][c] = this.human;
            //this.repaint();
            const req = JSON.stringify({Key: this.key, Action: "Move", Payload: JSON.stringify([r, c])});
            this.conn.send(req);
        }
    }

    get ctx() {
        return this.canvas.getContext("2d");
    }

    newGame(rows, cols, thresh, nplayers, human, depth, time) {
        this.conn = new WebSocket(`ws://${location.host}/ws`);
        this.nplayers = nplayers;
        this.thresh = thresh;
        this.human = human;
        this.board = [];
        for (let i=0; i<rows; i++) {
            this.board[i] = [];
            for (let j=0; j<cols; j++) {
                this.board[i][j] = -1;
            }
        }
        this.rows = rows;
        this.cols = cols;
        this.conn.onopen = e => {
            const game = JSON.stringify({Key: -1, Human: human, State: {Rows: rows, Cols: cols, Board: this.board, NPlayers: nplayers, WinThresh: thresh, Turn: 0}});
            const req = JSON.stringify({Key: -1, Action: "New", Payload: game, Depth: depth, TimeMillis: time});
            this.conn.send(req);
        }
        this.conn.onmessage = e => {
            const game = JSON.parse(e.data);
            this.board = game.State.Board;
            this.key = game.Key;
            this.repaint();
        }
    }

    repaint() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        if (!this.rows) {
            return;
        }
        const dx = this.canvas.width / this.cols;
        const dy = this.canvas.height / this.rows;
        for (let i=0; i<=this.rows; i++) {
            this.ctx.beginPath();
            this.ctx.moveTo(0, i*dy);
            this.ctx.lineTo(this.canvas.width, i*dy);
            this.ctx.stroke();
        }
        for (let i=0; i<=this.cols; i++) {
            this.ctx.beginPath();
            this.ctx.moveTo(i*dx, 0);
            this.ctx.lineTo(i*dx, this.canvas.height);
            this.ctx.stroke();
        }
        for (let i=0; i<this.rows; i++) {
            for (let j=0; j<this.cols; j++) {
                if (this.board[i][j] != -1) {
                    drawText(this.ctx, this.board[i][j]+1, {x: j*dx+dx/2, y: i*dy+dy/2+10}, "black", "30px serif", "black");
                }
            }
        }
    }
}

window.onload = () => {
    const canvas = $('canvas');
    const board = new Board(canvas);
    canvas.addEventListener("click", e => {
        board.click(e.offsetX, e.offsetY);
    });
    $('#start').addEventListener("click", e => {
        e.preventDefault();
        const rows = parseInt($('#rows').value);
        const cols = parseInt($('#cols').value);
        const thresh = parseInt($('#thresh').value);
        const nplayers = parseInt($('#nplayers').value);
        const player = $('#human').checked ? parseInt($('#player').value)-1 : -1;
        const depth = parseInt($('#depth').value);
        const time = parseInt($('#time').value);
        board.newGame(rows, cols, thresh, nplayers, player, depth, time);
        //board.repaint();
    });
    function updateThresh() {
        const rows = $('#rows').value;
        const cols = $('#cols').value;
        let thresh = $('#thresh').value;
        if (thresh > rows && thresh > cols) {
            thresh = rows > cols ? rows : cols;
            $('#thresh').value = thresh;
        }
    }
    function updatePlayer() {
        const nplayers = $('#nplayers').value;
        const player = $('#player').value;
        if (player > nplayers) {
            $('#player').value = nplayers;
        }
    }
    $('#rows').addEventListener("change", e => {
        updateThresh();
    });
    $('#cols').addEventListener("change", e => {
        updateThresh();
    });
    $('#thresh').addEventListener("change", e => {
        updateThresh();
    });
    $('#nplayers').addEventListener("change", e => {
        updatePlayer();
    });
    $('#player').addEventListener("change", e => {
        updatePlayer();
    });
    $('#time').addEventListener("input", e => {
        $('#time-value').innerText = $('#time').value + " ms";
    });
}
