import { createRoot } from "react-dom/client";
import { Routes } from "@generouted/react-router";
import "./lib/i18n/config"; // Initialize i18n

createRoot(document.getElementById("root")!).render(<Routes />);
