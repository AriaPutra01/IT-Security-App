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
import { getSags, deleteSag } from "../../../../services/sag.service";
import { SagForm } from "../../../components/Fragments/Services/Dokumen/Sag/SagForm";
import { HiOutlineExclamationCircle } from "react-icons/hi";
import { SearchInput } from "../../../components/Elements/SearchInput";

// mendefinisikan komponen utama Sag
export function SagPage() {
  const [sags, setSagsState] = useState([]);
  const [formModalOpen, setFormModalOpen] = useState(false); // State for form modal
  const [action, setAction] = useState("add");
  const [Tanggal, setTanggal] = useState("");
  const [NoMemo, setNoMemo] = useState("");
  const [Perihal, setPerihal] = useState("");
  const [Pic, setPic] = useState("");
  const [sag, setSag] = useState(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);

  // Function untuk fetch data sag dan update state
  const handleAdd = () => {
    setFormModalOpen(true);
    setAction("add");
    setTanggal("");
    setNoMemo("");
    setPerihal("");
    setPic("");
  };

  // Function untuk fetch data sag dan update state
  const handleEdit = (sag) => {
    setFormModalOpen(true);
    setAction("edit");
    const tanggal = new Date(sag.Tanggal);
    setTanggal(tanggal.toISOString().split("T")[0]);
    setNoMemo(sag.NoMemo);
    setPerihal(sag.Perihal);
    setPic(sag.Pic);
    setSag(sag);
  };

  // Function untuk meng handle inputan baik tambah maupun ubah
  const handleSubmit = async (event) => {
    event.preventDefault(); // Hentikan event default agar halaman tidak reload otomatis

    if (action === "add") {
      const newData = {
        tanggal: Tanggal.trim(),
        no_memo: NoMemo.trim(),
        perihal: Perihal.trim(),
        pic: Pic.trim(),
      };

      if (
        !newData.tanggal ||
        !newData.no_memo ||
        !newData.perihal ||
        !newData.pic
      ) {
        Swal.fire({
          icon: "info",
          title: "Gagal!",
          text: "Mohon untuk mengisi semua kolom",
        });
        return;
      }

      try {
        const response = await fetch("http://localhost:8080/sag", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(newData),
        });

        if (!response.ok) {
          throw new Error("Network response was not ok");
        }

        const data = await response.json();

        // Panggil SweetAlert setelah berhasil menambahkan data
        Swal.fire({
          icon: "success",
          title: "Berhasil!",
          text: "Data berhasil ditambahkan",
          showConfirmButton: false,
          timer: 1500,
        });

        setSagsState([...sags, data]);
        onCloseFormModal();
      } catch (error) {
        Swal.fire({
          icon: "error",
          title: "Gagal!",
          text: "error terjadi saat menyimpan data",
          showConfirmButton: false,
          timer: 1500,
        });
      }
    } else if (action === "edit") {
      const updatedData = {
        id: sag.ID,
        tanggal: Tanggal.trim(),
        no_memo: NoMemo.trim(),
        perihal: Perihal.trim(),
        pic: Pic.trim(),
      };

      try {
        const response = await fetch(`http://localhost:8080/sag/${sag.ID}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(updatedData),
        });

        if (!response.ok) {
          throw new Error("Network response was not ok");
        }

        const data = await response.json();

        Swal.fire({
          icon: "success",
          title: "Berhasil!",
          text: "Data berhasil diperbarui",
          showConfirmButton: false,
          timer: 1500,
        });

        setSagsState(sags.map((sag) => (sag.ID === data.ID ? data : sag)));
        onCloseFormModal();
      } catch (error) {
        console.error("Error updating data:", error);
        Swal.fire({
          icon: "error",
          title: "Gagal!",
          text: "Gagal mengubah data",
          showConfirmButton: false,
          timer: 1500,
        });
      }
    }
  };

  // Function untuk handle tutup form modal
  const onCloseFormModal = () => {
    setFormModalOpen(false);
    setTanggal("");
    setNoMemo("");
    setPerihal("");
    setPic("");
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
          setSagsState(sags.filter((sag) => sag.ID !== id));
          Swal.fire("Berhasil!", "Data berhasil dihapus", "success");
        } catch (error) {
          Swal.fire("Gagal!", "Error saat hapus data:", "error");
        }
      }
    });
  };

  // UseEffect untuk mengambil data saat komponen dimount dan di balikan urutan
  useEffect(() => {
    getSags((data) => {
      setSagsState(data.reverse());
    });
  }, []);

  // Function untuk menangani pagination
  const paginate = (pageNumber) => {
    setCurrentPage(pageNumber);
  };

  // Hitung indeks awal dan akhir untuk penomoran paginate
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;

  // Function untuk handle search
  const handleSearchChange = (event) => {
    setSearchTerm(event.target.value || "");
  };

  // Get data paginate dan filter search
  const paginatedSags = sags
    .filter((sag) => {
      // Menggunakan optional chaining dan default value
      const tanggal = FormatDate(sag.Tanggal)?.toLowerCase() || "";
      const noMemo = sag.NoMemo?.toLowerCase() || "";
      const perihal = sag.Perihal?.toLowerCase() || "";
      const pic = sag.Pic?.toLowerCase() || "";
      const search = searchTerm.toLowerCase();

      return (
        tanggal.includes(search) ||
        noMemo.includes(search) ||
        perihal.includes(search) ||
        pic.includes(search)
      );
    })
    .slice(startIndex, endIndex);

  return (
    <App services="sag">
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
              <Dropdown.Item>This Sheet</Dropdown.Item>
              <Dropdown.Item>All Sheet</Dropdown.Item>
            </Dropdown>
            <Dropdown
              color="success"
              label="Update Excel"
              dismissOnClick={false}
            >
              <Dropdown.Item>This Sheet</Dropdown.Item>
              <Dropdown.Item>All Sheet</Dropdown.Item>
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
                    <Table.Cell>{FormatDate(sag.Tanggal)}</Table.Cell>
                    <Table.Cell>{sag.NoMemo}</Table.Cell>
                    <Table.Cell>{sag.Perihal}</Table.Cell>
                    <Table.Cell>{sag.Pic}</Table.Cell>
                    <Table.Cell>
                      <div className="flex gap-2">
                        <a href="#" className="font-medium">
                          <Button
                            onClick={() => handleEdit(sag)}
                            action="edit"
                            color="warning"
                          >
                            Ubah
                          </Button>
                        </a>
                        <a href="#" className="font-medium">
                          <Button
                            onClick={() => {
                              setSag(sag);
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
            <SagForm
              onSubmit={handleSubmit}
              action={action}
              services="sag"
              tanggal={Tanggal}
              noMemo={NoMemo}
              perihal={Perihal}
              pic={Pic}
              setTanggal={setTanggal}
              setNoMemo={setNoMemo}
              setPerihal={setPerihal}
              setPic={setPic}
            />
          </Modal.Body>
        </Modal>
        {/* endModalForm */}
      </div>
    </App>
  );
}
