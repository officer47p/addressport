import "./style.css";
import { loadGraph } from "./loadGraph";

// import typescriptLogo from './typescript.svg'
// import viteLogo from '/vite.svg'
// import { setupCounter } from "./counter.ts";

document.querySelector<HTMLDivElement>("#app")!.innerHTML = `
  <div className="container">
      <svg id="container" width="1060" height="960" viewBox="0 0 960 960"></svg>
    </div>
`;

loadGraph();
