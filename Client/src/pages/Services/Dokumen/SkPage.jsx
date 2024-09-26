import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  addSk,
  deleteSk,
  getSks,
  updateSk,
} from "../../../../API/Dokumen/Sk.service";
import { useToken } from "../../../context/TokenContext";

export function SkPage() {
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
    services: "Sk",
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
          get={getSks}
          set={addSk}
          update={updateSk}
          remove={deleteSk}
          excel
          ExportExcel="exportSk"
          UpdateExcel="updateSk"
          ImportExcel="uploadSk"
          InfoColumn={true}
          UploadArsip={{
            get: "filesSk",
            upload: "uploadFileSk",
            download: "downloadSk",
            delete: "deleteSk",
          }}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
