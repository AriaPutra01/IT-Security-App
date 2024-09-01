import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
// Auth
import { LoginPage } from "./pages/Auth/LoginPage";
import { RegisterPage } from "./pages/Auth/RegisterPage";
// Welcome
import { WelcomePage } from "./pages/Welcome/Welcome";
// Dashboard
import ErrorPage from "./pages/Error/404";
import { DashboardPage } from "./pages/Dashboard/Dashboard";
// Dokumen
import { MemoPage } from "./pages/Services/Dokumen/MemoPage";
import { PerdinPage } from "./pages/Services/Dokumen/PerjalananDinasPage";
// Rencana Kerja
import { ProjectPage } from "./pages/Services/RencanaKerja/ProjectPage";
import { BaseProjectPage } from "./pages/Services/RencanaKerja/BaseProjectPage";
// Kegiatan Proses
import { TimelinePage } from "./pages/Services/KegiatanProses/TimelinePage";
import { BookingRapatPage } from "./pages/Services/KegiatanProses/BookingRapatPage";
import { JadwalRapatPage } from "./pages/Services/KegiatanProses/JadwalRapatPage";
import { JadwalCutiPage } from "./pages/Services/KegiatanProses/JadwalCutiPage";
// Data Informasi
import { SuratMasukPage } from "./pages/Services/DataInformasi/SuratMasukPage";
import { SuratKeluarPage } from "./pages/Services/DataInformasi/SuratKeluarPage";

import axios from "axios";
import { jwtDecode } from "jwt-decode";

axios.interceptors.request.use(function (config) {
  const token = localStorage.getItem("token");
  config.headers.Authorization = token ? `Bearer ${token}` : "";
  return config;
});

import { Navigate } from "react-router-dom";

const ProtectedRoute = ({ children, requiredRole }) => {
  const token = localStorage.getItem("token");
  if (!token) {
    // Redirect ke halaman login jika tidak ada token
    return <Navigate to="/" />;
  }
  const decoded = jwtDecode(token);
  if (requiredRole && decoded.role !== requiredRole) {
    return <Navigate to="/unauthorized" />;
  }
  return children;
};

const router = createBrowserRouter([
  // welcome
  { path: "/", element: <WelcomePage />, errorElement: <ErrorPage /> },
  // auth
  { path: "/login", element: <LoginPage /> },
  {
    path: "/register",
    element: (
      <ProtectedRoute>
        <RegisterPage />
      </ProtectedRoute>
    ),
  },
  // dashboard
  {
    path: "/dashboard",
    element: (
      <ProtectedRoute>
        <DashboardPage />
      </ProtectedRoute>
    ),
  },
  // Dokumen
  {
    path: "/memo",
    element: (
      <ProtectedRoute>
        <MemoPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/perjalanan-dinas",
    element: (
      <ProtectedRoute>
        <PerdinPage />
      </ProtectedRoute>
    ),
  },
  // Rencana Kerja
  {
    path: "/project",
    element: (
      <ProtectedRoute>
        <ProjectPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/base-project",
    element: (
      <ProtectedRoute>
        <BaseProjectPage />
      </ProtectedRoute>
    ),
  },
  // Kegiatan Proses
  {
    path: "/timeline",
    element: (
      <ProtectedRoute>
        <TimelinePage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/booking-rapat",
    element: (
      <ProtectedRoute>
        <BookingRapatPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/jadwal-rapat",
    element: (
      <ProtectedRoute>
        <JadwalRapatPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/jadwal-cuti",
    element: (
      <ProtectedRoute>
        <JadwalCutiPage />
      </ProtectedRoute>
    ),
  },
  // Data Informasi
  {
    path: "/surat-masuk",
    element: (
      <ProtectedRoute>
        <SuratMasukPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/surat-keluar",
    element: (
      <ProtectedRoute>
        <SuratKeluarPage />
      </ProtectedRoute>
    ),
  },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
