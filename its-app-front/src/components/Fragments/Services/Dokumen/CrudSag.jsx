import { Button, Table, Badge, Label, Modal, TextInput } from "flowbite-react";
import { useEffect, useState } from "react";
import { getSags } from "../../../../../services/sag.service";

export function IndexSag(props) {
  const { services } = props;
  const [sags, setSagsState] = useState([]);
  const [openModal, setOpenModal] = useState(false);
  const [action, setAction] = useState("");
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
  const [Tanggal, setTanggal] = useState("");
  const [NoMemo, setNoMemo] = useState("");
  const [Perihal, setPerihal] = useState("");
  const [Pic, setPic] = useState("");

  const handleAdd = () => {
    setOpenModal(true);
    setAction("add");
    setTanggal(""); // reset tanggal field
    setNoMemo(""); // reset noMemo field
    setPerihal(""); // reset perihal field
    setPic(""); // reset pic field
  };

  const handleEdit = (sag) => {
    setOpenModal(true);
    setAction("edit");
    setTanggal(sag.Tanggal); // set initial value of tanggal field
    setNoMemo(sag.NoMemo); // set initial value of noMemo field
    setPerihal(sag.Perihal); // set initial value of perihal field
    setPic(sag.Pic); // set initial value of pic field
  };

  const handleSubmit = () => {
    if (action === "add") {
      const newData = {
        tanggal: Tanggal,
        no_memo: NoMemo,
        perihal: Perihal,
        pic: Pic,
      };
      fetch("http://localhost:8080/sag", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newData),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          alert("Data added successfully:", data);
          setSagsState([...sags, data]);
          onCloseModal();
        })
        .catch((error) => {
          alert("Error adding data:", error);
        });
    } else if (action === "edit") {
      const updatedData = {
        id: sag.ID,
        tanggal: Tanggal,
        no_memo: NoMemo,
        perihal: Perihal,
        pic: Pic,
      };
      fetch(`http://localhost:8080/sag/${sag.ID}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedData),
      })
        .then((response) => response.json())
        .then((data) => {
          alert("Data updated successfully:", data);
          setSagsState(sags.map((sag) => (sag.ID === data.ID ? data : sag)));
          onCloseModal();
        })
        .catch((error) => {
          alert("Error updating data:", error);
        });
    }
  };

  const onCloseModal = () => {
    setOpenModal(false);
    setTanggal(""); // reset tanggal field
    setNoMemo(""); // reset noMemo field
    setPerihal(""); // reset perihal field
    setPic(""); // reset pic field
  };

  useEffect(() => {
    getSags((data) => {
      setSagsState(data);
    });
  }, []);

  return (
    <div className="overflow-x-auto p-2">
      <Button onClick={handleAdd} action="add" className="ml-2 mb-4">
        Tambah
      </Button>
      <Table hoverable>
        <Table.Head>
          <Table.HeadCell>No.</Table.HeadCell>
          <Table.HeadCell>Tanggal</Table.HeadCell>
          <Table.HeadCell>No Memo</Table.HeadCell>
          <Table.HeadCell>Perihal</Table.HeadCell>
          <Table.HeadCell>Pic</Table.HeadCell>
          <Table.HeadCell>
            <span className="sr-only">Edit</span>
          </Table.HeadCell>
        </Table.Head>
        {sags.length > 0 ? (
          <Table.Body className="divide-y">
            {sags.map((sag, index) => (
              <Table.Row
                key={sag.ID}
                className="bg-white dark:border-gray-700 dark:bg-gray-800"
              >
                <Table.Cell>
                  <Badge className="flex justify-center">{index + 1}</Badge>
                </Table.Cell>
                <Table.Cell> {formatDate(sag.Tanggal)}</Table.Cell>
                <Table.Cell>{sag.NoMemo}</Table.Cell>
                <Table.Cell>{sag.Perihal}</Table.Cell>
                <Table.Cell>{sag.Pic}</Table.Cell>
                <Table.Cell>
                  <div className="flex gap-2">
                    <a href="#" className="font-medium">
                      <Button onClick={handleEdit} action="add" color="yellow">
                        Ubah
                      </Button>
                    </a>
                    <a href="#" className="font-medium">
                      <Button color="red">Hapus</Button>
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
      {/* Modal */}
      <Modal show={openModal} size="md" onClose={onCloseModal} popup>
        <Modal.Header />
        <Modal.Body>
          <form onSubmit={handleSubmit}>
            <div className="space-y-6">
              <h3 className="flex gap-1 justify-center text-xl font-medium text-gray-900 dark:text-white">
                {action === "add" ? "Tambah Data" : "Ubah Data"}{" "}
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
              <div className="w-full flex gap-4">
                <Button color={"green"} type="submit">
                  {action === "add" ? "Tambah" : "Ubah"}
                </Button>
                <Button color={"yellow"} type="reset">
                  reset
                </Button>
              </div>
            </div>
          </form>
        </Modal.Body>
      </Modal>
      {/* endModal */}
    </div>
  );
}
