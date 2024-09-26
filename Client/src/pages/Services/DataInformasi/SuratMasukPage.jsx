import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  getSuratMasuks,
  addSuratMasuk,
  deleteSuratMasuk,
  updateSuratMasuk,
} from "../../../../API/DataInformasi/SuratMasuk.service";
import { useToken } from "../../../context/TokenContext";

export function SuratMasukPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "no_surat", label: "No Surat", type: "text", required: false },
      { name: "title", label: "Title Of Letter", type: "text", required: false },
      {
        name: "related_div",
        label: "Related Divisi",
        type: "text",
        required: false,
      },
      { name: "destiny_div", label: "Tujuan", type: "text", required: false },
      { name: "tanggal", label: "Date Issue", type: "date", required: false },
    ],
    services: "Surat Masuk",
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
          get={getSuratMasuks}
          set={addSuratMasuk}
          update={updateSuratMasuk}
          remove={deleteSuratMasuk}
          excel
          ExportExcel="exportSuratMasuk"
          UpdateExcel="updateSuratMasuk"
          ImportExcel="uploadSuratMasuk"
          InfoColumn={true}
          UploadArsip={{
            get: "filesSuratMasuk",
            upload: "uploadFileSuratMasuk",
            download: "downloadSuratMasuk",
            delete: "deleteSuratMasuk",
          }}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
