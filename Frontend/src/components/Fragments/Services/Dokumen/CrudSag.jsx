import React, { useState, useEffect } from "react";
import {
  Button,
  Table,
  Badge,
  Label,
  Modal,
  TextInput,
  Pagination,
} from "flowbite-react";
import { HiOutlineExclamationCircle } from "react-icons/hi";
import { getSags, deleteSag } from "../../../../../services/sag.service";

// Format tanggal ke string dengan format indonesia
const formatDate = (dateString) => {
  const date = new Date(dateString);
  const months = [
    "Januari",
    "Februari",
    "Maret",
    "April",
    "Mei",
    "Juni",
    "Juli",
    "Agustus",
    "September",
    "Oktober",
    "November",
    "Desember",
  ];
  return `${date.getDate()} ${months[date.getMonth()]} ${date.getFullYear()}`;
};

// mendefinisikan komponen utama Sag
export function IndexSag(props) {
  const { services } = props;
  const [sags, setSagsState] = useState([]);
  const [formModalOpen, setFormModalOpen] = useState(false); // State for form modal
  const [deleteModalOpen, setDeleteModalOpen] = useState(false); // State for delete modal
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
  const handleSubmit = async () => {
    if (action === "add") {
      const newData = {
        tanggal: Tanggal,
        no_memo: NoMemo,
        perihal: Perihal,
        pic: Pic,
      };

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
        alert("Data added successfully:", data);
        setSagsState([...sags, data]);
        onCloseFormModal();
      } catch (error) {
        alert("Error adding data:", error);
      }
    } else if (action === "edit") {
      const updatedData = {
        id: sag.ID,
        tanggal: Tanggal,
        no_memo: NoMemo,
        perihal: Perihal,
        pic: Pic,
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
        alert("Data updated successfully:", data);
        setSagsState(sags.map((sag) => (sag.ID === data.ID ? data : sag)));
        onCloseFormModal();
      } catch (error) {
        alert("Error updating data:", error);
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

  // Function untuk handle tutup delete modal
  const onCloseDeleteModal = () => {
    setDeleteModalOpen(false);
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

  // Function untuk hapus 1 data
  const handleDelete = async (id) => {
    try {
      await deleteSag(id);
      setSagsState(sags.filter((sag) => sag.ID !== id));
      onCloseDeleteModal();
    } catch (error) {
      alert("Error deleting data:", error);
    }
  };

  // Hitung indeks awal dan akhir untuk penomoran paginate
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;

  // Function untuk handle search
  const handleSearchChange = (event) => {
    setSearchTerm(event.target.value);
  };

  // Get data paginate dan filter search
  const paginatedSags = sags
    .filter(
      (sag) =>
        sag.NoMemo.toLowerCase().includes(searchTerm.toLowerCase()) ||
        sag.Perihal.toLowerCase().includes(searchTerm.toLowerCase()) ||
        sag.Pic.toLowerCase().includes(searchTerm.toLowerCase())
    )
    .slice(startIndex, endIndex);

  return (
    <div className="w-full h-full">
      <div className="flex justify-between mx-4 mb-2">
        <Button
          className="flex justify-center items-center"
          onClick={handleAdd}
          action="add"
          color="success"
        >
          Tambah
        </Button>
        <form className="w-2/5">
          <div className="relative">
            <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
              <svg
                className="w-4 h-4 text-gray-500 dark:text-gray-400"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 20 20"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
                />
              </svg>
            </div>
            <input
              type="search"
              id="default-search"
              className="block w-full p-4 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              placeholder="Search..."
              value={searchTerm}
              onChange={handleSearchChange}
              required
            />
          </div>
        </form>
      </div>
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
                  <Table.Cell>{formatDate(sag.Tanggal)}</Table.Cell>
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
                            setDeleteModalOpen(true);
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
                    Belum Ada Data
                  </Badge>
                </Table.Cell>
              </Table.Row>
            </Table.Body>
          )}
        </Table>
      </div>

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
          <form onSubmit={handleSubmit} className="space-y-6">
            <h3 className="flex gap-1 justify-center text-xl font-medium text-gray-900 dark:text-white">
              {action === "add" ? "Tambah Data" : `Ubah Data`}
              <div className="uppercase">{services}</div>
            </h3>
            <div>
              <div className="mb-2 block">
                <Label htmlFor="Tanggal" value="Tanggal" />
              </div>
              <TextInput
                type="date"
                name="Tanggal"
                id="Tanggal"
                value={Tanggal}
                onChange={(e) => setTanggal(e.target.value)}
                required
              />
            </div>
            <div>
              <div className="mb-2 block">
                <Label htmlFor="NoMemo" value="Nomor Memo" />
              </div>
              <TextInput
                id="NoMemo"
                name="NoMemo"
                type="text"
                value={NoMemo}
                onChange={(e) => setNoMemo(e.target.value)}
                required
              />
            </div>
            <div>
              <div className="mb-2 block">
                <Label htmlFor="Perihal" value="Perihal" />
              </div>
              <TextInput
                id="Perihal"
                name="Perihal"
                type="text"
                value={Perihal}
                onChange={(e) => setPerihal(e.target.value)}
                required
              />
            </div>
            <div>
              <div className="mb-2 block">
                <Label htmlFor="Pic" value="Pic" />
              </div>
              <TextInput
                id="Pic"
                name="Pic"
                type="text"
                value={Pic}
                onChange={(e) => setPic(e.target.value)}
                required
              />
            </div>
            <Button
              className="w-full"
              color={action === "add" ? "success" : "warning"}
              type="submit"
            >
              {action === "add" ? "Tambah" : "Ubah"}
            </Button>
          </form>
        </Modal.Body>
      </Modal>
      {/* endModalForm */}

      {/* ModalDelete */}
      <Modal
        show={deleteModalOpen}
        size="md"
        onClose={onCloseDeleteModal}
        popup
      >
        <Modal.Header />
        <Modal.Body>
          <div className="text-center">
            <HiOutlineExclamationCircle className="mx-auto mb-4 h-14 w-14 text-gray-400 dark:text-gray-200" />
            <h3 className="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
              pakah anda yakin untuk menghapus data ini?
            </h3>
            <div className="flex justify-center gap-4">
              <Button color="failure" onClick={() => handleDelete(sag.ID)}>
                {"Ya, saya yakin"}
              </Button>
              <Button color="gray" onClick={onCloseDeleteModal}>
                Batal
              </Button>
            </div>
          </div>
        </Modal.Body>
      </Modal>
      {/* endModalDelete */}
    </div>
  );
}
