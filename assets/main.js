import "vite/modulepreload-polyfill";

import "./css/app.css";

import "basecoat-css/all";

// Preact itself lives behind this import, so pages with no islands never
// download it.
if (document.querySelector("[data-island]")) {
  import("./islands/mount.jsx").then(({ mountIslands }) => mountIslands());
}
