import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';

interface Node extends d3.SimulationNodeDatum {
  id: string;
}

interface Link extends d3.SimulationLinkDatum<Node> {
  source: string;
  target: string;
}

interface Props {
  nodes: Node[];
  links: Link[];
}

const SimpleGraph: React.FC<Props> = ({ nodes, links }) => {
  const svgRef = useRef<SVGSVGElement | null>(null);

  useEffect(() => {
    if (!svgRef.current) return;

    const width = 500;
    const height = 500;

    const svg = d3.select(svgRef.current);
    svg.selectAll("*").remove();

    const simulation = d3.forceSimulation<Node>(nodes)
      .force("link", d3.forceLink<Node, Link>(links).id(d => d.id))
      .force("charge", d3.forceManyBody().strength(-50))
      .force("center", d3.forceCenter(width / 2, height / 2));

    const link = svg.append("g")
      .selectAll<SVGLineElement, Link>("line")
      .data(links)
      .join("line")
      .attr("stroke", "#999")
      .attr("stroke-opacity", 0.6);

    const node = svg.append("g")
      .selectAll<SVGCircleElement, Node>("circle")
      .data(nodes)
      .join("circle")
      .attr("r", 5)
      .attr("fill", "blue");

    const label = svg.append("g")
      .selectAll<SVGTextElement, Node>("text")
      .data(nodes)
      .join("text")
      .text(d => d.id)
      .attr("font-size", 0)
      .attr("dx", 8)
      .attr("dy", 4);

    simulation.on("tick", () => {
      link
        .attr("x1", d => (d.source as unknown as Node).x as number)
        .attr("y1", d => (d.source as unknown as Node).y as number)
        .attr("x2", d => (d.target as unknown as Node).x as number)
        .attr("y2", d => (d.target as unknown as Node).y as number);

      node
        .attr("cx", d => d.x as number)
        .attr("cy", d => d.y as number);

      label
        .attr("x", d => d.x as number)
        .attr("y", d => d.y as number);
    });

  }, [nodes, links]);

  return <svg ref={svgRef} width={500} height={500}></svg>;
};

export default SimpleGraph;