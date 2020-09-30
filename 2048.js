const go = new Go();
WebAssembly.instantiateStreaming(fetch("2048.wasm"), go.importObject).then((result) => go.run(result.instance));

// -------------------------
// 2048.js
// -------------------------
let isGameStared = false;

function start() {
    let panel = wasmStart();
    drawPanel(panel);

    isGameStared = true;
}

document.onkeydown = (e) => {
    if (isGameStared && (e.keyCode >= 37 && e.keyCode <= 40)) {
        let panel = wasmMove(e.keyCode);
        drawPanel(panel)
    }
};

function drawPanel(panel) {
    panel.forEach((val ,idx) => {
        let _div = document.getElementById('c'+idx);
        if (val !== 0) {
            _div.innerText = val;
            _div.className = 'cell cell'+val;
        } else {
            _div.innerText = '';
            _div.className = 'cell';
        }
    })
}
