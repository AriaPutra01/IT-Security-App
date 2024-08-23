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

import axios from 'axios';
import {jwtDecode} from 'jwt-decode';

axios.interceptors.request.use(function (config) {
  const token = localStorage.getItem('token');
  config.headers.Authorization =  token ? `Bearer ${token}` : '';
  return config;
});

import { Navigate } from 'react-router-dom';


const ProtectedRoute = ({ children, requiredRole }) => {
  const token = localStorage.getItem('token');
  if (!token) {
    // Redirect ke halaman login jika tidak ada token
    return <Navigate to="/login" />;
  }
  const decoded = jwtDecode(token);
  console.log('Current User Role:', decoded.role); // Log role saat ini
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
  { path: "/register", element: <ProtectedRoute> <RegisterPage /></ProtectedRoute> },
  // dashboard
  { path: "/dashboard", element: <ProtectedRoute><DashboardPage /></ProtectedRoute> },
  // Dokumen
  { path: "/sag", element: <ProtectedRoute><SagPage /></ProtectedRoute> },
  { path: "/iso", element: <ProtectedRoute><IsoPage /></ProtectedRoute> },
  { path: "/memo", element: <ProtectedRoute><MemoPage /></ProtectedRoute> },
  { path: "/surat", element: <ProtectedRoute><SuratPage /></ProtectedRoute> },
  { path: "/berita-acara", element: <ProtectedRoute><BeritaAcaraPage /></ProtectedRoute> },
  { path: "/sk", element: <ProtectedRoute><SkPage /></ProtectedRoute> },
  // Rencana Kerja
  { path: "/project", element: <ProtectedRoute><ProjectPage /></ProtectedRoute> },
  { path: "/base-project", element: <ProtectedRoute><BaseProjectPage /></ProtectedRoute> },
  // Kegiatan Proses
  { path: "/perjalanan-dinas", element: <ProtectedRoute><PerdinPage /></ProtectedRoute> },
  // Data Informasi
  { path: "/surat-masuk", element: <ProtectedRoute><SuratMasukPage /></ProtectedRoute> },
  { path: "/surat-keluar", element: <ProtectedRoute><SuratKeluarPage /></ProtectedRoute> },
]);

// const decoded = jwtDecode(token);
// console.log('Current User Role:', decoded.role); // Log role saat ini

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);