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
import { MemoPage } from "./pages/Services/Dokumen/MemoPage";
import { SuratPage } from "./pages/Services/Dokumen/SuratPage";
import { BeritaAcaraPage } from "./pages/Services/Dokumen/BeritaAcaraPage";
import { SkPage } from "./pages/Services/Dokumen/SkPage";
// Rencana Kerja
import { ProjectPage } from "./pages/Services/RencanaKerja/ProjectPage";
import { BaseProjectPage } from "./pages/Services/RencanaKerja/BaseProjectPage";
// Kegiatan Proses
import { PerdinPage } from "./pages/Services/KegiatanProses/PerjalananDinasPage";
// Data Informasi
import { SuratMasukPage } from "./pages/Services/DataInformasi/SuratMasukPage";
import { SuratKeluarPage } from "./pages/Services/DataInformasi/SuratKeluarPage";

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
  { path: "/memo", element: <MemoPage /> },
  { path: "/surat", element: <SuratPage /> },
  { path: "/berita-acara", element: <BeritaAcaraPage /> },
  { path: "/sk", element: <SkPage /> },
  // Rencana Kerja
  { path: "/project", element: <ProjectPage /> },
  { path: "/base-project", element: <BaseProjectPage /> },
  // Kegiatan Proses
  { path: "/perjalanan-dinas", element: <PerdinPage /> },
  // Data Informasi
  { path: "/surat-masuk", element: <SuratMasukPage /> },
  { path: "/surat-keluar", element: <SuratKeluarPage /> },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
