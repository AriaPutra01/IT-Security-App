import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  addBeritaAcara,
  deleteBeritaAcara,
  getBeritaAcaras,
  updateBeritaAcara,
} from "../../../../API/Dokumen/BeritaAcara.service";
import { useToken } from "../../../context/TokenContext";

export function BeritaAcaraPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "tanggal", label: "Tanggal", type: "date", required: false },
      {
        name: "no_surat",
        label: "Nomor Surat",
        type: "select",
        options: ["ITS-SAG", "ITS-ISO"], // Hanya kategori
        required: false,
      },
      { name: "perihal", label: "Perihal", type: "text", required: false },
      { name: "pic", label: "Pic", type: "text", required: false },
    ],
    services: "Berita Acara",
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
          get={getBeritaAcaras}
          set={addBeritaAcara}
          update={updateBeritaAcara}
          remove={deleteBeritaAcara}
          excel
          ExportExcel="exportBeritaAcara"
          UpdateExcel="updateBeritaAcara"
          ImportExcel="uploadBeritaAcara"
          InfoColumn={true}
          UploadArsip={{
            get: "filesBeritaAcara",
            upload: "uploadFileBeritaAcara",
            download: "downloadBeritaAcara",
            delete: "deleteBeritaAcara",
          }}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
