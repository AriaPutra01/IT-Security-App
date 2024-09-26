import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  addMeeting,
  deleteMeeting,
  getMeetings,
  updateMeeting,
} from "../../../../API/KegiatanProses/Meeting.service";
import { useToken } from "../../../context/TokenContext";

export function MeetingPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      {
        name: "task",
        label: "Task",
        type: "text",
        required: true,
      },
      { name: "tindak_lanjut", label: "Tindak Lanjut", type: "text", required: true },
      {
        name: "status",
        label: "Status",
        type: "select",
        options: ["Done", "On Progress", "Cancel"],
        required: true,
      },
      { name: "update_pengerjaan", label: "Update Pengerjaan", type: "text", required: false },
      { name: "pic", label: "Pic", type: "text", required: true },
      { name: "tanggal_target", label: "Tanggal Target", type: "date", required: true },
      { name: "tanggal_actual", label: "Tanggal Actual", type: "date", required: true },
    ],
    services: "Meeting",
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
          get={getMeetings}
          set={addMeeting}
          update={updateMeeting}
          remove={deleteMeeting}
          excel
          ExportExcel="exportMemo"
          UpdateExcel="updateMemo"
          ImportExcel="uploadMemo"
          InfoColumn={true}
          StatusColumn={true}
          UploadArsip={{
            get: "filesMeeting",
            upload: "uploadFileMeeting",
            download: "downloadMeeting",
            delete: "deleteMeeting",
          }}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
