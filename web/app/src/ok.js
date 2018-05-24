const NONE = 0;
const BLACK = 1;
const WHITE = 2;

var Board = function(id) {
    this.id = id;
};

Board.prototype.load = function(cb) {
    // TODO are we already loading?
    this.loading = new XMLHttpRequest();
    // TODO handle error states
    this.loading.onreadystatechange = () => {
        if (this.loading.readyState == 4 && this.loading.status == 200) {
            const data = JSON.parse(this.loading.responseText);
            this.loading = null;
            console.log('Got API response', data);

            this.data = mapApiResponse(data);
            cb(null);
        }
    };
    this.loading.open('GET', `http://localhost:8080/games/${this.id}`, true);
    this.loading.send();
}

const mapApiResponse = (res) => {
    const boardSize = res.board.size;
    const stones = [];
    for (let i = 0; i < boardSize * boardSize; i++) {
        stones.push(0);
    }

    for (let j = 0; j < res.board.stones; j++) {
        const stone = res.board.stones[j];
        const i = index(boardSize, stone.x, stone.y);
        switch (stone.colour) {
            case "BLACK": stones[i] = BLACK; break;
            case "WHITE": stones[i] = WHITE; break;
        }
    }

    return {
        ...res,
        board: {
            size: res.board.size,
            stones: stones, // Format changed
        },
    };
}

Board.prototype.render = function(element) {
    if (this.loading) {
        return renderLoading(element);
    }

    // Add a random stone, just cause
    if (Math.random() > 0.98) {
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
    const boardSize = data.board.size;

    const table = document.createElement('table');
    table.className = 'ok-board';

    const thead = document.createElement('thead');
    const theadrow = document.createElement('tr');
    const theaddata = document.createElement('th');
    theaddata.colSpan = boardSize;
    theaddata.appendChild(createPlayerNameEl(data.black, BLACK));
    theaddata.appendChild(document.createTextNode(' vs '));
    theaddata.appendChild(createPlayerNameEl(data.white, WHITE));
    theadrow.appendChild(theaddata);
    thead.appendChild(theadrow);
    table.appendChild(thead);
    
    const tbody = document.createElement('tbody');

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

const createPlayerNameEl = (playerName, colour) => {
    const el = document.createElement('a');
    el.href = '#'; // TODO link to player profile
    el.innerHTML = playerName;
    switch (colour) {
        case BLACK: el.className = 'black'; break;
        case WHITE: el.className = 'white'; break;
    }
    return el;
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