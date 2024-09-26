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
      {
        name: "kode_project",
        label: "Kode Project",
        type: "text",
        required: false,
      },
      {
        name: "jenis_pengadaan",
        label: "Jenis Pengadaan",
        type: "text",
        required: false,
      },
      {
        name: "nama_pengadaan",
        label: "Nama Pengadaan",
        type: "text",
        required: false,
      },
      {
        name: "div_inisiasi",
        label: "Div Inisisasi",
        type: "text",
        required: false,
      },
      { name: "bulan", label: "Bulan", type: "date", required: false },
      {
        name: "sumber_pendanaan",
        label: "Sumber Pendanaan",
        type: "text",
        required: false,
      },
      { name: "anggaran", label: "Anggaran", type: "number", required: false },
      { name: "no_izin", label: "No Izin Prinsip", type: "text", required: false },
      {
        name: "tanggal_izin",
        label: "Tanggal Izin",
        type: "date",
        required: false,
      },
      {
        name: "tanggal_tor",
        label: "Tanggal Tor",
        type: "date",
        required: false,
      },
      { name: "pic", label: "PIC", type: "text", required: false },
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
          ImportExcel="uploadProject"
          InfoColumn={true}
          UploadArsip={{
            get: "filesProject",
            upload: "uploadFileProject",
            download: "downloadProject",
            delete: "deleteProject",
          }}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
