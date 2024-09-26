import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import {
  createBrowserRouter,
  Navigate,
  RouterProvider,
} from "react-router-dom";
// Auth
import { TokenProvider, useToken } from "./context/TokenContext";
import { LoginPage } from "./pages/Auth/LoginPage";
import { RegisterPage } from "./pages/Auth/RegisterPage";
// User
import { UserPage } from "./pages/Services/Users/UserPage";
// Welcome
import { WelcomePage } from "./pages/Welcome/Welcome";
// Dashboard
import ErrorPage from "./pages/Error/404";
import { DashboardPage } from "./pages/Dashboard/Dashboard";
// Dokumen
import { MemoPage } from "./pages/Services/Dokumen/MemoPage";
import { BeritaAcaraPage } from "./pages/Services/Dokumen/BeritaAcaraPage";
import { SkPage } from "./pages/Services/Dokumen/SkPage";
import { SuratPage } from "./pages/Services/Dokumen/SuratPage";
import { PerdinPage } from "./pages/Services/Dokumen/PerjalananDinasPage";
// Rencana Kerja
import { ProjectPage } from "./pages/Services/RencanaKerja/ProjectPage";
import { BaseProjectPage } from "./pages/Services/RencanaKerja/BaseProjectPage";
// Kegiatan Proses
import { TimelineProjectPage } from "./pages/Services/KegiatanProses/TimelineProjectPage";
import { TimelineDesktopPage } from "./pages/Services/KegiatanProses/TimelineDesktopPage";
import { MeetingPage } from "./pages/Services/KegiatanProses/MeetingPage";
import { MeetingListPage } from "./pages/Services/KegiatanProses/MeetingListPage";
import { BookingRapatPage } from "./pages/Services/KegiatanProses/BookingRapatPage";
import { JadwalRapatPage } from "./pages/Services/KegiatanProses/JadwalRapatPage";
import { JadwalCutiPage } from "./pages/Services/KegiatanProses/JadwalCutiPage";
// Data Informasi
import { SuratMasukPage } from "./pages/Services/DataInformasi/SuratMasukPage";
import { SuratKeluarPage } from "./pages/Services/DataInformasi/SuratKeluarPage";
import { ArsipPage } from "./pages/Services/DataInformasi/ArsipPage";

import axios from "axios";
axios.defaults.withCredentials = true; // Izinkan pengiriman cookie

const ProtectedRoute = ({ children, requiredRole }) => {
  const { token, userDetails } = useToken(); // Ambil token dan userDetails dari context

  if (!token) {
    // Redirect ke halaman login jika tidak ada token
    return <Navigate to="/login" />;
  }

  if (requiredRole && userDetails.role !== requiredRole) {
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
    path: "/add-user",
    element: (
      <ProtectedRoute>
        <RegisterPage />
      </ProtectedRoute>
    ),
  },
  //user
  {
    path: "/user",
    element: (
      <ProtectedRoute>
        <UserPage />
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
    path: "/berita-acara",
    element: (
      <ProtectedRoute>
        <BeritaAcaraPage />
      </ProtectedRoute>
    ),  
  },
  {
    path: "/sk",
    element: (
      <ProtectedRoute>
        <SkPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/surat",
    element: (
      <ProtectedRoute>
        <SuratPage />
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
    path: "/timeline-project",
    element: (
      <ProtectedRoute>
        <TimelineProjectPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/timeline-desktop",
    element: (
      <ProtectedRoute>
        <TimelineDesktopPage />
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
  {
    path: "/meeting",
    element: (
      <ProtectedRoute>
        <MeetingPage />
      </ProtectedRoute>
    ),
  },
  {
    path: "/meeting-schedule",  
    element: (
      <ProtectedRoute>
        <MeetingListPage />
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
  {
    path: "/arsip",
    element: (
      <ProtectedRoute>
        <ArsipPage />
      </ProtectedRoute>
    ),
  },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <TokenProvider>
      <RouterProvider router={router} />
    </TokenProvider>
  </React.StrictMode>
);
