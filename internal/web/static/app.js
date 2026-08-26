document.addEventListener("keydown", (event) => {
  if (event.key === "/" && document.activeElement?.tagName !== "INPUT") {
    event.preventDefault();
    document.querySelector("#search")?.focus();
  }
});

const player = document.querySelector("[data-trace]");
if (player) {
  fetch(player.dataset.trace)
    .then((response) => response.json())
    .then((trace) => createTracePlayer(player, trace))
    .catch(() => {
      player.querySelector("[data-narration]").textContent = "无法加载执行轨迹。";
    });
}

function createTracePlayer(root, trace) {
  let index = 0;
  let timer = null;
  const board = root.querySelector("[data-board]");
  const narration = root.querySelector("[data-narration]");
  const variables = root.querySelector("[data-variables]");
  const code = root.querySelector("[data-code]");
  const count = root.querySelector("[data-step-count]");
  const play = root.querySelector("[data-play]");

  code.replaceChildren(...trace.pseudocode.map((line, lineIndex) => {
    const item = document.createElement("li");
    item.textContent = line;
    item.dataset.line = String(lineIndex);
    return item;
  }));

  const stop = () => {
    if (timer) window.clearInterval(timer);
    timer = null;
    play.textContent = "播放";
  };
  const render = () => {
    const frame = trace.frames[index];
    narration.textContent = frame.narration;
    count.textContent = `${index + 1} / ${trace.frames.length}`;
    variables.replaceChildren(...Object.entries(frame.variables).flatMap(([key, value]) => {
      const term = document.createElement("dt"); term.textContent = key;
      const definition = document.createElement("dd"); definition.textContent = value;
      return [term, definition];
    }));
    code.querySelectorAll("li").forEach((line) => line.classList.toggle("active", Number(line.dataset.line) === frame.activeLine));
    renderIntervals(board, frame.state);
  };
  const setIndex = (next) => {
    index = Math.max(0, Math.min(trace.frames.length - 1, next));
    render();
  };
  root.querySelector("[data-first]").addEventListener("click", () => { stop(); setIndex(0); });
  root.querySelector("[data-prev]").addEventListener("click", () => { stop(); setIndex(index - 1); });
  root.querySelector("[data-next]").addEventListener("click", () => { stop(); setIndex(index + 1); });
  root.querySelector("[data-last]").addEventListener("click", () => { stop(); setIndex(trace.frames.length - 1); });
  play.addEventListener("click", () => {
    if (timer) return stop();
    play.textContent = "暂停";
    timer = window.setInterval(() => {
      if (index >= trace.frames.length - 1) return stop();
      setIndex(index + 1);
    }, 900);
  });
  render();
}

function renderIntervals(board, state) {
  board.replaceChildren(...state.intervals.map((interval) => {
    const row = document.createElement("div");
    row.className = "interval-row";
    const label = document.createElement("span");
    label.textContent = interval.label;
    const track = document.createElement("div");
    track.className = "interval-track";
    const bar = document.createElement("div");
    bar.className = `interval interval--${interval.status}`;
    bar.style.left = `${(interval.start / 9) * 100}%`;
    bar.style.width = `${((interval.end - interval.start) / 9) * 100}%`;
    bar.textContent = `[${interval.start}, ${interval.end})`;
    track.append(bar);
    row.append(label, track);
    return row;
  }));
}
