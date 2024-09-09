import { useState, useEffect } from "react";
import Swal from "sweetalert2";
import App from "../../../components/Layouts/App";
import { Calendar } from "../../../components/Fragments/Services/Calendar/Calendar";
import {
  getRapats,
  addRapat,
  deleteRapat,
} from "../../../../API/KegiatanProses/JadwalRapat.service";
import { Button, Label, Modal, TextInput } from "flowbite-react";
import { ColorPick } from "../../../Utilities/ColorPick";

export function JadwalRapatPage() {
  const [formModalOpen, setFormModalOpen] = useState(false);
  const [formData, setFormData] = useState({});
  const [currentEvents, setCurrentEvents] = useState([]);
  const onCloseFormModal = () => {
    setFormModalOpen(false);
    setFormData({});
  };

  // Fetch events
  useEffect(() => {
    getRapats((data) => {
      setCurrentEvents(data);
    });
  }, []);

  const handleFormChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  };

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    const newEvent = {
      title: formData.title,
      start: formData.start,
      end: formData.end,
      allDay: formData.allDay,
      color: formData.color,
    };
    try {
      await addRapat(newEvent);
      setCurrentEvents((prevEvents) => [...prevEvents, newEvent]);
      onCloseFormModal();
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Gagal!",
        text: "Error saat menyimpan data: " + error.message,
        showConfirmButton: false,
        timer: 1500,
      });
    }
  };

  // Handle date click to add new event
  const handleDateClick = async (selected) => {
    setFormModalOpen(true);
    setFormData({
      title: "",
      start: selected.startStr,
      end: selected.endStr,
      allDay: selected.allDay,
      color: "black",
    });
  };

  // Handle event click to delete event
  const handleEventClick = async (selected) => {
    Swal.fire({
      title: "Apakah Anda yakin?",
      text: `Anda akan menghapus data ${selected.event.title}?`,
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          await deleteRapat(selected.event.id);
          setCurrentEvents((prevEvents) =>
            prevEvents.filter((event) => event.id !== selected.event.id)
          );
        } catch (error) {
          Swal.fire({
            icon: "error",
            title: "Gagal!",
            text: "Error saat menghapus data: " + error.message, // Menampilkan pesan error
            showConfirmButton: false,
            timer: 1500,
          });
        }
      }
    });
  };
  return (
    <App services="Jadwal Rapat">
      <Calendar
        view="timeGridWeek"
        currentEvents={currentEvents}
        handleDateClick={handleDateClick}
        handleEventClick={handleEventClick}
      />
      {/* ModalForm */}
      <Modal show={formModalOpen} size="xl" onClose={onCloseFormModal} popup>
        <Modal.Header />
        <Modal.Body>
          <form onSubmit={handleFormSubmit}>
            <div className="grid grid-cols-4 gap-4">
              <div className="flex flex-col col-span-3">
                <Label htmlFor="title" value="Title" />
                <TextInput
                  id="title"
                  name="title"
                  type="text"
                  className="block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 appearance-none  focus:outline-none focus:ring-0 focus:border-blue-600 peer"
                  placeholder="masukan event"
                  value={formData.title}
                  onChange={handleFormChange}
                  required
                />
              </div>
              <div className="flex flex-col gap-2 justify-start col-span-1">
                <Label htmlFor="color" value="Color" />
                <ColorPick
                  name="color"
                  onChange={handleFormChange}
                  className="w-full h-full mb-2 p-[2px]"
                />
              </div>
              <Button className="col-span-4" type="submit">
                Simpan
              </Button>
            </div>
          </form>
        </Modal.Body>
      </Modal>
      {/* endModalForm */}
    </App>
  );
}
