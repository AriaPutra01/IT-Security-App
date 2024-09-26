import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  addMemo,
  deleteMemo,
  getMemos,
  updateMemo,
} from "../../../../API/Dokumen/MemoSag.service";
import { useToken } from "../../../context/TokenContext";

export function MemoPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "tanggal", label: "Tanggal", type: "date", required: false },
      {
        name: "no_memo",
        label: "Nomor Memo/Surat",
        type: "select",
        options: ["ITS-SAG", "ITS-ISO"], // Hanya kategori
        required: false,
      },
      { name: "perihal", label: "Perihal", type: "text", required: false },
      { name: "pic", label: "Pic", type: "text", required: false },
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
          ImportExcel="uploadMemo"
          InfoColumn={true}
          UploadArsip={{
            get: "filesMemo",
            upload: "uploadFileMemo",
            download: "downloadMemo",
            delete: "deleteMemo",
          }}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
