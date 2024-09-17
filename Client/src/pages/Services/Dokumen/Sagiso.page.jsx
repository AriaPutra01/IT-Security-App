import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  addSagiso,
  deleteSagiso,
  getSagisos,
  updateSagiso,
} from "../../../../API/Dokumen/sagiso.service";
import { useToken } from "../../../context/TokenContext";

export function SagisoPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "tanggal", label: "Tanggal", type: "date", required: true },
      {
        name: "no_memo",
        label: "Nomor Memo/Surat",
        type: "text",
        required: true,
      },
      { name: "perihal", label: "Perihal", type: "text", required: true },
      { name: "pic", label: "Pic", type: "text", required: true },
      {
        name: "kategori",
        label: "Kategori",
        type: "select",
        options: ["Sag", "Iso"],
        required: true,
      },
    ],
    services: "Memos",
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
          get={getSagisos}
          set={addSagiso}
          update={updateSagiso}
          remove={deleteSagiso}
          excel
          ExportExcel="exportMemo"
          UpdateExcel="updateMemo"
          importExcel="uploadMemo"
          InfoColumn={true}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
