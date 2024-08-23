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
<<<<<<< HEAD
import { PerdinPage } from "./pages/Services/KegiatanProses/PerjalananDinasPage";
=======
import { RuangRapatPage } from "./pages/Services/KegiatanProses/RuangRapatPage";
import { PerdinPage } from "./pages/Services/KegiatanProses/PerjalananDinasPage";
import { JadwalCutiPage } from "./pages/Services/KegiatanProses/JadwalCutiPage";
>>>>>>> 1b0c2169936079c580b8a8ef08b3251fb59f00df
// Data Informasi
import { SuratMasukPage } from "./pages/Services/DataInformasi/SuratMasukPage";
import { SuratKeluarPage } from "./pages/Services/DataInformasi/SuratKeluarPage";

<<<<<<< HEAD
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


=======
>>>>>>> 1b0c2169936079c580b8a8ef08b3251fb59f00df
const router = createBrowserRouter([
  // welcome
  { path: "/", element: <WelcomePage />, errorElement: <ErrorPage /> },
  // auth
  { path: "/login", element: <LoginPage /> },
<<<<<<< HEAD
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

=======
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
  { path: "/ruang-rapat", element: <RuangRapatPage /> },
  { path: "/perjalanan-dinas", element: <PerdinPage /> },
  { path: "/jadwal-cuti", element: <JadwalCutiPage /> },
  // Data Informasi
  { path: "/surat-masuk", element: <SuratMasukPage /> },
  { path: "/surat-keluar", element: <SuratKeluarPage /> },
]);

>>>>>>> 1b0c2169936079c580b8a8ef08b3251fb59f00df
ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
<<<<<<< HEAD
);
=======
);
>>>>>>> 1b0c2169936079c580b8a8ef08b3251fb59f00df
