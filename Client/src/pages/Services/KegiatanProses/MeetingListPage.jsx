import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  addMeetingList,
  deleteMeetingList,
  getMeetingList,
  updateMeetingList,
} from "../../../../API/KegiatanProses/MeetingList.service";
import { useToken } from "../../../context/TokenContext";

export function MeetingListPage() {
  const [formConfig, setFormConfig] = useState({
    fields: [
      {
        name: "hari",
        label: "Hari",
        type: "text",
        required: true,
      },
      { name: "tanggal", label: "Tanggal", type: "date", required: true },
      { name: "perihal", label: "Perihal", type: "text", required: true },
      {
        name: "waktu",
        label: "Waktu",
        type: "time",
        required: true,
      },
      { name: "tempat", label: "Tempat", type: "text", required: true },
      { name: "pic", label: "Pic", type: "text", required: true },
    ],
    services: "MeetingList",
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
          get={getMeetingList}
          set={addMeetingList}
          update={updateMeetingList}
          remove={deleteMeetingList}
          excel
          ExportExcel="exportMeetingList"
          UpdateExcel="updateMeetingList"
          importExcel="uploadMeetingList"
          InfoColumn={true}
        />
        {/* End Table */}
      </div>
    </App>
  );
}
