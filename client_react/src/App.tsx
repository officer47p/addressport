import React, { useEffect, useState } from 'react';
import SimpleGraph from './SimpleGraph';
import axios from 'axios';

interface TreeNode {
  address: string;
  children?: TreeNode[];
}

interface Node {
  id: string;
}

interface Link {
  source: string;
  target: string;
}

function processData(data: TreeNode): { nodes: Node[], links: Link[] } {
  const nodes: Node[] = [];
  const links: Link[] = [];

  function traverse(node: TreeNode) {
    // Add current node
    nodes.push({ id: node.address });

    if (node.children && Array.isArray(node.children))     {// Process children
    for (const child of node.children) {
      // Add child node
      nodes.push({ id: child.address });

      // Add link from current node to child
      links.push({ source: node.address, target: child.address });

      // Recursively process child's children
      traverse(child);
    }}
  }

  traverse(data);

  // Remove duplicate nodes
  const uniqueNodes = Array.from(new Set(nodes.map(n => n.id))).map(id => ({ id }));

  return { nodes: uniqueNodes, links };
}


class Data {
  constructor(public nodes: Node[], public links: Link[] ){}
}

const App: React.FC =  () => {
  const [data, setData] = useState<Data|undefined>(undefined);

  useEffect(() => {
    // React advises to declare the async function directly inside useEffect
    async function getAddresses() {
      const data = await axios.get("http://localhost:3001/api/v1/investigation/tools/address-association/0xC1D8E8f14b6AA1cf2F2321348Cbb51d94dc73152?depth=4&#")
      const {nodes, links} = processData(data.data)
      const d = new Data(nodes, links)
      setData(d);
      console.log(data.data)
    }

    // You need to restrict it at some point
    // This is just dummy code and should be replaced by actual
    if (data == undefined) {
      getAddresses();
    }
  }, [data]);


  // const nodes = [
  //   { id: 'A' },
  //   { id: 'B' },
  //   { id: 'C' },
  //   { id: 'D' },
  // ];

  // const links = [
  //   { source: 'A', target: 'B' },
  //   { source: 'B', target: 'C' },
  //   { source: 'C', target: 'D' },
  //   { source: 'D', target: 'A' },
  // ];

  if(!data) {
    return <div></div>
  }
  return (
    <div>
      <h1>Simple Graph Network</h1>
      <SimpleGraph nodes={data.nodes} links={data.links} />
    </div>
  );
};

export default App;