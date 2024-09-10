import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  addMemo,
  deleteMemo,
  getMemos,
  updateMemo,
} from "../../../../API/Dokumen/memo.service";
import { useToken } from "../../../context/TokenContext";

export function MemoPage() {
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
        options: ["sag", "iso", "surat", "berita_acara", "sk"],
        required: true,
      },
    ],
    services: "Memo",
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
          get={getMemos}
          set={addMemo}
          update={updateMemo}
          remove={deleteMemo}
          excel
          ExportExcel="exportMemo"
          UpdateExcel="updateMemo"
          importExcel="uploadMemo"
        />
        {/* End Table */}
      </div>
    </App>
  );
}
