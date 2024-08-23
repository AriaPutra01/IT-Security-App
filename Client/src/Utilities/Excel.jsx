import axios from "axios";
import Swal from "sweetalert2";
import { useState } from "react";
import { Badge, Button, Dropdown, FileInput } from "flowbite-react";

export const Excel = (props) => {
  const { linkExportThis, linkUpdateThis, importExcel } = props;
  const [file, setFile] = useState(null);
  const [updateThisHref, setUpdateThisHref] = useState("#");
  const [updateAllHref, setUpdateAllHref] = useState("#");
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const UpdateThis = async () => {
    const result = await Swal.fire({
      title: "Update Sheet Ini",
      text: "Anda akan mengupdate sheet ini?",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    });
    if (result.isConfirmed) {
      try {
        setUpdateThisHref(`http://localhost/8080/${linkUpdateThis}`);
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
    } else {
      setUpdateThisHref("#");
    }
  };

  const UpdateAll = async () => {
    const result = await Swal.fire({
      title: "Update Semua Sheet",
      text: "Anda akan mengupdate semua sheet ?",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    });
    if (result.isConfirmed) {
      try {
        setUpdateAllHref("http://localhost:8080/updateAll");
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
    } else {
      setUpdateAllHref("#");
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

        await axios
          .post(importExcel, formData, {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          })
          .then(() => {
            Swal.fire({
              icon: "success",
              title: "Berhasil!",
              text: "Data berhasil dimport",
              showConfirmButton: false,
              timer: 1500,
            });

            setTimeout(() => {
              window.location.reload();
            }, 1000);
          });
      } catch (error) {
        Swal.fire({
          icon: "error",
          title: "gagal!",
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
              <a href={`http://localhost/8080/${linkExportThis}`}>This Sheet</a>
            </Dropdown.Item>
            <Dropdown.Item>
              <a href="http://localhost:8080/exportAll">All Sheets</a>
            </Dropdown.Item>
          </Dropdown>
          <Dropdown color="warning" label="UPDATE" dismissOnClick={false}>
            <Dropdown.Item>
              <a
                href={updateThisHref}
                onClick={() => {
                  UpdateThis();
                }}
              >
                This Sheet
              </a>
            </Dropdown.Item>
            <Dropdown.Item>
              <a
                href={updateAllHref}
                onClick={() => {
                  UpdateAll();
                }}
              >
                All Sheet
              </a>
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
