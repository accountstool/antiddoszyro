import { BrowserRouter } from "react-router-dom";

import { AppRouter } from "../routes/AppRouter";
import { AuthProvider } from "./auth";
import { UIProvider } from "./ui";

export default function App() {
  return (
    <BrowserRouter>
      <UIProvider>
        <AuthProvider>
          <AppRouter />
        </AuthProvider>
      </UIProvider>
    </BrowserRouter>
  );
}
