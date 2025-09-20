import ReactDOM from "react-dom/client";
import AppRoutes from "./js/routes";
import React from "react";

export default function Root(root: HTMLElement) {
  return ReactDOM.createRoot(root).render(<AppRoutes />);
}