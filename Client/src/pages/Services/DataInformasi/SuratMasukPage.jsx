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
      { name: "no_surat", label: "No Surat", type: "text", required: true },
      { name: "title", label: "Title Of Letter", type: "text", required: true },
      {
        name: "related_div",
        label: "Related Divisi",
        type: "text",
        required: true,
      },
      { name: "destiny_div", label: "Destiny", type: "text", required: true },
      { name: "tanggal", label: "Date Issue", type: "date", required: true },
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
          importExcel="uploadSuratMasuk"
        />
        {/* End Table */}
      </div>
    </App>
  );
}
