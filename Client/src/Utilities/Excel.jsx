import axios from "axios";
import Swal from "sweetalert2";
import { useState } from "react";
import { Button, Dropdown, FileInput } from "flowbite-react";

export const Excel = (props) => {
  const { linkExportThis, linkUpdateThis, importExcel } = props;
  const [file, setFile] = useState(null);
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const UpdateThis = async () => {
    const result = await Swal.fire({
      icon: "info",
      title: "Update Sheet Ini",
      text: "Anda akan mengupdate sheet ini?",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    });
    if (result.isConfirmed) {
      try {
        await axios.get(`http://localhost:8080/${linkUpdateThis}`);
        Swal.fire({
          icon: "success",
          title: "Berhasil!",
          text: "Data berhasil diupdate ke Excel",
          showConfirmButton: false,
          timer: 1500,
        });
      } catch (error) {
        Swal.fire("Gagal!", "Error saat update data:", "error");
      }
    }
  };

  const UpdateAll = async () => {
    const result = await Swal.fire({
      icon: "info",
      title: "Update Semua Sheet",
      text: "Anda akan mengupdate semua sheet ke excel?",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    });
    if (result.isConfirmed) {
      try {
        await axios.get("http://localhost:8080/updateAll");
        Swal.fire({
          icon: "success",
          title: "Berhasil!",
          text: "Data berhasil diupdate",
          showConfirmButton: false,
          timer: 1500,
        });
      } catch (error) {
        Swal.fire("Gagal!", "Error saat update data:", "error");
      }
    }
  };

  const handleImport = async (event) => {
    event.preventDefault();
    if (!file) {
      alert("Mohon Untuk Menambahkan File.");
      return;
    }
    const fileReader = new FileReader();
    fileReader.readAsBinaryString(file);
    fileReader.onload = async (event) => {
      try {
        const data = event.target.result;
        const fileExtension = file.name.split(".").pop().toLowerCase();
        if (fileExtension !== "xlsx") {
          throw new Error("File format harus berupa .xlsx");
        }
        const formData = new FormData();
        formData.append("file", file);
        console.log("URL untuk import:", `http://localhost:8080/${importExcel}`); // Tambahkan log ini
        await axios.post(`http://localhost:8080/${importExcel}`, formData, { // Pastikan importExcel digunakan di sini
          headers: {
            "Content-Type": "multipart/form-data",
          },
        });
        Swal.fire({
          icon: "success",
          title: "Berhasil!",
          text: "Data berhasil diimport",
          showConfirmButton: false,
          timer: 1500,
        });
        setTimeout(() => {
          window.location.reload();
        }, 1500);
      } catch (error) {
        console.error("Error saat mengimport data:", error); // Tambahkan log error ini
        Swal.fire({
          icon: "error",
          title: "Gagal!",
          text: "Mohon untuk memasukkan file.xlsx",
          showConfirmButton: false,
          timer: 1500,
        });
      }
    };
  };

  return (
    <div className="flex gap-1.5 items-center justify-center">
      <Dropdown color="success" label="Excel" dismissOnClick={false}>
        <Dropdown.Item className="flex justify-between">
          <Dropdown color="info" label="EXPORT" dismissOnClick={false}>
            <Dropdown.Item>
              <a href={`http://localhost:8080/${linkExportThis}`}>This Sheet</a>
            </Dropdown.Item>
            <Dropdown.Item>
              <a href="http://localhost:8080/exportAll">All Sheets</a>
            </Dropdown.Item>
          </Dropdown>
          <Dropdown color="warning" label="UPDATE" dismissOnClick={false}>
            <Dropdown.Item>
              <a onClick={UpdateThis}>This Sheet</a>
            </Dropdown.Item>
            <Dropdown.Item>
              <a onClick={UpdateAll}>All Sheet</a>
            </Dropdown.Item>
          </Dropdown>
        </Dropdown.Item>
        <Dropdown.Item className="flex flex-col gap-2">
          <FileInput onChange={handleFileChange} />
          <Button onClick={handleImport} color="success" className="w-full">
            Import
          </Button>
        </Dropdown.Item>
      </Dropdown>
    </div>
  );
};