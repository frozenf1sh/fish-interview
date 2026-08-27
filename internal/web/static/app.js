document.addEventListener("keydown", (event) => {
  if (event.key === "/" && document.activeElement?.tagName !== "INPUT") {
    event.preventDefault();
    document.querySelector("#search")?.focus();
  }
});

const treePayload = document.querySelector("#tree-data");
const treeCanvas = document.querySelector("[data-tree-canvas]");
if (treePayload && treeCanvas) {
  const treeRoots = JSON.parse(treePayload.textContent).roots;
  const rootByID = new Map(treeRoots.map((root) => [root.id, root]));
  let currentTreeID = treeCanvas.dataset.treeId;

  const drawTree = () => {
    const root = rootByID.get(currentTreeID) || treeRoots[0];
    currentTreeID = root.id;
    const positions = new Map();
    let leafIndex = 0;
    let maxDepth = 0;
    const layout = (node, depth) => {
      maxDepth = Math.max(maxDepth, depth);
      const children = node.children || [];
      let y;
      if (children.length === 0) {
        y = 38 + leafIndex * 72;
        leafIndex += 1;
      } else {
        const childYs = children.map((child) => layout(child, depth + 1));
        y = childYs.reduce((sum, value) => sum + value, 0) / childYs.length;
      }
      positions.set(node.id, { x: 65 + depth * 164, y, node });
      return y;
    };
    layout(root, 0);
    const width = Math.max(330, 130 + maxDepth * 164);
    const height = Math.max(150, leafIndex * 72 + 18);
    treeCanvas.setAttribute("viewBox", `0 0 ${width} ${height}`);
    treeCanvas.setAttribute("height", String(height));
    treeCanvas.setAttribute("width", String(width));
    treeCanvas.replaceChildren();

    const svg = (tag) => document.createElementNS("http://www.w3.org/2000/svg", tag);
    const drawEdges = (node) => {
      const from = positions.get(node.id);
      for (const child of node.children || []) {
        const to = positions.get(child.id);
        const line = svg("path");
        const middle = (from.x + to.x) / 2;
        line.setAttribute("d", `M ${from.x + 54} ${from.y} H ${middle} V ${to.y} H ${to.x - 54}`);
        line.setAttribute("class", "tree-edge");
        treeCanvas.append(line);
        drawEdges(child);
      }
    };
    drawEdges(root);
    [...positions.values()].forEach(({ node, x, y }) => {
      const group = svg("g");
      const active = node.card && node.card === treeCanvas.dataset.activeId;
      group.setAttribute("class", `tree-node${node.card ? " tree-node--card" : " tree-node--group"}${active ? " is-active" : ""}`);
      group.setAttribute("transform", `translate(${x - 54} ${y - 22})`);
      group.setAttribute("tabindex", node.card ? "0" : "-1");
      group.setAttribute("aria-label", node.card ? `打开：${node.title}` : node.title);
      if (node.card) group.setAttribute("role", "link");
      const rect = svg("rect");
      rect.setAttribute("width", "108");
      rect.setAttribute("height", "44");
      rect.setAttribute("rx", "8");
      group.append(rect);
      const lines = splitTreeLabel(node.title);
      lines.forEach((line, index) => {
        const label = svg("text");
        label.setAttribute("x", "54");
        label.setAttribute("y", String(lines.length === 1 ? 27 : 19 + index * 13));
        label.setAttribute("text-anchor", "middle");
        label.textContent = line;
        group.append(label);
      });
      const title = svg("title");
      title.textContent = node.title;
      group.append(title);
      const open = () => {
        if (node.card) window.location.href = `/cards/${node.card}?tree=${currentTreeID}`;
      };
      group.addEventListener("click", open);
      group.addEventListener("keydown", (event) => {
        if (event.key === "Enter" || event.key === " ") { event.preventDefault(); open(); }
      });
      treeCanvas.append(group);
    });
  };
  document.querySelectorAll("[data-tree-switch]").forEach((button) => {
    button.addEventListener("click", () => {
      currentTreeID = button.dataset.treeId;
      document.querySelectorAll("[data-tree-switch]").forEach((item) => item.setAttribute("aria-selected", String(item === button)));
      const url = new URL(window.location.href);
      url.searchParams.set("tree", currentTreeID);
      window.history.replaceState({}, "", url);
      drawTree();
    });
  });
  drawTree();
}

function splitTreeLabel(label) {
  if (label.length <= 10) return [label];
  return [label.slice(0, 10), `${label.slice(10, 19)}${label.length > 19 ? "…" : ""}`];
}

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
