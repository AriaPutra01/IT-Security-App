import { Button, Table, Badge, Label, Modal, TextInput } from "flowbite-react";
import { useState } from "react";
export function IndexDokumen(props) {
  const { services } = props;
  const [openModal, setOpenModal] = useState(false);
  const [action, setAction] = useState("");
  const [tanggal, setTanggal] = useState("");
  const kolom = services === "sag" || services === "iso" || services === "memo";
  const handleAdd = () => {
    setOpenModal(true);
    setAction("add");
  };

  const handleEdit = () => {
    setOpenModal(true);
    setAction("edit");
  };
  const onCloseModal = () => {
    setOpenModal(false);
    setEmail("");
  };

  return (
    <div className="overflow-x-auto p-2">
      <Button onClick={handleAdd} className="ml-2 mb-4">
        Tambah
      </Button>
      <Table hoverable>
        <Table.Head>
          <Table.HeadCell>No.</Table.HeadCell>
          <Table.HeadCell>Tanggal</Table.HeadCell>
          <Table.HeadCell>{kolom ? "No Memo" : "No Surat"}</Table.HeadCell>
          <Table.HeadCell>Perihal</Table.HeadCell>
          <Table.HeadCell>Pic</Table.HeadCell>
          <Table.HeadCell>
            <span className="sr-only">Edit</span>
          </Table.HeadCell>
        </Table.Head>
        <Table.Body className="divide-y">
          <Table.Row className="bg-white dark:border-gray-700 dark:bg-gray-800">
            <Table.Cell>
              <Badge className="flex justify-center">1</Badge>
            </Table.Cell>
            <Table.Cell>Apple</Table.Cell>
            <Table.Cell>Sliver</Table.Cell>
            <Table.Cell>Laptop</Table.Cell>
            <Table.Cell>wow</Table.Cell>
            <Table.Cell>
              <div className="flex gap-2">
                <a href="#" className="font-medium">
                  <Button onClick={handleEdit} color="yellow">
                    Ubah
                  </Button>
                </a>
                <a href="#" className="font-medium">
                  <Button color="red">Hapus</Button>
                </a>
              </div>
            </Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>
      {/* Modal */}
      <Modal show={openModal} size="md" onClose={onCloseModal} popup>
        <Modal.Header />
        <Modal.Body>
          <form action="">
            <div className="space-y-6">
              <h3 className="flex gap-1 justify-center text-xl font-medium text-gray-900 dark:text-white">
                {action === "add" ? "Tambah Data" : "Ubah Data"}{" "}
                <div className="uppercase">{services}</div>
              </h3>
              <div>
                <div className="mb-2 block">
                  <Label htmlFor="tanggal" value="Tanggal" />
                </div>
                <TextInput
                  type="date"
                  id="tanggal"
                  value={tanggal}
                  onChange={(event) => setTanggal(event.target.value)}
                  required
                />
              </div>
              <div>
                <div className="mb-2 block">
                  <Label htmlFor="noMemo" value="Nomor Memo" />
                </div>
                <TextInput id="noMemo" type="text" required />
              </div>
              <div>
                <div className="mb-2 block">
                  <Label htmlFor="perihal" value="Perihal" />
                </div>
                <TextInput id="perihal" type="text" required />
              </div>
              <div>
                <div className="mb-2 block">
                  <Label htmlFor="pic" value="PIC" />
                </div>
                <TextInput id="pic" type="text" required />
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
