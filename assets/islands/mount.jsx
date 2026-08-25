import { render } from "preact";

// Each island is a dynamic import, so a page only downloads the islands it
// actually renders. Vite emits one chunk per entry here.
const islands = {
  "signup-form": () => import("./signup-form.jsx"),
};

function readProps(el) {
  if (!el.dataset.props) return {};

  try {
    return JSON.parse(el.dataset.props);
  } catch (error) {
    console.error(`island "${el.dataset.island}": invalid data-props`, error);
    return {};
  }
}

export function mountIslands(root = document) {
  root.querySelectorAll("[data-island]").forEach(async (el) => {
    const name = el.dataset.island;
    const load = islands[name];

    if (!load) {
      console.error(`island "${name}" is not registered in mount.js`);
      return;
    }

    const { default: Component } = await load();

    // render, not hydrate: Go ships the shell, Preact owns everything inside
    // this node. Anything the island renders that uses basecoat markup gets
    // picked up by basecoat's MutationObserver automatically.
    render(<Component {...readProps(el)} />, el);
  });
}
