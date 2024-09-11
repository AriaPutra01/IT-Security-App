import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  getArsip,
  addArsip,
  deleteArsip,
  updateArsip,
} from "../../../../API/DataInformasi/Arsip.service";
import { useToken } from "../../../context/TokenContext";

export function ArsipPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "no_arsip", label: "No Arsip", type: "text", required: true },
      {
        name: "jenis_dokumen",
        label: "Jenis Dokumen",
        type: "text",
        required: true,
      },
      { name: "no_dokumen", label: "No Dokumen", type: "text", required: true }, // Diubah dari "From" menjadi "No Dokumen"
      { name: "perihal", label: "Perihal", type: "text", required: true }, // Tambahkan field ini
      { name: "no_box", label: "No Box", type: "text", required: true }, // Tambahkan field ini
      {
        name: "keterangan",
        label: "Keterangan",
        type: "text",
        required: false,
      }, // Tambahkan field ini
      {
        name: "tanggal_dokumen",
        label: "Tanggal Dokumen",
        type: "date",
        required: true,
      }, // Tambahkan field ini
      {
        name: "tanggal_penyerahan",
        label: "Tanggal Penyerahan",
        type: "date",
        required: true,
      }, // Tambahkan field ini
    ],
    services: "Arsip",
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
          get={getArsip}
          set={addArsip}
          update={updateArsip}
          remove={deleteArsip}
          excel
          ExportExcel="exportArsip"
          UpdateExcel="updateArsip"
          importExcel="uploadArsip"
          UploadArsip
        />
        {/* End Table */}
      </div>
    </App>
  );
}
