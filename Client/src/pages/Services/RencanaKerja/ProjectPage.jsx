import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  getProjects,
  addProject,
  updateProject,
  deleteProject,
} from "../../../../API/RencanaKerja/Project.service";
import { useToken } from "../../../context/TokenContext";

export function ProjectPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "bulan", label: "Bulan", type: "date", required: true },
      {
        name: "kode_project",
        label: "Kode Project",
        type: "text",
        required: true,
      },
      {
        name: "jenis_pengadaan",
        label: "Jenis Pengadaan",
        type: "text",
        required: true,
      },
      {
        name: "nama_pengadaan",
        label: "Nama Pengadaan",
        type: "text",
        required: true,
      },
      {
        name: "div_inisiasi",
        label: "Div Inisisasi",
        type: "text",
        required: true,
      },
      {
        name: "sumber_pendanaan",
        label: "Sumber Pendanaan",
        type: "text",
        required: true,
      },
      { name: "anggaran", label: "Anggaran", type: "number", required: true },
      { name: "no_izin", label: "No Izin Prinsip", type: "text", required: true },
      {
        name: "tanggal_izin",
        label: "Tanggal Izin",
        type: "date",
        required: true,
      },
      {
        name: "tanggal_tor",
        label: "Tanggal Tor",
        type: "date",
        required: true,
      },
      { name: "pic", label: "PIC", type: "text", required: true },
    ],
    services: "Project",
  });
  const { token } = useToken(); // Ambil token dari context
  let userRole = "";
  if (token) {
    const decoded = jwtDecode(token);
    userRole = decoded.role;
  }

  return (
    <App services={formConfig.services}>
      <div className="overflow-auto">
        {/* Table */}
        <ReusableTable
          formConfig={formConfig}
          setFormConfig={setFormConfig}
          get={getProjects}
          set={addProject}
          update={updateProject}
          remove={deleteProject}
          excel
          ExportExcel="exportProject"
          UpdateExcel="updateProject"
          importExcel="uploadProject"
        />
        {/* End Table */}
      </div>
    </App>
  );
}
