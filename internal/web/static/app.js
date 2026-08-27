document.addEventListener("keydown", (event) => {
  if (event.key === "/" && document.activeElement?.tagName !== "INPUT") {
    event.preventDefault();
    document.querySelector("#search")?.focus();
  }
});

let treeController = null;

const treePayload = document.querySelector("#tree-data");
const treeCanvas = document.querySelector("[data-tree-canvas]");
if (treePayload && treeCanvas) {
  const treeRoots = JSON.parse(treePayload.textContent).roots;
  const rootByID = new Map(treeRoots.map((root) => [root.id, root]));
  const expandedByTree = new Map();
  let currentTreeID = treeCanvas.dataset.treeId;
  const canvasWrap = treeCanvas.closest(".tree-canvas-wrap");
  let treeZoom = 1;

  const expanded = () => {
    if (!expandedByTree.has(currentTreeID)) expandedByTree.set(currentTreeID, new Set([currentTreeID]));
    return expandedByTree.get(currentTreeID);
  };
  const expandPathToActive = (node, target) => {
    if (!target) return false;
    if (node.card === target || node.id === target) return true;
    for (const child of node.children || []) {
      if (expandPathToActive(child, target)) {
        expanded().add(node.id);
        return true;
      }
    }
    return false;
  };
  const drawTree = () => {
    const root = rootByID.get(currentTreeID) || treeRoots[0];
    currentTreeID = root.id;
    expandPathToActive(root, treeCanvas.dataset.activeId);
    const positions = new Map();
    let leafIndex = 0;
    let maxDepth = 0;
    const visibleChildren = (node) => expanded().has(node.id) ? (node.children || []) : [];
    const layout = (node, depth) => {
      maxDepth = Math.max(maxDepth, depth);
      const children = visibleChildren(node);
      let x;
      if (children.length === 0) {
        x = 72 + leafIndex * 142;
        leafIndex += 1;
      } else {
        const childXs = children.map((child) => layout(child, depth + 1));
        x = childXs.reduce((sum, value) => sum + value, 0) / childXs.length;
      }
      positions.set(node.id, { x, y: 42 + depth * 106, node, children });
      return x;
    };
    layout(root, 0);
    const width = Math.max(330, leafIndex * 142 + 18);
    const height = Math.max(165, (maxDepth + 1) * 106 + 18);
    const scrollLeft = canvasWrap.scrollLeft;
    const scrollTop = canvasWrap.scrollTop;
    treeCanvas.setAttribute("viewBox", `0 0 ${width} ${height}`);
    treeCanvas.setAttribute("height", String(height));
    treeCanvas.setAttribute("width", String(width));
    treeCanvas.style.width = `${width * treeZoom}px`;
    treeCanvas.style.height = `${height * treeZoom}px`;
    treeCanvas.replaceChildren();

    const svg = (tag) => document.createElementNS("http://www.w3.org/2000/svg", tag);
    const drawEdges = (node) => {
      const from = positions.get(node.id);
      for (const child of from.children) {
        const to = positions.get(child.id);
        const line = svg("path");
        const middle = (from.y + to.y) / 2;
        line.setAttribute("d", `M ${from.x} ${from.y + 23} V ${middle} H ${to.x} V ${to.y - 23}`);
        line.setAttribute("class", "tree-edge");
        treeCanvas.append(line);
        drawEdges(child);
      }
    };
    drawEdges(root);
    [...positions.values()].forEach(({ node, x, y, children }) => {
      const group = svg("g");
      const active = node.card && node.card === treeCanvas.dataset.activeId;
      group.setAttribute("class", `tree-node${node.card ? " tree-node--card" : " tree-node--group"}${active ? " is-active" : ""}`);
      group.setAttribute("transform", `translate(${x - 54} ${y - 22})`);
      group.setAttribute("tabindex", node.card || node.children?.length ? "0" : "-1");
      group.setAttribute("aria-label", node.card ? `打开：${node.title}` : node.title);
      if (node.card) group.setAttribute("role", "link");
      const rect = svg("rect");
      rect.setAttribute("width", "108");
      rect.setAttribute("height", "44");
      rect.setAttribute("rx", "8");
      group.append(rect);
      splitTreeLabel(node.title).forEach((line, index, lines) => {
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
        if (node.card) loadCard(node.card);
        else if (node.children?.length) toggleNode(node.id);
      };
      group.addEventListener("click", open);
      group.addEventListener("keydown", (event) => {
        if (event.key === "Enter" || event.key === " ") { event.preventDefault(); open(); }
      });
      treeCanvas.append(group);
      if (node.children?.length) {
        const badge = svg("g");
        badge.setAttribute("class", "tree-fold-badge");
        badge.setAttribute("transform", `translate(${x + 37} ${y - 27})`);
        badge.setAttribute("role", "button");
        badge.setAttribute("tabindex", "0");
        badge.setAttribute("aria-label", `${expanded().has(node.id) ? "收起" : "展开"}${node.title} 的 ${node.children.length} 个子节点`);
        const circle = svg("circle");
        circle.setAttribute("r", "12");
        badge.append(circle);
        const text = svg("text");
        text.setAttribute("text-anchor", "middle");
        text.setAttribute("y", "4");
        text.textContent = expanded().has(node.id) ? "−" : `+${node.children.length}`;
        badge.append(text);
        const toggle = (event) => { event.stopPropagation(); toggleNode(node.id); };
        badge.addEventListener("click", toggle);
        badge.addEventListener("keydown", (event) => {
          if (event.key === "Enter" || event.key === " ") { event.preventDefault(); toggle(event); }
        });
        treeCanvas.append(badge);
      }
    });
    canvasWrap.scrollLeft = scrollLeft;
    canvasWrap.scrollTop = scrollTop;
  };
  const toggleNode = (id) => {
    if (expanded().has(id)) expanded().delete(id);
    else expanded().add(id);
    drawTree();
  };
  const setTreeZoom = (nextZoom) => {
    const next = Math.max(0.8, Math.min(2, nextZoom));
    if (next === treeZoom) return;
    const oldZoom = treeZoom;
    const centerX = (canvasWrap.scrollLeft + canvasWrap.clientWidth / 2) / oldZoom;
    const centerY = (canvasWrap.scrollTop + canvasWrap.clientHeight / 2) / oldZoom;
    treeZoom = next;
    drawTree();
    canvasWrap.scrollLeft = centerX * treeZoom - canvasWrap.clientWidth / 2;
    canvasWrap.scrollTop = centerY * treeZoom - canvasWrap.clientHeight / 2;
  };
  let drag = null;
  const activePointers = new Map();
  let pinch = null;
  let suppressTreeClick = false;
  const finishDrag = (event) => {
    if (!drag || event.pointerId !== drag.pointerID) return;
    if (drag.captured && canvasWrap.hasPointerCapture(event.pointerId)) canvasWrap.releasePointerCapture(event.pointerId);
    if (drag.moved) {
      suppressTreeClick = true;
      window.setTimeout(() => { suppressTreeClick = false; }, 0);
    }
    drag = null;
    canvasWrap.classList.remove("is-dragging");
  };
  canvasWrap.addEventListener("pointerdown", (event) => {
    if (event.button !== 0) return;
    activePointers.set(event.pointerId, { x: event.clientX, y: event.clientY });
    if (activePointers.size === 2) {
      const [first, second] = [...activePointers.values()];
      pinch = { distance: Math.hypot(first.x - second.x, first.y - second.y), zoom: treeZoom };
      drag = null;
      return;
    }
    drag = { pointerID: event.pointerId, startX: event.clientX, startY: event.clientY, scrollLeft: canvasWrap.scrollLeft, scrollTop: canvasWrap.scrollTop, moved: false, captured: false };
  });
  canvasWrap.addEventListener("pointermove", (event) => {
    if (!activePointers.has(event.pointerId)) return;
    activePointers.set(event.pointerId, { x: event.clientX, y: event.clientY });
    if (activePointers.size >= 2) {
      const [first, second] = [...activePointers.values()];
      if (!pinch) pinch = { distance: Math.hypot(first.x - second.x, first.y - second.y), zoom: treeZoom };
      const distance = Math.hypot(first.x - second.x, first.y - second.y);
      if (pinch.distance > 0) {
        const scale = distance / pinch.distance;
        setTreeZoom(pinch.zoom * (1 + (scale - 1) * 3));
      }
      suppressTreeClick = true;
      return;
    }
    if (!drag || event.pointerId !== drag.pointerID) return;
    const deltaX = event.clientX - drag.startX;
    const deltaY = event.clientY - drag.startY;
    if (Math.abs(deltaX) > 4 || Math.abs(deltaY) > 4) {
      drag.moved = true;
      canvasWrap.classList.add("is-dragging");
      if (!drag.captured) {
        canvasWrap.setPointerCapture(event.pointerId);
        drag.captured = true;
      }
    }
    if (drag.moved) {
      canvasWrap.scrollLeft = drag.scrollLeft - deltaX;
      canvasWrap.scrollTop = drag.scrollTop - deltaY;
    }
  });
  const releasePointer = (event) => {
    const wasPinching = pinch !== null;
    activePointers.delete(event.pointerId);
    if (activePointers.size < 2) pinch = null;
    if (wasPinching) {
      suppressTreeClick = true;
      window.setTimeout(() => { suppressTreeClick = false; }, 0);
    }
    finishDrag(event);
  };
  canvasWrap.addEventListener("pointerup", releasePointer);
  canvasWrap.addEventListener("pointercancel", releasePointer);
  canvasWrap.addEventListener("wheel", (event) => {
    if (!event.ctrlKey) return;
    event.preventDefault();
    setTreeZoom(treeZoom * Math.exp(-event.deltaY * 0.006));
  }, { passive: false });
  canvasWrap.addEventListener("dragstart", (event) => event.preventDefault());
  canvasWrap.addEventListener("click", (event) => {
    if (!suppressTreeClick) return;
    event.preventDefault();
    event.stopImmediatePropagation();
    suppressTreeClick = false;
  }, true);
  treeController = {
    currentTree: () => currentTreeID,
    activate: (id) => { treeCanvas.dataset.activeId = id; drawTree(); },
    draw: drawTree,
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
  const focusTree = (active) => {
    document.body.classList.toggle("tree-focus", active);
    const button = document.querySelector("[data-tree-focus]");
    button?.setAttribute("aria-pressed", String(active));
    if (button) button.textContent = active ? "收起知识树" : "展开知识树";
  };
  document.querySelector("[data-tree-focus]")?.addEventListener("click", () => focusTree(!document.body.classList.contains("tree-focus")));
  document.querySelector("[data-tree-focus-close]")?.addEventListener("click", () => focusTree(false));
  drawTree();
}

function splitTreeLabel(label) {
  if (label.length <= 10) return [label];
  return [label.slice(0, 10), `${label.slice(10, 19)}${label.length > 19 ? "…" : ""}`];
}

async function loadCard(id, pushHistory = true) {
  const treeID = treeController?.currentTree() || "";
  const response = await fetch(`/partials/cards/${id}?tree=${encodeURIComponent(treeID)}`);
  if (!response.ok) {
    window.location.href = `/cards/${id}?tree=${encodeURIComponent(treeID)}`;
    return;
  }
  const panel = document.querySelector("[data-content-panel]");
  if (!panel) return;
  panel.outerHTML = await response.text();
  if (pushHistory) window.history.pushState({}, "", `/cards/${id}?tree=${encodeURIComponent(treeID)}`);
  treeController?.activate(id);
  const nextPanel = document.querySelector("[data-content-panel]");
  initTracePlayers(nextPanel);
}

document.addEventListener("click", (event) => {
  const link = event.target.closest('[data-content-panel] a[href^="/cards/"]');
  if (!link || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
  event.preventDefault();
  const id = new URL(link.href).pathname.split("/").pop();
  loadCard(id);
});

window.addEventListener("popstate", () => {
  const match = window.location.pathname.match(/^\/cards\/([^/]+)$/);
  if (match) loadCard(match[1], false);
});

function initTracePlayers(scope = document) {
  scope?.querySelectorAll("[data-trace]").forEach((player) => {
    if (player.dataset.initialized) return;
    player.dataset.initialized = "true";
    fetch(player.dataset.trace)
      .then((response) => response.json())
      .then((trace) => createTracePlayer(player, trace))
      .catch(() => { player.querySelector("[data-narration]").textContent = "无法加载执行轨迹。"; });
  });
}

initTracePlayers();

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
    appendHighlightedGoSyntax(item, line);
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
    renderTraceBoard(board, trace.kind, frame.state);
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

function appendHighlightedGoSyntax(target, line) {
  const tokens = /(\/\/.*$|"(?:\\.|[^"\\])*"|`[^`]*`|\b(?:break|case|continue|default|defer|else|for|func|go|if|range|return|select|switch|type|var)\b|\b(?:append|len|make|max|min)\b|\b\d+\b|:=|==|!=|<=|>=|&&|&|\|\||[+\-*/=<>])/g;
  let offset = 0;
  for (const match of line.matchAll(tokens)) {
    target.append(document.createTextNode(line.slice(offset, match.index)));
    const token = match[0];
    const span = document.createElement("span");
    if (token.startsWith("//")) span.className = "go-comment";
    else if (token.startsWith('"') || token.startsWith("`")) span.className = "go-string";
    else if (/^(break|case|continue|default|defer|else|for|func|go|if|range|return|select|switch|type|var)$/.test(token)) span.className = "go-keyword";
    else if (/^(append|len|make|max|min)$/.test(token)) span.className = "go-builtin";
    else if (/^\d+$/.test(token)) span.className = "go-number";
    else span.className = "go-operator";
    span.textContent = token;
    target.append(span);
    offset = match.index + token.length;
  }
  target.append(document.createTextNode(line.slice(offset)));
}

function renderTraceBoard(board, kind, state) {
  if (kind === "dp-table") return renderDPTable(board, state);
  if (kind === "dp-grid") return renderDPGrid(board, state);
  if (kind === "rolling-dependency") return renderRollingDependency(board, state);
  if (kind === "bitmask-state") return renderBitmaskState(board, state);
  if (kind === "binary-red-blue") return renderRedBlue(board, state);
  return renderIntervals(board, state);
}

function renderIntervals(board, state) {
  board.className = "trace-board trace-board--intervals";
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

function renderDPTable(board, state) {
  board.className = "trace-board trace-board--dp";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = "dp[i]：到达第 i 级的方案数";
  const cells = document.createElement("div");
  cells.className = "dp-cells";
  state.cells.forEach((cell) => {
    const item = document.createElement("div");
    item.className = `dp-cell dp-cell--${cell.state}${cell.dependency ? " is-dependency" : ""}`;
    const key = document.createElement("small");
    key.textContent = `dp[${cell.index}]`;
    const value = document.createElement("strong");
    value.textContent = String(cell.value);
    item.append(key, value);
    cells.append(item);
  });
  board.replaceChildren(heading, cells);
}

function renderRollingDependency(board, state) {
  board.className = "trace-board trace-board--rolling";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = `第 ${state.index} 轮：先读取，再写 current，最后覆盖滚动变量`;
  const cells = document.createElement("div");
  cells.className = "rolling-cells";
  const values = [
    ["previousTwo", state.previousTwo, state.stage === "read" || state.stage === "write" ? "dependency" : "ready"],
    ["previousOne", state.previousOne, state.stage === "read" || state.stage === "write" ? "dependency" : "ready"],
    ["current", state.hasCurrent ? state.current : "—", !state.hasCurrent ? "pending" : state.stage === "write" ? "current" : "ready"],
  ];
  values.forEach(([name, value, status]) => {
    const item = document.createElement("div");
    item.className = `dp-cell dp-cell--${status}${status === "dependency" ? " is-dependency" : ""}`;
    const key = document.createElement("small");
    key.textContent = name;
    const output = document.createElement("strong");
    output.textContent = String(value);
    item.append(key, output);
    cells.append(item);
  });
  board.replaceChildren(heading, cells);
}

function renderDPGrid(board, state) {
  board.className = "trace-board trace-board--grid";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = state.title;
  const table = document.createElement("table");
  table.className = "dp-grid";
  const head = document.createElement("tr");
  head.append(document.createElement("th"));
  state.columns.forEach((column) => {
    const label = document.createElement("th");
    label.textContent = column;
    head.append(label);
  });
  table.append(head);
  const cells = new Map(state.cells.map((cell) => [`${cell.row}:${cell.column}`, cell]));
  state.rows.forEach((row, rowIndex) => {
    const line = document.createElement("tr");
    const label = document.createElement("th");
    label.textContent = row;
    line.append(label);
    state.columns.forEach((_, columnIndex) => {
      const cell = cells.get(`${rowIndex}:${columnIndex}`);
      const value = document.createElement("td");
      value.className = `dp-grid-cell dp-grid-cell--${cell?.state || "pending"}${cell?.dependency ? " is-dependency" : ""}`;
      value.textContent = cell?.state === "unused" ? "·" : String(cell?.value ?? 0);
      line.append(value);
    });
    table.append(line);
  });
  board.replaceChildren(heading, table);
}

function renderBitmaskState(board, state) {
  board.className = "trace-board trace-board--bitmask";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = "mask 的每一位对应一个城市；高亮位已经访问";
  const cities = document.createElement("div");
  cities.className = "bitmask-cities";
  state.names.forEach((name, index) => {
    const city = document.createElement("div");
    const visited = (state.mask & (1 << index)) !== 0;
    city.className = `bitmask-city${visited ? " is-visited" : ""}${state.previousLast === index ? " is-dependency" : ""}${state.last === index ? " is-last" : ""}`;
    city.textContent = `城市 ${name}`;
    cities.append(city);
  });
  const result = document.createElement("p");
  result.className = "bitmask-result";
  result.textContent = `当前累计代价：${state.cost}`;
  board.replaceChildren(heading, cities, result);
}

function renderRedBlue(board, state) {
  board.className = "trace-board trace-board--binary";
  const numbers = document.createElement("div");
  numbers.className = "binary-numbers";
  state.numbers.forEach((number) => {
    const item = document.createElement("span");
    item.textContent = String(number);
    numbers.append(item);
  });
  const range = document.createElement("div");
  range.className = "binary-range";
  const minimum = state.minimum;
  const maximum = state.maximum;
  const toPercent = (value) => ((value - minimum) / (maximum - minimum)) * 100;
  const window = document.createElement("div");
  window.className = "binary-window";
  window.style.left = `${toPercent(state.red)}%`;
  window.style.width = `${Math.max(1, toPercent(state.blue) - toPercent(state.red))}%`;
  range.append(window);
  [["red", state.red], ["mid", state.mid], ["blue", state.blue]].forEach(([name, value]) => {
    if (value < minimum || value > maximum) return;
    if (name === "mid" && (value === state.red || value === state.blue)) return;
    const marker = document.createElement("span");
    marker.className = `binary-marker binary-marker--${name}`;
    marker.style.left = `${toPercent(value)}%`;
    marker.textContent = String(value);
    marker.setAttribute("aria-label", `${name}=${value}`);
    range.append(marker);
  });
  const result = document.createElement("p");
  result.className = `binary-result${state.feasible ? " is-feasible" : " is-infeasible"}`;
  result.textContent = state.mid < 0 ? "区间 (red, blue]：左开、右闭" : `mid=${state.mid}：${state.feasible ? "蓝色可行" : "红色不可行"}，需要 ${state.groups} 组`;
  board.replaceChildren(numbers, range, result);
}
