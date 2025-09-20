import React from "react";
import { NavLink } from "react-router";


export default function App() {
  return (
    <div>
      <NavLink to="/about" end>
        About
      </NavLink>
    </div>
  )
}