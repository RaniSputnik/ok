const NONE = 0;
const BLACK = 1;
const WHITE = 2;

var Board = function(id) {
    this.id = id;
};

Board.prototype.load = function(cb) {
    // TODO are we already loading?
    this.loading = setTimeout(() => {
        this.loading = false;
        const boardSize = 9;
        const stones = [];
        for (let i = 0; i < boardSize * boardSize; i++) {
            stones.push(0);
        }

        this.data = {
            created_by: 'RaniSputnik',
            created_at: new Date(),
            black: 'RaniSputnik',
            white: 'Davezilla',
            board: {
                size: 9,
                stones: stones,
            },
        }
        cb(null);
    }, 1000);
}

Board.prototype.render = function(element) {
    if (this.loading) {
        return renderLoading(element);
    }

    // Add a random stone, just cause
    if (Math.random() > 0.9) {
        const boardSize = this.data.board.size;
        const x = Math.floor(Math.random() * boardSize);
        const y = Math.floor(Math.random() * boardSize);
        const i = index(boardSize, x,y);
        if (this.data.board.stones[i] == NONE) {
            this.data.board.stones[i] = Math.round(Math.random()) + 1;
        }
    }

    if (!this.refs) {
        this.refs = renderRefs(element, this.data);
    }
    renderStones(this.refs, this.data.board.stones);
}

const index = (size, x, y) => y * size + x;

const renderLoading = (element) => {
    element.innerHTML = 'Loading...';
}

const renderRefs = (element, data) => {
    const table = document.createElement('table');
    table.className = 'ok-board';
    const tbody = document.createElement('tbody');

    const boardSize = data.board.size;
    const refs = [];
    refs.length = boardSize * boardSize;
    for (let y = 0; y < boardSize; y++) {
        const row = document.createElement('tr');
        for (let x = 0; x < boardSize; x++) {
            const td = document.createElement('td');
            const cell = document.createElement('div');
            cell.className = 'cell';
            const i = index(boardSize, x,y);
            td.appendChild(cell)
            row.appendChild(td);
            refs[i] = cell;
        }
        tbody.appendChild(row);
    }
    table.appendChild(tbody);

    element.innerHTML = '';
    element.appendChild(table);

    return refs;
}

const renderStones = (refs, stones) => {
    assert(refs.length == stones.length);

    for (let i = 0; i < stones.length; i++) {
        switch (stones[i]) {
            case BLACK: refs[i].className = 'cell black'; break;
            case WHITE: refs[i].className = 'cell white'; break;
            case NONE: /* Do nothing */ break;
        }
    }
}

const assert = (condition, msg) => {
    if (!condition) {
        if (!msg) {
            msg = 'Failed assertion';
        }
        throw new Error(msg);
    }
}

window.ok = {
    play: function(boardId, element) {
        var b = new Board(boardId);
        b.load(function(err) {
            const animationFrame = () => {
                b.render(element);
                window.requestAnimationFrame(animationFrame);
            }
            window.requestAnimationFrame(animationFrame);
        });
        b.render(element);
    },
};