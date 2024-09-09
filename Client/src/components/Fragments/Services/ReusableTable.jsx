import React, { useState, useEffect } from "react";
import DataTable from "react-data-table-component";
import { Button } from "flowbite-react";
import { useToken } from "../../../context/TokenContext";
import { Modal } from "flowbite-react"; 
import Swal from "sweetalert2";
import { jwtDecode } from "jwt-decode";
import { FaEdit, FaTrash, FaUpload, FaDownload } from 'react-icons/fa'; // Tambahkan import untuk ikon

const FileUploadModal = ({ isOpen, onClose, onUpload, uploadedFiles, fetchFiles, selectedId }) => { // Tambahkan selectedId sebagai prop
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
            <li key={index} className="flex items-center justify-between mb-2"> {/* Tambahkan mb-2 untuk margin bawah */}
              <span>{file}</span>
              <div className="flex gap-2">
                <Button
                  onClick={() => handleDeleteFile(file)}
                  color="failure"
                >
                  Hapus
                </Button>
                <Button
                  onClick={() => {
                    if (selectedId) {
                      fetch(`http://localhost:8080/download/${file}`)
                        .then(response => {
                          if (response.ok) {
                            return response.blob();
                          }
                          throw new Error('Network response was not ok.');
                        })
                        .then(blob => {
                          const url = window.URL.createObjectURL(blob);
                          const a = document.createElement('a');
                          a.href = url;
                          a.download = file;
                          document.body.appendChild(a);
                          a.click();
                          a.remove();
                        })
                        .catch(error => {
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

export const ReusableTable = (props) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [uploadedFiles, setUploadedFiles] = useState([]);
  const [selectedId, setSelectedId] = useState(null);

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
    showUploadButton, // Tambahkan prop ini
  } = props;

  const fetchUploadedFiles = async (id) => {
    try {
      const response = await fetch(`http://localhost:8080/files/${id}`);
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

  const { token } = useToken(); 
  let userRole = "";
  if (token) {
    const decoded = jwtDecode(token);
    userRole = decoded.role;
  }

  const header = formConfig.fields.map((field) => ({
    name: field.label,
    selector: (row) => row[field.name],
    sortable: true,
  }));

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
            <FaEdit />
          </Button>
          <Button
            className="w-full"
            onClick={() => handleDelete(data.ID)}
            color="failure"
          >
            <FaTrash />
          </Button>
          {showUploadButton && ( // Tampilkan tombol upload hanya jika showUploadButton true
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

  const uploadFileToServer = async (file) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("id", selectedId); // Gunakan selectedId saat mengupload file

    try {
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
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
    } 
    catch (error) {
      Swal.fire({
        icon: "error",
        title: "Oops...",
        text: "Error occurred during the file upload.",
      });
      console.error("Error:", error);
    }
  };

  return (
    <div className="w-full rounded-lg p-2 overflow-auto">
      <FileUploadModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onUpload={uploadFileToServer}
        uploadedFiles={uploadedFiles}
        fetchFiles={() => fetchUploadedFiles(selectedId)}
        selectedId={selectedId} // Tambahkan selectedId saat memanggil FileUploadModal
      />
      <div className="overflow-auto">
        <DataTable
          title={`Tabel ${formConfig.services}`}
          columns={columns}
          data={MainData}
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
