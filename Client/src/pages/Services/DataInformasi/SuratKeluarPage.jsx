import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  getSuratKeluars,
  addSuratKeluar,
  deleteSuratKeluar,
  updateSuratKeluar,
} from "../../../../API/DataInformasi/SuratKeluar.service";
import { useToken } from "../../../context/TokenContext";

export function SuratKeluarPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "no_surat", label: "No Surat", type: "text", required: false },
      { name: "title", label: "Title Of Letter", type: "text", required: false },
      { name: "from", label: "From", type: "text", required: false },
      { name: "pic", label: "PIC", type: "text", required: false },
      { name: "tanggal", label: "Date Issue", type: "date", required: false },
    ],
    services: "Surat Keluar",
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
          get={getSuratKeluars}
          set={addSuratKeluar}
          update={updateSuratKeluar}
          remove={deleteSuratKeluar}
          excel
          ExportExcel="exportSuratKeluar"
          UpdateExcel="updateSuratKeluar"
          ImportExcel="uploadSuratKeluar"
          InfoColumn={true}
          UploadArsip={{
            get: "filesSuratKeluar",
            upload: "uploadFileSuratKeluar",
            download: "downloadSuratKeluar",
            delete: "deleteSuratKeluar",
          }}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
