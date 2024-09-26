import React, { useState } from "react";
import App from "../../../components/Layouts/App";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import {
  getUsers,
  updateUser,
  deleteUser,
} from "../../../../API/Users/Users.service";
import { useToken } from "../../../context/TokenContext";

export const UserPage = () => {
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "Username", label: "Username", type: "text", required: true },
      { name: "Email", label: "Email", type: "text", required: true },
      {
        name: "Role",
        label: "Role",
        type: "select",
        options: ["user", "admin"],
        required: true,
      },
    ],
    services: "Users",
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
          get={getUsers}
          CustomHandleAdd={() => {
            window.location.href = "/add-user";
          }}
          update={updateUser}
          remove={deleteUser}
          InfoColumn={false}
        />
        {/* End Table */}
      </div>
    </App>
  );
};
