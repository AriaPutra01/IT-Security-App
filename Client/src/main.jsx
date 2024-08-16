import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
// Auth
import { LoginPage, RegisterPage } from "./pages/Auth/AuthPage";
// Welcome
import { WelcomePage } from "./pages/Welcome/Welcome";
// Dashboard
import ErrorPage from "./pages/Error/404";
import DashboardPage from "./pages/Dashboard/Dashboard";
// Dokumen
import { SagPage } from "./pages/Services/Dokumen/SagPage";
import { IsoPage } from "./pages/Services/Dokumen/IsoPage";

const router = createBrowserRouter([
  // welcome
  { path: "/", element: <WelcomePage />, errorElement: <ErrorPage /> },
  // auth
  { path: "/login", element: <LoginPage /> },
  { path: "/register", element: <RegisterPage /> },
  // dashboard
  { path: "/dashboard", element: <DashboardPage /> },
  // Dokumen
  { path: "/sag", element: <SagPage /> },
  { path: "/iso", element: <IsoPage /> },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
