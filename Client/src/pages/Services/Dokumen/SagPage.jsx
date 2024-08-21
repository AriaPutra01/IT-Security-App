import React, { useState, useEffect } from "react";
import App from "../../../components/Layouts/App";
import Swal from "sweetalert2";
import {
  Button,
  Dropdown,
  Table,
  Badge,
  Modal,
  Pagination,
} from "flowbite-react";
import { FormatDate } from "../../../Utilities/FormatDate";
import {
  getSags,
  addSag,
  updateSag,
  deleteSag,
} from "../../../../services/sag.service";
import { SearchInput } from "../../../components/Elements/SearchInput";
import { ReusableForm } from "../../../components/Fragments/Form/ReusableForm";

// mendefinisikan komponen utama Sag
export function SagPage() {
  const [sags, setSagsState] = useState([]);
  const [formModalOpen, setFormModalOpen] = useState(false); // State for form modal
  const [formData, setFormData] = useState({});
  const [formConfig, setFormConfig] = useState({
    fields: [
      { name: "tanggal", label: "Tanggal", type: "date", required: true },
      { name: "no_memo", label: "Nomor Memo", type: "text", required: true },
      { name: "perihal", label: "Perihal", type: "text", required: true },
      { name: "pic", label: "Pic", type: "text", required: true },
    ],
    services: "sag",
  });
  const [searchTerm, setSearchTerm] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [updateModalOpen, setUpdateModalOpen] = useState(false); // State for update modal

  // UseEffect untuk mengambil data saat komponen dimount dan di balikan urutan
  useEffect(() => {
    getSags((data) => {
      setSagsState(data.reverse());
    });
  }, []);

  // Function untuk handle tutup form modal
  const onCloseFormModal = () => {
    setFormModalOpen(false);
    setFormData({});
  };

  // Function untuk handle tutup update modal
  const onCloseUpdateModal = () => {
    setUpdateModalOpen(false);
  };

  // Function untuk fetch data sag dan update state
  const handleAdd = () => {
    setFormModalOpen(true);
    setFormConfig((prevConfig) => ({
      ...prevConfig,
      action: "add",
      onSubmit: (data) => AddSubmit(data),
    }));
  };

  // Function untuk fetch data sag dan update state
  const handleEdit = (sag) => {
    setFormModalOpen(true);
    setFormConfig((prevConfig) => ({
      ...prevConfig,
      action: "edit",
      onSubmit: (data) => EditSubmit(data),
    }));
    setFormData({ ...sag });
  };

  // tambah data
  const AddSubmit = async (data) => {
    try {
      await addSag(data);
      Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data berhasil ditambahkan",
        showConfirmButton: false,
        timer: 1000,
      }).then(() => {
        setSagsState([...sags, data]);
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
    const sag = sags.find((s) => s.ID === data.ID);

    if (!data.tanggal) {
      data.tanggal = sag.tanggal;
    }

    // Log data untuk memeriksa sebelum dikirim
    console.log('Data yang dikirim:', data);

    try {
      const response = await updateSag(data.ID, data);

      Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data berhasil diperbarui",
        showConfirmButton: false,
        timer: 1500,
      }).then(() => {
        setSagsState(
          sags.map((sag) => {
            return sag.ID === response.ID ? data : sag;
          })
        );
      });
    } catch (error) {
      console.error('Error saat mengubah data:', error);

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
          await deleteSag(id);
          setSagsState((prevSags) => prevSags.filter((sag) => sag.ID !== id));
          Swal.fire({
            icon: "info",
            title: "Berhasil!",
            text: "Data berhasil dihapus",
            showConfirmButton: false,
          });
          setTimeout(() => {
          }, 1000);
        } catch (error) {
          Swal.fire("Gagal!", "Error saat hapus data:", "error");
        }
      }
    });
  };

  // Function untuk handle search
  const handleSearchChange = (event) => {
    setSearchTerm(event.target.value || "");
  };

  // Hitung indeks awal dan akhir untuk penomoran paginate
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;

  // Get data paginate dan filter search
  const paginatedSags = sags
    .filter((sag) => {
      // Menggunakan optional chaining dan default value
      const tanggal = FormatDate(sag.tanggal)?.toLowerCase() || "";
      const noMemo = sag.no_memo?.toLowerCase() || "";
      const perihal = sag.perihal?.toLowerCase() || "";
      const pic = sag.pic?.toLowerCase() || "";
      const search = searchTerm.toLowerCase();
      return (
        tanggal.includes(search) ||
        noMemo.includes(search) ||
        perihal.includes(search) ||
        pic.includes(search)
      );
    })
    .slice(startIndex, endIndex);

  // Function untuk menangani pagination
  const paginate = (pageNumber) => {
    setCurrentPage(pageNumber);
  };

  const handleExport = () => {
    window.location.href = "http://localhost:8080/exportSag";
  };

  const handleExportAll = () => {
    window.location.href = "http://localhost:8080/exportAll";
  };

  const handleUpdate = () => {
    setUpdateModalOpen(true);
  };

  const handleUpdateSubmit = () => {
    // Logic untuk submit update
    console.log("Update submitted");
    window.location.href = "http://localhost:8080/updateSag"; // Ganti dengan endpoint yang sesuai
    onCloseUpdateModal();
  };

  const handleUpdateAll = () => {
    window.location.href = "http://localhost:8080/updateAll";
  };

  return (
    <App services={formConfig.services}>
      <div className="w-full h-full">
        {/* page title */}
        <div className="flex justify-between items-center mx-2 mb-2 ">
          <div className="flex gap-1.5 items-center justify-center">
            <Button
              className="flex justify-center items-center"
              onClick={handleAdd}
              action="add"
              color="info"
            >
              Tambah
            </Button>
            <Dropdown
              color="success"
              label="Export Excel"
              dismissOnClick={false}
            >
              <Dropdown.Item onClick={handleExport}><a href="http://localhost:8080/exportSag">Export This Sheet</a></Dropdown.Item>
              <Dropdown.Item onClick={handleExportAll}><a href="http://localhost:8080/exportAll"> All Sheet</a></Dropdown.Item>
            </Dropdown>
            <Dropdown
              color="success"
              label="Update Excel"
              dismissOnClick={false}
            >
              <Dropdown.Item onClick={handleUpdate}>This Sheet</Dropdown.Item>
              <Dropdown.Item onClick={handleUpdateAll}><a href="http://localhost:8080/updateAll">All Sheet</a></Dropdown.Item>
            </Dropdown>
            <Button color="warning">Import Excel</Button>
          </div>
          <SearchInput value={searchTerm} onChange={handleSearchChange} />
        </div>
        {/* End page title */}

        {/* table */}
        <div className="h-2/3 overflow-y-auto p-2">
          <Table hoverable>
            <Table.Head>
              <Table.HeadCell>No.</Table.HeadCell>
              <Table.HeadCell>Tanggal</Table.HeadCell>
              <Table.HeadCell>No Memo</Table.HeadCell>
              <Table.HeadCell>Perihal</Table.HeadCell>
              <Table.HeadCell>Pic</Table.HeadCell>
              <Table.HeadCell>
                <span>Action</span>
              </Table.HeadCell>
            </Table.Head>
            {paginatedSags.length > 0 ? (
              <Table.Body className="divide-y">
                {paginatedSags.map((sag, index) => (
                  <Table.Row key={sag.ID}>
                    <Table.Cell>
                      <Badge className="flex justify-center">
                        {(currentPage - 1) * itemsPerPage + index + 1}
                      </Badge>
                    </Table.Cell>
                    <Table.Cell>{FormatDate(sag.tanggal)}</Table.Cell>
                    <Table.Cell>{sag.no_memo}</Table.Cell>
                    <Table.Cell>{sag.perihal}</Table.Cell>
                    <Table.Cell>{sag.pic}</Table.Cell>
                    <Table.Cell>
                      <div className="flex gap-2">
                        <a href="#" className="font-medium">
                          <Button
                            onClick={() => {
                              handleEdit(sag);
                            }}
                            action="edit"
                            color="warning"
                          >
                            Ubah
                          </Button>
                        </a>
                        <a href="#" className="font-medium">
                          <Button
                            onClick={() => {
                              setSagsState(sag);
                              handleDelete(sag.ID);
                            }}
                            color="failure"
                          >
                            Hapus
                          </Button>
                        </a>
                      </div>
                    </Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            ) : (
              <Table.Body className="divide-y">
                <Table.Row>
                  <Table.Cell colSpan={6} className="text-center">
                    <Badge className="p-4 font-bold" color="red">
                      Belum Ada Data yang Tersedia
                    </Badge>
                  </Table.Cell>
                </Table.Row>
              </Table.Body>
            )}
          </Table>
        </div>
        {/* End Table */}

        {/* Pagination */}
        <div className="flex overflow-x-auto m-2">
          <Pagination
            currentPage={currentPage}
            totalPages={Math.ceil(sags.length / itemsPerPage)}
            onPageChange={(newPage) => paginate(newPage)}
          />
        </div>
        {/* End Pagination */}

        {/* ModalForm */}
        <Modal show={formModalOpen} size="md" onClose={onCloseFormModal} popup>
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

        {/* Update Modal */}
        <Modal show={updateModalOpen} size="md" onClose={onCloseUpdateModal} popup>
          <Modal.Header />
          <Modal.Body>
            <div className="flex flex-col gap-4">
              <h3 className="text-lg font-medium text-gray-900">Update This Sheet</h3>
              <Button color="success" onClick={handleUpdateSubmit}>Submit</Button>
            </div>
          </Modal.Body>
        </Modal>
        {/* endUpdateModal */}
      </div>
    </App>
  );
}