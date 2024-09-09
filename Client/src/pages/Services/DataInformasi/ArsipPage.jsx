import React, { useState, useEffect } from "react";
import App from "../../../components/Layouts/App";
import Swal from "sweetalert2";
import { ReusableForm } from "../../../components/Fragments/Services/ReusableForm";
import { ReusableTable } from "../../../components/Fragments/Services/ReusableTable";
import { jwtDecode } from "jwt-decode";
import { Modal, Button } from "flowbite-react"; // Tambahkan Button di sini
import {
  getArsip,
  addArsip,
  deleteArsip,
  updateArsip,
} from "../../../../API/DataInformasi/Arsip.service";
import { useToken } from "../../../context/TokenContext";
import { Excel } from "../../../Utilities/Excel";
import { FaUpload } from "react-icons/fa"; // Tambahkan ikon upload

export function ArsipPage() {
  const [MainData, setMainData] = useState([]);
  const [formModalOpen, setFormModalOpen] = useState(false);
  const [formData, setFormData] = useState({});
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "no_arsip", label: "No Arsip", type: "text", required: true },
      { name: "jenis_dokumen", label: "Jenis Dokumen", type: "text", required: true },
      { name: "no_dokumen", label: "No Dokumen", type: "text", required: true }, // Diubah dari "From" menjadi "No Dokumen"
      { name: "perihal", label: "Perihal", type: "text", required: true }, // Tambahkan field ini
      { name: "no_box", label: "No Box", type: "text", required: true }, // Tambahkan field ini
      { name: "keterangan", label: "Keterangan", type: "text", required: false }, // Tambahkan field ini
      { name: "tanggal_dokumen", label: "Tanggal Dokumen", type: "date", required: true }, // Tambahkan field ini
      { name: "tanggal_penyerahan", label: "Tanggal Penyerahan", type: "date", required: true }, // Tambahkan field ini
    ],
    services: "Arsip",
  });
  const [selectedIds, setSelectedIds] = useState([]);
  const { token } = useToken(); // Ambil token dari context
  let userRole = "";
  if (token) {
    const decoded = jwtDecode(token);
    userRole = decoded.role;
  }

  const [uploadedFiles, setUploadedFiles] = useState([]);

  // UseEffect untuk mengambil data saat komponen dimount
  useEffect(() => {
    getArsip((data) => {
      if (Array.isArray(data)) {
        setMainData(data.reverse());
      } else if (data.posts) {
        setMainData(data.posts.reverse());
      } else {
        console.error("Data tidak dalam format yang diharapkan:", data);
      }
    });

  }, []);

  // Function untuk handle tutup form modal
  const onCloseFormModal = () => {
    setFormModalOpen(false);
    setFormData({});
  };

  // Function untuk fetch data dan update state
  const handleAdd = () => {
    setFormModalOpen(true);
    setFormConfig((prevConfig) => ({
      ...prevConfig,
      action: "add",
      onSubmit: (data) => AddSubmit(data),
    }));
  };

  // Function untuk fetch data dan update state
  const handleEdit = (MainData) => {
    setFormModalOpen(true);
    setFormConfig((prevConfig) => ({
      ...prevConfig,
      action: "edit",
      onSubmit: (data) => EditSubmit(data),
    }));
    setFormData({ ...MainData });
  };

  // tambah data
  const AddSubmit = async (data) => {
    console.log("Data yang dikirim ke server:", data);
    try {
      await addArsip(data); // tambah data ke API
      Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data berhasil ditambahkan",
        showConfirmButton: false,
        timer: 1500,
      }).then(() => {
        window.location.reload();
        setMainData([...MainData, data]);
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

  // ubah data
  const EditSubmit = async (data) => {
    try {
      await updateArsip(data.ID, data); // edit data ke API
      Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data berhasil diperbarui",
        showConfirmButton: false,
        timer: 1500,
      }).then(() => {
        setMainData(
          MainData.map((item) => {
            return item.ID === data.ID ? data : item;
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

  // Function untuk hapus 1 data
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
          await deleteArsip(id); // hapus data di API
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
          await Promise.all(selectedIds.map((id) => deleteArsip(id))); // hapus data di API
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

  const [uploadModalOpen, setUploadModalOpen] = useState(false);
  const [file, setFile] = useState(null);

  const handleUpload = (id) => {
    setSelectedId(id); // Set selectedId untuk file yang diupload
    setUploadModalOpen(true); // Buka modal upload
  };

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleFileUpload = async () => {
    const formData = new FormData();
    formData.append("file", file);

    try {
      await uploadFile(formData); // Panggil fungsi uploadFile yang akan dibuat
      Swal.fire("Berhasil!", "File berhasil diupload", "success");
    } catch (error) {
      Swal.fire("Gagal!", "Error saat mengupload file", "error");
    } finally {
      setUploadModalOpen(false);
      setFile(null);
    }
  };

  // Tambahkan fungsi uploadFile
  const uploadFile = async (data) => {
    await fetch("http://localhost:8080/upload", {
      method: "POST",
      body: data,
    });
  };

  // Tambahkan fungsi handleDeleteFile
  const handleDeleteFile = async (filename) => {
    Swal.fire({
      title: "Apakah Anda yakin?",
      text: "Anda akan menghapus file ini!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          const response = await deleteFile(filename); // Panggil fungsi deleteFile
          if (response.ok) { // Periksa status respon
            setUploadedFiles((prevFiles) => prevFiles.filter((file) => file !== filename)); // Update state
            Swal.fire("Berhasil!", "File berhasil dihapus", "success");
          } else {
            throw new Error("Gagal menghapus file"); // Tangani error jika respon tidak ok
          }
        } catch (error) {
          Swal.fire("Gagal!", "Error saat menghapus file", "error");
        }
      }
    });
  };

  return (
    <App services={formConfig.services}>
      <div className="overflow-auto">
        {/* Table */}
        <ReusableTable
          MainData={MainData}
          formConfig={formConfig}
          handleAdd={handleAdd}
          handleEdit={handleEdit}
          handleDelete={handleDelete}
          handleSelect={handleSelect}
          selectedIds={selectedIds}
          handleBulkDelete={handleBulkDelete}
          uploadedFiles={uploadedFiles}
          excel={
            <Excel
              linkExportThis="exportArsip"
              linkUpdateThis="updateArsip"
              importExcel="uploadArsip"
            />
          }
          handleDeleteFile={handleDeleteFile} // Tambahkan props ini
          handleUpload={handleUpload} // Tambahkan props ini
          showUploadButton={true} // Tambahkan ini
        />
        {/* Modal untuk upload file */}
        <Modal show={uploadModalOpen} size="md" onClose={() => setUploadModalOpen(false)}>
          <Modal.Header>Upload File</Modal.Header>
          <Modal.Body>
            <input type="file" onChange={handleFileChange} />
          </Modal.Body>
          <Modal.Footer>
            <Button onClick={handleFileUpload}>Upload</Button>
            <Button onClick={() => setUploadModalOpen(false)}>Batal</Button>
          </Modal.Footer>
        </Modal>
        {/* End Modal */}
      </div>
    </App>
  );
}