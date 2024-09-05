import { FormatDate } from "../../../Utilities/FormatDate";
import React, { useState } from "react";
import { jwtDecode } from "jwt-decode";
import DataTable from "react-data-table-component";
import { Button } from "flowbite-react";
import { SearchInput } from "../../Elements/SearchInput";
import { useToken } from "../../../context/TokenContext";

export const ReusableTable = (props) => {
  const {
    excel,
    MainData,
    formConfig,
    handleAdd,
    handleEdit,
    handleDelete,
    handleSelect,
    selectedIds,
    handleBulkDelete,
  } = props;
  const [globalFilterText, setGlobalFilterText] = useState("");
  const { token } = useToken(); // Ambil token dari context
  let userRole = "";
  if (token) {
    const decoded = jwtDecode(token);
    userRole = decoded.role;
  }

  const renderCellContent = (field, value) => {
    switch (field.type) {
      case "number":
        return `Rp. ${new Intl.NumberFormat("id-ID").format(value)}`;
      case "date":
        return FormatDate(value);
      default:
        return value;
    }
  };

  const header = formConfig.fields.map((field) => {
    return {
      name: field.label,
      selector: (row) => row[field.name],
      sortable: true,
    };
  });

  const columns = [
    ...header,
    {
      name: "Action",
      cell: (data) => (
        <div className="flex gap-1">
          <Button
            className="w-full"
            onClick={() => handleEdit(data)}
            action="edit"
            color="warning"
          >
            <svg
              className="w-6 h-6"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="m14.304 4.844 2.852 2.852M7 7H4a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h11a1 1 0 0 0 1-1v-4.5m2.409-9.91a2.017 2.017 0 0 1 0 2.853l-6.844 6.844L8 14l.713-3.565 6.844-6.844a2.015 2.015 0 0 1 2.852 0Z"
              />
            </svg>
          </Button>
          <Button
            className="w-full"
            onClick={() => handleDelete(data.ID)}
            color="failure"
          >
            <svg
              className="w-6 h-6"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                fillRule="evenodd"
                d="M8.586 2.586A2 2 0 0 1 10 2h4a2 2 0 0 1 2 2v2h3a1 1 0 1 1 0 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V8a1 1 0 0 1 0-2h3V4a2 2 0 0 1 .586-1.414ZM10 6h4V4h-4v2Zm1 4a1 1 0 1 0-2 0v8a1 1 0 1 0 2 0v-8Zm4 0a1 1 0 1 0-2 0v8a1 1 0 1 0 2 0v-8Z"
                clipRule="evenodd"
              />
            </svg>
          </Button>
        </div>
      ),
    },
  ];

  const filteredData = MainData.filter((data) => {
    return Object.values(data).some(
      (value) =>
        value !== null &&
        value !== undefined &&
        value.toString().toLowerCase().includes(globalFilterText.toLowerCase())
    );
  });

  return (
    <div className="w-full rounded-lg p-2 overflow-auto">
      <div className="flex justify-between">
        <div className="flex gap-1.5 items-center mx-2 mb-2">
          {userRole === "user" ? (
            <Button
              className="flex justify-center items-center"
              onClick={handleAdd}
              action="add"
              color="info"
            >
              Tambah
            </Button>
          ) : (
            <>
              <Button
                className="flex justify-center items-center"
                onClick={handleAdd}
                action="add"
                color="info"
              >
                Tambah
              </Button>
              {excel}
              <Button
                className="w-max"
                color="failure"
                onClick={handleBulkDelete}
                disabled={selectedIds.length === 0}
              >
                Hapus Data dipilih
              </Button>
            </>
          )}
        </div>
        <SearchInput
          type="text"
          value={globalFilterText}
          onChange={(e) => setGlobalFilterText(e.target.value || "")}
          placeholder="Search..."
        />
      </div>
      <div className="overflow-auto">
        <DataTable
          title={`Tabel ${formConfig.services}`}
          columns={columns}
          data={filteredData}
          onSelectedRowsChange={handleSelect}
          selectableRows
          pagination
          highlightOnHover
          striped
          responsive
          pointerOnHover
          fixedHeader
        />
      </div>
    </div>
  );
};
