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
    loadTracePlayer(player, 0);
  });
}

async function loadTracePlayer(player, attempt) {
  try {
    const response = await fetch(player.dataset.trace, { cache: "no-store" });
    if (!response.ok) throw new Error(`trace request failed: ${response.status}`);
    const trace = await response.json();
    validateTracePayload(trace);
    createTracePlayer(player, trace);
  } catch (error) {
    if (attempt === 0) {
      window.setTimeout(() => loadTracePlayer(player, 1), 120);
      return;
    }
    player.dataset.initialized = "";
    player.querySelector("[data-narration]").textContent = "无法加载执行轨迹；点击此处重试。";
    player.querySelector("[data-narration]").onclick = () => {
      player.dataset.initialized = "true";
      player.querySelector("[data-narration]").onclick = null;
      player.querySelector("[data-narration]").textContent = "正在加载执行轨迹…";
      loadTracePlayer(player, 0);
    };
  }
}

function validateTracePayload(trace) {
  if (!trace || !Array.isArray(trace.pseudocode) || trace.pseudocode.length === 0 || !Array.isArray(trace.frames) || trace.frames.length === 0) {
    throw new Error("trace has no renderable frames");
  }
  trace.frames.forEach((frame, index) => {
    if (!frame || frame.state == null || typeof frame.narration !== "string" || !frame.variables || !Number.isInteger(frame.activeLine) || frame.activeLine < 0 || frame.activeLine >= trace.pseudocode.length) {
      throw new Error(`trace frame ${index} violates player contract`);
    }
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
  if (kind === "flow-steps") return renderFlowSteps(board, state);
  if (kind === "example-state") return renderExampleState(board, state);
  if (kind === "greedy-range" || kind === "window-range") return renderGreedyRange(board, state);
  if (kind === "matrix-state") return renderMatrixState(board, state);
  if (kind === "node-link-state") return renderNodeLinkState(board, state);
  if (kind === "cycle-list-state") return renderCycleListState(board, state);
  if (kind === "linked-list-merge") return renderLinkedListMerge(board, state);
  if (kind === "linked-list-merge-sort") return renderLinkedListMergeSort(board, state);
  if (kind === "sequence-tails") return renderSequenceTails(board, state);
  if (kind === "row-gravity") return renderRowGravity(board, state);
  if (kind === "bitmask-state") return renderBitmaskState(board, state);
  if (kind === "linked-list") return renderLinkedList(board, state);
  if (kind === "binary-red-blue") return renderRedBlue(board, state);
  return renderIntervals(board, state);
}

function renderExampleState(board, state) {
  board.className = "trace-board trace-board--example";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = state.caption;
  const lanes = document.createElement("div");
  lanes.className = "example-lanes";
  state.lanes.forEach((lane) => {
    const row = document.createElement("div");
    row.className = "example-lane";
    const label = document.createElement("small");
    label.textContent = lane.label;
    const tokens = document.createElement("div");
    tokens.className = "example-tokens";
    lane.items.forEach((token) => {
      const item = document.createElement("span");
      item.className = `example-token is-${token.state || "ready"}`;
      item.textContent = token.label;
      tokens.append(item);
    });
    row.append(label, tokens);
    lanes.append(row);
  });
  board.replaceChildren(heading, lanes);
}

function renderGreedyRange(board, state) {
  board.className = "trace-board trace-board--greedy-range";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = state.caption;
  const visual = document.createElement("div");
  visual.className = "greedy-range-visual";
  const scale = (value) => ((value - state.min) / Math.max(1, state.max - state.min)) * 100;
  const axis = document.createElement("div");
  axis.className = "greedy-range-axis";
  [state.min, Math.round((state.min + state.max) / 2), state.max].forEach((value) => {
    const tick = document.createElement("span");
    tick.style.left = `${scale(value)}%`;
    tick.textContent = String(value);
    axis.append(tick);
  });
  visual.append(axis);
  state.tracks.forEach((track) => {
    const row = document.createElement("div");
    row.className = "greedy-range-row";
    const label = document.createElement("small");
    label.textContent = track.label;
    const line = document.createElement("div");
    line.className = "greedy-range-line";
    const stacked = track.segments.length > 1 && track.segments.every((segment) => (segment.kind || "range") === "range");
    if (stacked) line.style.height = `${Math.max(34, track.segments.length * 30 + 4)}px`;
    track.segments.forEach((segment, index) => {
      const bar = document.createElement("span");
      bar.className = `greedy-range-segment is-${segment.state || "ready"} is-${segment.kind || "range"}`;
      bar.style.left = `${scale(segment.start)}%`;
      bar.style.width = `${Math.max(1.5, scale(segment.end) - scale(segment.start))}%`;
      if (stacked) bar.style.top = `${3 + index * 30}px`;
      bar.textContent = segment.label;
      line.append(bar);
    });
    state.markers.filter((marker) => marker.track === track.label).forEach((marker) => {
      const pin = document.createElement("span");
      pin.className = `greedy-range-marker is-${marker.state || "current"}`;
      pin.style.left = `${scale(marker.position)}%`;
      pin.textContent = marker.label;
      line.append(pin);
    });
    row.append(label, line);
    visual.append(row);
  });
  board.replaceChildren(heading, visual);
}

function renderMatrixState(board, state) {
  board.className = "trace-board trace-board--matrix";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = state.caption;
  const grid = document.createElement("div");
  grid.className = "matrix-grid";
  grid.style.gridTemplateColumns = `repeat(${state.columns}, minmax(40px, 1fr))`;
  const byPosition = new Map(state.cells.map((cell) => [`${cell.row}:${cell.column}`, cell]));
  for (let row = 0; row < state.rows; row++) {
    for (let column = 0; column < state.columns; column++) {
      const cell = byPosition.get(`${row}:${column}`) || { label: "", state: "pending" };
      const item = document.createElement("div");
      item.className = `matrix-cell is-${cell.state || "pending"}`;
      item.textContent = cell.label;
      grid.append(item);
    }
  }
  board.replaceChildren(heading, grid);
}

function renderNodeLinkState(board, state) {
  board.className = "trace-board trace-board--node-link";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = state.caption;
  const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
  svg.setAttribute("class", "node-link-canvas");
  svg.setAttribute("viewBox", "0 0 360 210");
  const nodesByID = new Map(state.nodes.map((node) => [node.id, node]));
  const activeLinks = new Set((state.activeLinks || []).map((link) => `${link.from}->${link.to}`));
  state.links.forEach((link) => {
    const from = nodesByID.get(link.from);
    const to = nodesByID.get(link.to);
    if (!from || !to) return;
    const line = document.createElementNS(svg.namespaceURI, "line");
    line.setAttribute("x1", String(from.x)); line.setAttribute("y1", String(from.y));
    line.setAttribute("x2", String(to.x)); line.setAttribute("y2", String(to.y));
    line.setAttribute("class", `node-link-edge${activeLinks.has(`${link.from}->${link.to}`) ? " is-active" : ""}`);
    svg.append(line);
  });
  state.nodes.forEach((node) => {
    const group = document.createElementNS(svg.namespaceURI, "g");
    group.setAttribute("class", `node-link-node is-${node.state || "ready"}`);
    const circle = document.createElementNS(svg.namespaceURI, "circle");
    circle.setAttribute("cx", String(node.x)); circle.setAttribute("cy", String(node.y)); circle.setAttribute("r", "21");
    const text = document.createElementNS(svg.namespaceURI, "text");
    text.setAttribute("x", String(node.x)); text.setAttribute("y", String(node.y + 4)); text.setAttribute("text-anchor", "middle");
    text.textContent = node.label;
    group.append(circle, text); svg.append(group);
  });
  const detail = document.createElement("div");
  detail.className = "node-link-detail";
  if (state.callStack?.length) {
    const stack = document.createElement("div");
    stack.className = "node-link-detail-row";
    stack.append(Object.assign(document.createElement("small"), { textContent: "递归栈" }));
    const stackItems = document.createElement("div");
    stackItems.className = "node-link-chips";
    state.callStack.forEach((entry, index) => {
      const chip = document.createElement("span");
      chip.className = `node-link-chip${index === state.callStack.length - 1 ? " is-current" : " is-dependency"}`;
      chip.textContent = entry;
      stackItems.append(chip);
    });
    stack.append(stackItems);
    detail.append(stack);
  }
  if (state.path?.length) {
    const path = document.createElement("div");
    path.className = "node-link-detail-row";
    path.append(Object.assign(document.createElement("small"), { textContent: "当前路径" }));
    const pathValue = document.createElement("strong");
    pathValue.textContent = state.path.join(" → ");
    path.append(pathValue);
    detail.append(path);
  }
  if (state.values && Object.keys(state.values).length) {
    const values = document.createElement("div");
    values.className = "node-link-values";
    Object.entries(state.values).forEach(([key, value]) => {
      const chip = document.createElement("span");
      chip.textContent = `${key}: ${value}`;
      values.append(chip);
    });
    detail.append(values);
  }
  board.replaceChildren(heading, document.createElement("div"));
  const visual = board.lastChild;
  visual.className = "node-link-visual";
  visual.append(svg, detail);
}

function renderCycleListState(board, state) {
  board.className = "trace-board trace-board--cycle-list";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = state.caption;
  const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
  svg.setAttribute("class", "cycle-list-canvas");
  svg.setAttribute("viewBox", "0 0 360 180");
  const defs = document.createElementNS(svg.namespaceURI, "defs");
  const marker = document.createElementNS(svg.namespaceURI, "marker");
  marker.setAttribute("id", "cycle-arrow"); marker.setAttribute("markerWidth", "8"); marker.setAttribute("markerHeight", "8"); marker.setAttribute("refX", "7"); marker.setAttribute("refY", "4"); marker.setAttribute("orient", "auto");
  const tip = document.createElementNS(svg.namespaceURI, "path");
  tip.setAttribute("d", "M0,0 L8,4 L0,8 z"); tip.setAttribute("fill", "#6f8a77"); marker.append(tip); defs.append(marker); svg.append(defs);
  const nodesByID = new Map(state.nodes.map((node) => [node.id, node]));
  state.links.forEach((link) => {
    const from = nodesByID.get(link.from), to = nodesByID.get(link.to);
    if (!from || !to) return;
    const edge = document.createElementNS(svg.namespaceURI, "line");
    edge.setAttribute("x1", String(from.x)); edge.setAttribute("y1", String(from.y)); edge.setAttribute("x2", String(to.x)); edge.setAttribute("y2", String(to.y)); edge.setAttribute("class", "cycle-list-edge"); edge.setAttribute("marker-end", "url(#cycle-arrow)"); svg.append(edge);
  });
  state.nodes.forEach((node) => {
    const group = document.createElementNS(svg.namespaceURI, "g");
    group.setAttribute("class", `cycle-list-node is-${node.state || "ready"}`);
    const circle = document.createElementNS(svg.namespaceURI, "circle"); circle.setAttribute("cx", String(node.x)); circle.setAttribute("cy", String(node.y)); circle.setAttribute("r", "23");
    const label = document.createElementNS(svg.namespaceURI, "text"); label.setAttribute("x", String(node.x)); label.setAttribute("y", String(node.y + 5)); label.setAttribute("text-anchor", "middle"); label.textContent = node.label;
    group.append(circle, label); svg.append(group);
  });
  const previousPointers = board._cyclePointerPositions || new Map();
  const pointerTargets = new Map();
  Object.entries(state.pointers || {}).forEach(([key, value]) => {
    const node = nodesByID.get(value);
    if (!node) return;
    const sameTarget = Object.entries(state.pointers || {}).filter(([, target]) => target === value).map(([name]) => name);
    const offset = sameTarget.length > 1 ? (sameTarget.indexOf(key) === 0 ? -34 : 34) : -30;
    const pointer = document.createElementNS(svg.namespaceURI, "g");
    pointer.setAttribute("class", `cycle-list-pointer cycle-list-pointer--${key}`);
    const text = document.createElementNS(svg.namespaceURI, "text");
    const x = node.x + (key === "slow" ? -24 : 24);
    const y = node.y + offset;
    text.setAttribute("x", String(x));
    text.setAttribute("y", String(y));
    text.setAttribute("text-anchor", "middle");
    text.textContent = key === "slow" ? "慢" : "快";
    pointer.append(text);
    const previous = previousPointers.get(key);
    if (previous && (previous.x !== x || previous.y !== y)) {
      pointer.style.transform = `translate(${previous.x - x}px, ${previous.y - y}px)`;
      pointer.style.transition = "none";
      requestAnimationFrame(() => {
        pointer.style.transition = "transform 360ms ease";
        pointer.style.transform = "translate(0, 0)";
      });
    }
    pointerTargets.set(key, {x, y});
    svg.append(pointer);
  });
  board._cyclePointerPositions = pointerTargets;
  const pointers = document.createElement("div"); pointers.className = "node-link-values";
  Object.entries(state.pointers || {}).forEach(([key, value]) => { const chip = document.createElement("span"); chip.textContent = `${key === "slow" ? "慢指针" : key === "fast" ? "快指针" : key} → ${value}`; pointers.append(chip); });
  if (state.callStack?.length) { const stack = document.createElement("div"); stack.className = "node-link-stack"; stack.textContent = `循环检查：${state.callStack.join(" → ")}`; pointers.append(stack); }
  board.replaceChildren(heading, svg, pointers);
}

function renderLinkedListMerge(board, state) {
  board.className = "trace-board trace-board--linked-merge";
  const heading = document.createElement("p"); heading.className = "trace-board-label"; heading.textContent = state.caption;
  const renderLane = (label, entries, active) => {
    const row = document.createElement("div"); row.className = "merge-list-row";
    const title = document.createElement("small"); title.textContent = label; row.append(title);
    const items = document.createElement("div"); items.className = "merge-list-items";
    const visibleEntries = label === "结果" ? [{ label: "dummy", state: "ready" }, ...entries] : entries;
    visibleEntries.forEach((entry, index) => { const item = document.createElement("span"); item.className = `example-token is-${entry.state || "ready"}${active === entry.label ? " is-current" : ""}`; item.textContent = entry.label; items.append(item); if (index < visibleEntries.length - 1) items.append(Object.assign(document.createElement("span"), { className: "linked-arrow", textContent: "→" })); });
    row.append(items); return row;
  };
  const result = renderLane("结果", state.result || [], state.tail);
  const inputs = document.createElement("div"); inputs.className = "merge-list-inputs";
  inputs.append(renderLane("A", state.left || [], state.chosen?.startsWith("A:") ? state.chosen.slice(2) : ""), renderLane("B", state.right || [], state.chosen?.startsWith("B:") ? state.chosen.slice(2) : ""));
  const note = document.createElement("p"); note.className = "merge-list-note"; note.textContent = state.chosen ? `本轮接入 ${state.chosen}，tail=${state.tail || "dummy"}` : "dummy 固定结果头，tail 指向结果尾";
  board.replaceChildren(heading, result, inputs, note);
}

function renderLinkedListMergeSort(board, state) {
  board.className = "trace-board trace-board--linked-merge";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = `${state.caption} · ${state.phase || ""}`;
  const renderLane = (label, entries, current = false) => {
    const row = document.createElement("div"); row.className = "merge-list-row";
    const title = document.createElement("small"); title.textContent = label; row.append(title);
    const items = document.createElement("div"); items.className = "merge-list-items";
    const visible = label === "结果" ? [{label: "dummy", state: "ready"}, ...(entries || [])] : (entries || []);
    visible.forEach((entry, index) => {
      const token = document.createElement("span");
      token.className = `example-token is-${entry.state || "ready"}${current && index === 0 ? " is-current" : ""}`;
      token.textContent = entry.label;
      items.append(token);
      if (index < visible.length - 1) items.append(Object.assign(document.createElement("span"), {className: "linked-arrow", textContent: "→"}));
    });
    row.append(items); return row;
  };
  const source = renderLane("当前子链", state.source || [], true);
  const halves = document.createElement("div"); halves.className = "merge-list-inputs";
  halves.append(renderLane("左半", state.left || []), renderLane("右半", state.right || []));
  const result = renderLane("结果", state.result || []);
  const stack = document.createElement("div"); stack.className = "merge-list-note";
  stack.textContent = state.stack?.length ? `递归栈：${state.stack.join("  →  ")}` : "递归栈已清空，返回最终结果";
  board.replaceChildren(heading, source, halves, result, stack);
}

function renderSequenceTails(board, state) {
  board.className = "trace-board trace-board--sequence";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = "橙色是当前数字；tails[k] 是长度 k+1 的递增子序列能取得的最小结尾";
  const numbers = document.createElement("div");
  numbers.className = "sequence-numbers";
  state.numbers.forEach((number, index) => {
    const item = document.createElement("span");
    item.className = `sequence-number${index < state.current ? " is-seen" : ""}${index === state.current ? " is-current" : ""}`;
    item.textContent = String(number);
    numbers.append(item);
  });
  if (Number.isInteger(state.probe) && state.probe >= 0) {
    const probe = document.createElement("p");
    probe.className = `sequence-probe is-${state.probeState || "dependency"}`;
    probe.textContent = `二分检查 tails[${state.probe}]：${state.probeState === "rejected" ? "结尾小于当前值，向右找" : "第一个不小于当前值"}`;
    numbers.append(probe);
  }
  const tailsLabel = document.createElement("p");
  tailsLabel.className = "sequence-tails-label";
  tailsLabel.textContent = "tails";
  const tails = document.createElement("div");
  tails.className = "sequence-tails";
  (state.tails || []).forEach((value, index) => {
    const item = document.createElement("div");
    item.className = `sequence-tail is-${state.tailStates?.[index] || "ready"}`;
    item.textContent = `长度 ${index + 1}: ${value}`;
    tails.append(item);
  });
  board.replaceChildren(heading, numbers, tailsLabel, tails);
}

function renderRowGravity(board, state) {
  board.className = "trace-board trace-board--gravity";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = "# 为可下落块，* 为固定障碍；蓝色是下一个落点";
  const row = document.createElement("div");
  row.className = "gravity-row";
  state.cells.forEach((cell, index) => {
    const item = document.createElement("span");
    item.className = `gravity-cell${cell === "#" ? " is-block" : ""}${cell === "*" ? " is-wall" : ""}${index === state.cursor ? " is-current" : ""}${index === state.write ? " is-write" : ""}`;
    item.textContent = cell;
    row.append(item);
  });
  board.replaceChildren(heading, row);
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

function renderFlowSteps(board, state) {
  board.className = "trace-board trace-board--flow";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = "绿色为已完成步骤，橙色为当前步骤";
  const steps = document.createElement("div");
  steps.className = "flow-steps";
  state.steps.forEach((label, index) => {
    const item = document.createElement("div");
    item.className = `flow-step${index < state.current ? " is-done" : ""}${index === state.current ? " is-current" : ""}`;
    const number = document.createElement("small");
    number.textContent = String(index + 1);
    const text = document.createElement("strong");
    text.textContent = label;
    item.append(number, text);
    steps.append(item);
    if (index < state.steps.length - 1) {
      const arrow = document.createElement("span");
      arrow.className = "flow-arrow";
      arrow.textContent = "→";
      steps.append(arrow);
    }
  });
  board.replaceChildren(heading, steps);
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
  if (state.candidates?.length) {
    const candidates = document.createElement("div");
    candidates.className = "bitmask-candidates";
    state.candidates.forEach((candidate) => { const item = document.createElement("span"); item.textContent = candidate; candidates.append(item); });
    cities.append(candidates);
  }
  if (state.states?.length) {
    const states = document.createElement("div"); states.className = "bitmask-states";
    state.states.forEach((entry) => { const item = document.createElement("span"); item.textContent = entry; states.append(item); });
    cities.append(states);
  }
  const result = document.createElement("p");
  result.className = "bitmask-result";
  result.textContent = `当前累计代价：${state.cost}`;
  board.replaceChildren(heading, cities, result);
}

function renderLinkedList(board, state) {
  const previousPositions = new Map([...board.querySelectorAll("[data-list-value]")].map((node) => [node.dataset.listValue, node.getBoundingClientRect()]));
  board.className = "trace-board trace-board--linked-list";
  const heading = document.createElement("p");
  heading.className = "trace-board-label";
  heading.textContent = "主链上的箭头是当前 Next 指针；橙色节点暂时从主链断开";
  const chain = document.createElement("div");
  chain.className = "linked-chain";
  const renderNode = (value, detached = false) => {
    const node = document.createElement("div");
    const labels = Object.entries(state.pointers || {}).filter(([, target]) => target === value).map(([name]) => name);
    const highlighted = state.highlight?.includes(value);
    node.className = `linked-node${value === "D" ? " is-dummy" : ""}${detached ? " is-detached" : ""}${labels.length ? " is-pointed" : ""}${highlighted ? " is-active" : ""}${!highlighted && !detached ? " is-stable" : ""}`;
    node.dataset.listValue = value;
    const nodeValue = document.createElement("strong");
    nodeValue.textContent = value === "D" ? "dummy" : value;
    node.append(nodeValue);
    if (labels.length) {
      const pointer = document.createElement("small");
      pointer.textContent = labels.join(" · ");
      node.append(pointer);
    }
    return node;
  };
  state.chain.forEach((value, index) => {
    chain.append(renderNode(value));
    if (index < state.chain.length - 1) {
      const arrow = document.createElement("span");
      arrow.className = "linked-arrow";
      arrow.textContent = "→";
      chain.append(arrow);
    }
  });
  if (state.detached?.length) {
    const detached = document.createElement("div");
    detached.className = "linked-detached";
    detached.append(document.createTextNode("暂离主链："));
    state.detached.forEach((value) => detached.append(renderNode(value, true)));
    board.replaceChildren(heading, chain, detached);
  } else {
    board.replaceChildren(heading, chain);
  }
  board.querySelectorAll("[data-list-value]").forEach((node) => {
    const before = previousPositions.get(node.dataset.listValue);
    if (!before) return;
    const after = node.getBoundingClientRect();
    const deltaX = before.left - after.left;
    const deltaY = before.top - after.top;
    if (deltaX === 0 && deltaY === 0) return;
    node.style.transition = "none";
    node.style.transform = `translate(${deltaX}px, ${deltaY}px)`;
    void node.offsetWidth;
    requestAnimationFrame(() => {
      node.style.transition = "transform 300ms ease";
      node.style.transform = "";
    });
  });
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
  const scan = document.createElement("div"); scan.className = "binary-scan";
  (state.scanned || []).forEach((entry) => { const item = document.createElement("span"); item.textContent = entry; scan.append(item); });
  const segments = document.createElement("p"); segments.className = "binary-segments"; segments.textContent = state.segments?.length ? `分组扫描：${state.segments.join(" | ")}` : "";
  board.replaceChildren(numbers, range, scan, segments, result);
}
