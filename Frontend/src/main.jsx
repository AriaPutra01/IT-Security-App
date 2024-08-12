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
import { SagPage } from "./pages/Services/Dokumen";
import {
  IsoPage,
  MemoPage,
  SuratPage,
  BeritaAcaraPage,
  SkPage,
} from "./pages/Services/Dokumen";
// Rencana Kerja
import { ProjectPage, BaseProjectPage } from "./pages/Services/RencanaKerja";
// Kegiatan & Proses
import {
  JadwalCutiPage,
  PerjalananDinasPage,
  RuangRapatPage,
} from "./pages/Services/KegiatanProses";
// Data & Informasi
import {
  SuratMasukPage,
  SuratKeluarPage,
} from "./pages/Services/DataInformasi";

const router = createBrowserRouter([
  // welcome
  { path: "/", element: <WelcomePage />, errorElement: <ErrorPage /> },
  // auth
  { path: "/login", element: <LoginPage /> },
  { path: "/register", element: <RegisterPage /> },
  // dashboard
  { path: "/dashboard", element: <DashboardPage /> },
  // Services Pages
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
  // Kegiatan & Proses
  { path: "/ruang-rapat", element: <RuangRapatPage /> },
  { path: "/perjalanan-dinas", element: <PerjalananDinasPage /> },
  { path: "/jadwal-cuti", element: <JadwalCutiPage /> },
  // Data Informasi
  { path: "/surat-masuk", element: <SuratMasukPage /> },
  { path: "/surat-keluar", element: <SuratKeluarPage /> },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
