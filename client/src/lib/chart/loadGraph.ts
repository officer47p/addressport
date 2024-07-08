import {
  select,
  Selection,
  forceSimulation,
  forceManyBody,
  forceLink,
  forceCenter,
  drag,
  zoom,
  SimulationNodeDatum,
  SimulationLinkDatum,
  D3DragEvent,
  ZoomBehavior,
} from "d3";

const colors: string[][] = [
  ["#9D4452", "#E6A6B0", "#BE6B78", "#812836", "#5B0D1A"],
  ["#A76C48", "#F4CAAF", "#C99372", "#884E2A", "#602E0E"],
  ["#2E6B5E", "#719D93", "#498175", "#1B584A", "#093E32"],
  ["#538E3D", "#A6D096", "#75AC61", "#3A7424", "#1F520C"],
];

const MAIN_NODE_SIZE = 40;
const CHILD_NODE_SIZE = 15;
const DEFAULT_DISTANCE = 90;
const MAIN_NODE_DISTANCE = 90;
const MANY_BODY_STRENGTH = -180;

let i = 0;

interface Node extends SimulationNodeDatum {
  id: string;
  size: number;
  color: string;
  showChildren: boolean;
  fx?: number;
  fy?: number;
}

interface Link extends SimulationLinkDatum<Node> {
  source: string;
  target: string;
  distance: number;
  color: string;
}

export let nodes: Node[] = [
  {
    id: "Karachi",
    size: MAIN_NODE_SIZE,
    color: "#BE6B78",
    showChildren: true,
  },
];
export let links: Link[] = [];

export const loadGraph = (): void => {
  const svg: Selection<SVGSVGElement, unknown, HTMLElement, any> =
    select("#container");
  const width = +svg.attr("width");
  const height = +svg.attr("height");
  const centerX = width / 2;
  const centerY = height / 2;

  const childeNodesHashMap: Map<string, Node[]> = new Map();
  const childLinksHashMap: Map<string, Link[]> = new Map();

  // Populate childeNodesHashMap
  childeNodesHashMap.set("Karachi", [
    {
      id: "Clifton",
      size: MAIN_NODE_SIZE,
      color: colors[i++ % colors.length][1],
      showChildren: true,
    },
    {
      id: "Defence",
      size: MAIN_NODE_SIZE,
      color: colors[i++ % colors.length][1],
      showChildren: true,
    },
    {
      id: "Saddar",
      size: MAIN_NODE_SIZE,
      color: colors[i++ % colors.length][1],
      showChildren: true,
    },
    {
      id: "Gulshan",
      size: MAIN_NODE_SIZE,
      color: colors[i++ % colors.length][1],
      showChildren: true,
    },
  ]);

  childeNodesHashMap.set("Clifton", [
    {
      id: "Boat Basin",
      size: MAIN_NODE_SIZE,
      color: childeNodesHashMap.get("Karachi")![0].color,
      showChildren: true,
    },
    {
      id: "Sea View",
      size: CHILD_NODE_SIZE,
      color: childeNodesHashMap.get("Karachi")![0].color,
      showChildren: true,
    },
    // ... other Clifton nodes
  ]);

  // ... Set nodes for Defence, Saddar, Gulshan, and Boat Basin

  // Populate childLinksHashMap
  childLinksHashMap.set("Karachi", [
    {
      source: "Karachi",
      target: "Clifton",
      distance: MAIN_NODE_DISTANCE,
      color: nodes[0].color,
    },
    {
      source: "Karachi",
      target: "Defence",
      distance: MAIN_NODE_DISTANCE,
      color: nodes[0].color,
    },
    {
      source: "Karachi",
      target: "Saddar",
      distance: MAIN_NODE_DISTANCE,
      color: nodes[0].color,
    },
    {
      source: "Karachi",
      target: "Gulshan",
      distance: MAIN_NODE_DISTANCE,
      color: nodes[0].color,
    },
  ]);

  childLinksHashMap.set("Clifton", [
    {
      source: "Clifton",
      target: "Boat Basin",
      distance: DEFAULT_DISTANCE,
      color: childeNodesHashMap.get("Karachi")![0].color,
    },
    {
      source: "Clifton",
      target: "Sea View",
      distance: DEFAULT_DISTANCE,
      color: childeNodesHashMap.get("Karachi")![0].color,
    },
    // ... other Clifton links
  ]);

  // ... Set links for Defence, Saddar, Gulshan, and Boat Basin

  console.log("After population:");
  console.log(
    "childeNodesHashMap keys:",
    Array.from(childeNodesHashMap.keys())
  );
  console.log("childLinksHashMap keys:", Array.from(childLinksHashMap.keys()));

  // Safely add nodes
  const addNodesIfExist = (key: string) => {
    const newNodes = childeNodesHashMap.get(key);
    if (newNodes) {
      nodes.push(...newNodes);
    }
  };

  addNodesIfExist("Karachi");
  addNodesIfExist("Clifton");
  addNodesIfExist("Defence");
  addNodesIfExist("Saddar");
  addNodesIfExist("Gulshan");
  addNodesIfExist("Boat Basin");

  // Safely add links
  const addLinksIfExist = (key: string) => {
    const newLinks = childLinksHashMap.get(key);
    if (newLinks) {
      links.push(...newLinks);
    }
  };

  addLinksIfExist("Karachi");
  addLinksIfExist("Clifton");
  addLinksIfExist("Defence");
  addLinksIfExist("Saddar");
  addLinksIfExist("Gulshan");
  addLinksIfExist("Boat Basin");

  const g: Selection<SVGGElement, unknown, HTMLElement, any> = svg.append("g");

  let simulation = forceSimulation(nodes)
    .force("charge", forceManyBody().strength(MANY_BODY_STRENGTH))
    .force(
      "link",
      forceLink<Node, Link>(links)
        .id((d) => d.id)
        .distance((link) => link.distance)
    )
    .force("center", forceCenter(centerX, centerY));

  let dragInteraction = drag<SVGCircleElement, Node>().on(
    "drag",
    (event: D3DragEvent<SVGCircleElement, Node, Node>, node: Node) => {
      node.fx = event.x;
      node.fy = event.y;
      simulation.alpha(1);
      simulation.restart();
    }
  );

  svg.call(
    zoom<SVGSVGElement, unknown>().on("zoom", (event) => {
      g.attr("transform", event.transform);
    }) as ZoomBehavior<SVGSVGElement, unknown>
  );

  let lines = g
    .selectAll<SVGLineElement, Link>("line")
    .data(links)
    .enter()
    .append("line")
    .attr("stroke", (link) => link.color || "black");

  let circles = g
    .selectAll<SVGCircleElement, Node>("circle")
    .data(nodes)
    .enter()
    .append("circle")
    .attr("fill", (node) => node.color || "gray")
    .attr("r", (node) => node.size)
    .style("cursor", "pointer")
    .call(dragInteraction);

  circles
    .on("mouseover", (event: MouseEvent, node: Node) => {
      // Add tooltip on mouseover
      if (node.id !== "") {
        const tooltip = select("body")
          .append("div")
          .attr("class", "tooltip")
          .style("position", "absolute")
          .style("z-index", "10")
          .style("background-color", "white")
          .style("padding", "10px")
          .style("border", "1px solid #ccc")
          .style("border-radius", "5px")
          .style("font-size", "14px")
          .style("visibility", "visible")
          .html(`Node: ${node.id}`);

        tooltip
          .style("top", `${event.pageY}px`)
          .style("left", `${event.pageX + 10}px`);
      }
    })
    .on("mouseout", () => {
      select(".tooltip").remove();
    });

  circles.on("click", (event: MouseEvent, node: Node) => {
    showHideChildren(node);
  });

  let text = g
    .selectAll<SVGTextElement, Node>("text")
    .data(nodes)
    .enter()
    .append("text")
    .attr("text-anchor", "middle")
    .attr("alignment-baseline", "middle")
    .style("pointer-events", "none")
    .style("font-size", "14px")
    .style("font-weight", "bold")
    .text((node) => node.id);

  simulation.on("tick", () => {
    circles.attr("cx", (node) => node.x!).attr("cy", (node) => node.y!);
    text.attr("x", (node) => node.x!).attr("y", (node) => node.y!);

    lines
      .attr("x1", (link) => (link.source as unknown as Node).x!)
      .attr("y1", (link) => (link.source as unknown as Node).y!)
      .attr("x2", (link) => (link.target as unknown as Node).x!)
      .attr("y2", (link) => (link.target as unknown as Node).y!);
  });

  function showHideChildren(node: Node): void {
    console.log("Called showHideChildren for node:", node.id);
    console.log(
      "childeNodesHashMap keys:",
      Array.from(childeNodesHashMap.keys())
    );
    console.log(
      "childLinksHashMap keys:",
      Array.from(childLinksHashMap.keys())
    );

    const childNodes = childeNodesHashMap.get(node.id);
    const childLinks = childLinksHashMap.get(node.id);

    console.log("childNodes:", childNodes);
    console.log("childLinks:", childLinks);

    if (childNodes && childLinks) {
      if (node.showChildren) {
        const descendantNodes = getDescendantNodes(node.id);
        descendantNodes.forEach((node) => (node.showChildren = false));
        nodes = nodes.filter(
          (n) => !descendantNodes.some((d) => n.id === d.id)
        );
        const descendantLinks = getDescendantLinks(node.id);
        links = links.filter(
          (l) =>
            !descendantLinks.some(
              (d) => l.source === d.source && l.target === d.target
            )
        );
        node.showChildren = false;
      } else {
        nodes.push(...childNodes);
        links.push(...childLinks);
        node.showChildren = true;
      }
      updateGraph();
    }
  }

  function getDescendantNodes(parentId: string): Node[] {
    let descendants: Node[] = [];
    const queue: string[] = [parentId];
    while (queue.length > 0) {
      const currentId = queue.shift()!;
      const children = childeNodesHashMap.get(currentId);
      if (children) {
        descendants.push(...children);
        queue.push(...children.map((child) => child.id));
      }
    }
    return descendants;
  }

  function getDescendantLinks(parentId: string): Link[] {
    let descendants: Link[] = [];
    const queue: string[] = [parentId];
    while (queue.length > 0) {
      const currentId = queue.shift()!;
      const children = childLinksHashMap.get(currentId);
      if (children) {
        descendants.push(...children);
        queue.push(...children.map((child) => child.target));
      }
    }
    return descendants;
  }

  function updateGraph(): void {
    simulation.nodes(nodes);
    simulation.force<any>("link")!.links(links);

    circles = g
      .selectAll<SVGCircleElement, Node>("circle")
      .data(nodes, (node) => node.id);
    circles.exit().remove();
    circles = circles
      .enter()
      .append("circle")
      .attr("fill", (d) => d.color || "gray")
      .attr("r", (d) => d.size)
      .call(dragInteraction)
      .on("click", (event: MouseEvent, clickedNode: Node) => {
        showHideChildren(clickedNode);
      })
      .on("mouseover", (event: MouseEvent, node: Node) => {
        // Add tooltip (same as before)
      })
      .on("mouseout", () => {
        select(".tooltip").remove();
      })
      .merge(circles as any);

    text = g
      .selectAll<SVGTextElement, Node>("text")
      .data(nodes, (node) => node.id);
    text.exit().remove();
    text = text
      .enter()
      .append("text")
      .attr("text-anchor", "middle")
      .attr("alignment-baseline", "middle")
      .style("pointer-events", "none")
      .text((d) => d.id)
      .merge(text as any);

    lines = g.selectAll<SVGLineElement, Link>("line").data(links);
    lines.exit().remove();
    lines = lines
      .enter()
      .append("line")
      .attr("stroke", (link) => link.color || "black")
      .merge(lines as any);

    simulation.alpha(1).restart();
  }
};
