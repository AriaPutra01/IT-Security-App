import { FormatDate } from "../../../Utilities/FormatDate";
import { ReusableForm } from "../../Fragments/Services/ReusableForm";
import React, { useEffect, useState } from "react";
import { jwtDecode } from "jwt-decode";
import DataTable from "react-data-table-component";
import { Button, Modal } from "flowbite-react";
import { SearchInput } from "../../Elements/SearchInput";
import { useToken } from "../../../context/TokenContext";
import { Excel } from "../../../Utilities/Excel";
import { FaUpload } from "react-icons/fa"; // Tambahkan ikon upload
import Swal from "sweetalert2";

export const ReusableTable = ({
  get,
  set,
  CustomHandleAdd,
  UploadArsip,
  update,
  remove,
  excel,
  ExportExcel,
  UpdateExcel,
  ImportExcel,
  formConfig,
  setFormConfig,
  InfoColumn,
}) => {
  const { token } = useToken(); // Ambil token dari context
  let userRole = "";
  if (token) {
    const decoded = jwtDecode(token);
    userRole = decoded.role;
  }
  const [globalFilterText, setGlobalFilterText] = useState("");

  // File
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [uploadedFiles, setUploadedFiles] = useState([]);
  const [selectedId, setSelectedId] = useState(null);

  const fetchUploadedFiles = async (id) => {
    try {
      const response = await fetch(`http://localhost:8080/files/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        credentials: "include", // Sertakan cookie dalam permintaan
      });
      const data = await response.json();
      if (data.files) {
        setUploadedFiles(data.files);
      } else {
        setUploadedFiles([]); // Set ke array kosong jika tidak ada file
      }
    } catch (error) {
      console.error("Error fetching uploaded files:", error);
    }
  };

  useEffect(() => {
    if (selectedId) {
      fetchUploadedFiles(selectedId);
    }
  }, [selectedId]);

  const uploadFileToServer = async (file) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("id", selectedId); // Gunakan selectedId saat mengupload file
    try {
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
        headers: {
          Authorization: `Bearer ${token}`,
        },
        credentials: "include", // Sertakan cookie dalam permintaan
      });
      if (response.ok) {
        Swal.fire({
          icon: "success",
          title: "File uploaded successfully!",
          showConfirmButton: false,
          timer: 1500,
        });
        fetchUploadedFiles(selectedId);
      } else {
        Swal.fire({
          icon: "error",
          title: "Oops...",
          text: "Error uploading file.",
        });
      }
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Oops...",
        text: "Error occurred during the file upload.",
      });
      console.error("Error:", error);
    }
  };
  // File

  // FORM
  const [MainData, setMainData] = useState([]);
  const [formModalOpen, setFormModalOpen] = useState(false);
  const [formData, setFormData] = useState({});
  const [selectedIds, setSelectedIds] = useState([]);

  useEffect(() => {
    get((data) => {
      // ambil dari API
      setMainData(data || []);
    });
  }, []);

  const onCloseFormModal = () => {
    setFormModalOpen(false);
    setFormData({});
  };

  const handleAdd = () => {
    setFormModalOpen(true);
    setFormConfig((prevConfig) => ({
      ...prevConfig,
      action: "add",
      onSubmit: (data) => AddSubmit(data),
    }));
  };

  const handleEdit = (MainData) => {
    setFormModalOpen(true);
    setFormConfig((prevConfig) => ({
      ...prevConfig,
      action: "edit",
      onSubmit: (data) => EditSubmit(data),
    }));
    setFormData({ ...MainData });
  };

  const AddSubmit = async (data) => {
    // Pastikan data.anggaran adalah number dan tidak undefined
    data.anggaran = data.anggaran ? data.anggaran.toString() : "0";
    try {
      const response = await set(data); // tambah data ke API
      Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data berhasil ditambahkan",
        showConfirmButton: false,
        timer: 1500,
      }).then(() => {
        setMainData([...MainData, response]);
      });
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Gagal!",
        text: "Error saat menyimpan data",
        showConfirmButton: false,
        timer: 1500,
      });
    } finally {
      onCloseFormModal();
    }
  };

  const EditSubmit = async (data) => {
    // Ubah anggaran menjadi string dan pastikan tidak undefined
    data.anggaran = data.anggaran ? data.anggaran.toString() : "0";
    try {
      const response = await update(data.ID, data); // edit data ke API
      Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data berhasil diperbarui",
        showConfirmButton: false,
        timer: 1500,
      }).then(() => {
        setMainData(
          MainData.map((item) => {
            return item.ID === data.ID ? response : item;
          })
        );
      });
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Gagal!",
        text: "Error saat mengubah data",
        showConfirmButton: false,
        timer: 1500,
      });
    } finally {
      onCloseFormModal();
    }
  };
  const handleDelete = async (id) => {
    Swal.fire({
      title: "Apakah Anda yakin?",
      text: "Anda akan menghapus data ini!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          await remove(id);
          setMainData((prevData) => prevData.filter((data) => data.ID !== id));
          Swal.fire({
            icon: "info",
            title: "Berhasil!",
            text: "Data berhasil dihapus",
            showConfirmButton: false,
            timer: 1500,
          });
        } catch (error) {
          Swal.fire("Gagal!", "Error saat hapus data:", "error");
        }
      }
    });
  };

  // handle select
  const handleSelect = ({ selectedRows }) => {
    const id = selectedRows.map((data) => data.ID);
    setSelectedIds(id);
  };

  // Function untuk hapus multi select checkbox
  const handleBulkDelete = async () => {
    Swal.fire({
      title: "Apakah Anda yakin?",
      text: "Anda akan menghapus data yang dipilih!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          await Promise.all(selectedIds.map((id) => remove(id))); // hapus data di API
          setMainData((prevData) =>
            prevData.filter((data) => !selectedIds.includes(data.ID))
          );
          setSelectedIds([]);
          Swal.fire({
            icon: "info",
            title: "Berhasil!",
            text: "Data berhasil dihapus",
            showConfirmButton: false,
            timer: 1500,
          });
        } catch (error) {
          Swal.fire("Gagal!", "Error saat hapus data:", "error");
        }
      }
    });
  };

  // const renderCellContent = (field, value) => {
  //   switch (field.type) {
  //     case "number":
  //       return `Rp. ${new Intl.NumberFormat("id-ID").format(value)}`;
  //     case "date":
  //       return FormatDate(value);
  //     default:
  //       return value;
  //   }
  // };

  const renderCellContent = (field, value) => {
    switch (field.type) {
      case "number":
        return `Rp. ${new Intl.NumberFormat("id-ID").format(value)}`; // Format currency dengan "Rp"
      case "date":
        return FormatDate(value);
      default:
        return value;
    }
  };

  const header = formConfig.fields.map((field) => {
    return {
      name: field.label,
      selector: (row) => {
        return renderCellContent(field, row[field.name]); // Panggil renderCellContent untuk semua field
      },
      sortable: true,
    };
  });

  const columns = [
    ...header,
    ...(InfoColumn
      ? [
          {
            name: "Info", // Menambahkan kolom Info
            selector: (row) => row.info, // Ganti 'info' dengan nama field yang sesuai dari data
            sortable: true,
          },
        ]
      : []), // Tambahkan kolom Info hanya jika showInfoColumn true
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
          {UploadArsip && (
            <Button
              className="w-full"
              onClick={() => {
                setSelectedId(data.ID);
                setIsModalOpen(true);
              }}
              color="info"
            >
              <FaUpload />
            </Button>
          )}
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
  // FORM

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
                onClick={CustomHandleAdd || handleAdd}
                action="add"
                color="info"
              >
                Tambah
              </Button>
              {excel && (
                <Excel
                  linkExportThis={ExportExcel}
                  linkUpdateThis={UpdateExcel}
                  importExcel={ImportExcel}
                />
              )}
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
      {/* ModalFile */}
      <FileUploadModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onUpload={uploadFileToServer}
        uploadedFiles={uploadedFiles}
        fetchFiles={() => fetchUploadedFiles(selectedId)}
        selectedId={selectedId} // Tambahkan selectedId saat memanggil FileUploadModal
      />
      {/* ModalFile */}
      {/* ModalForm */}
      <Modal show={formModalOpen} size="xl" onClose={onCloseFormModal} popup>
        <Modal.Header />
        <Modal.Body>
          <ReusableForm
            config={formConfig}
            formData={formData}
            setFormData={setFormData}
          />
        </Modal.Body>
      </Modal>
      {/* endModalForm */}
    </div>
  );
};

const FileUploadModal = ({
  isOpen,
  onClose,
  onUpload,
  uploadedFiles,
  fetchFiles,
  selectedId,
}) => {
  const { token } = useToken(); // Ambil token dari context
  let userRole = "";
  if (token) {
    const decoded = jwtDecode(token);
    userRole = decoded.role;
  }
  // Tambahkan selectedId sebagai prop
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleSubmit = () => {
    if (file) {
      onUpload(file); // Panggil onUpload dari parent (ReusableTable)
    } else {
      Swal.fire({
        icon: "error",
        title: "Oops...",
        text: "Please select a file before uploading!",
      });
    }
  };

  const handleDeleteFile = async (fileName) => {
    try {
      const response = await fetch(`http://localhost:8080/delete/${fileName}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        credentials: "include", // Sertakan cookie dalam permintaan
      });

      if (response.ok) {
        Swal.fire({
          icon: "success",
          title: "File deleted successfully!",
          showConfirmButton: false,
          timer: 1500,
        });
        fetchFiles(); // Refresh daftar file setelah penghapusan
      } else {
        Swal.fire({
          icon: "error",
          title: "Oops...",
          text: "Error deleting file.",
        });
      }
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Oops...",
        text: "Error occurred during the file deletion.",
      });
      console.error("Error:", error);
    }
  };

  return (
    <Modal show={isOpen} onClose={onClose}>
      <Modal.Header>Upload File</Modal.Header>
      <Modal.Body>
        <input type="file" onChange={handleFileChange} />
        <ul>
          {uploadedFiles.map((file, index) => (
            <li key={index} className="flex items-center justify-between mb-2">
              {/* Tambahkan mb-2 untuk margin bawah */}
              <span>{file}</span>
              <div className="flex gap-2">
                <Button onClick={() => handleDeleteFile(file)} color="failure">
                  Hapus
                </Button>
                <Button
                  onClick={() => {
                    if (selectedId) {
                      fetch(`http://localhost:8080/download/${file}`, {
                        headers: {
                          Authorization: `Bearer ${token}`,
                        },
                        credentials: "include", // Sertakan cookie dalam permintaan
                      })
                        .then((response) => {
                          if (response.ok) {
                            return response.blob();
                          }
                          throw new Error("Network response was not ok.");
                        })
                        .then((blob) => {
                          const url = window.URL.createObjectURL(blob);
                          const a = document.createElement("a");
                          a.href = url;
                          a.download = file;
                          document.body.appendChild(a);
                          a.click();
                          a.remove();
                        })
                        .catch((error) => {
                          Swal.fire({
                            icon: "error",
                            title: "Oops...",
                            text: "Error downloading file.",
                          });
                          console.error("Error:", error);
                        });
                    } else {
                      Swal.fire({
                        icon: "error",
                        title: "Oops...",
                        text: "No file selected for download.",
                      });
                    }
                  }}
                  color="success"
                >
                  Download
                </Button>
              </div>
            </li>
          ))}
        </ul>
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={onClose}>Batal</Button>
        <Button onClick={handleSubmit}>Kirim</Button>
      </Modal.Footer>
    </Modal>
  );
};
