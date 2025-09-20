
import 'vite/modulepreload-polyfill'

import "./css/app.css"

// import 'basecoat-css/all';
import Root from './root';

document.addEventListener("DOMContentLoaded", () => {
  const root = document.getElementById("app")
  Root(root!)
})