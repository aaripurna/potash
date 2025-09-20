import { BrowserRouter, Routes, Route } from "react-router";
import App from "./pages/app"
import NotFound from "./pages/NotFound";
import ApiContext from "./context/api-context";
import Bir from "./pages/foobar/bir"

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="foobar">
          <Route path="bir" element={<Bir />} />
        </Route>

        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  )
}
