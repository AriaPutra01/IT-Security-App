import { useState, useEffect } from "react";
import Swal from "sweetalert2";
import { Button, Label, Modal, TextInput } from "flowbite-react";
import { ColorPick } from "../../../../Utilities/ColorPick";
import FullCalendar from "@fullcalendar/react";
import { formatDate } from "@fullcalendar/core";
import dayGridPlugin from "@fullcalendar/daygrid";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
import listPlugin from "@fullcalendar/list";

export const Calendar = ({ view, get, add, remove }) => {
  const [formModalOpen, setFormModalOpen] = useState(false);
  const [formData, setFormData] = useState({});
  const [currentEvents, setCurrentEvents] = useState([]);
  const onCloseFormModal = () => {
    setFormModalOpen(false);
    setFormData({});
  };

  // Fetch events
  useEffect(() => {
    get((data) => {
      setCurrentEvents(data.reverse());
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
      const response = await add(newEvent);
      setCurrentEvents([response, ...currentEvents]);
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
  const handleDateClick = (selected) => {
    setFormModalOpen(true);
    setFormData({
      title: "",
      start: selected.startStr,
      end: selected.endStr,
      allDay: selected.allDay,
      color: "#2596be",
    });
  };

  // Handle event click to delete event
  const handleEventClick = (selected) => {
    Swal.fire({
      title: "Apakah Anda yakin?",
      text: `Anda akan menghapus data ${selected.event.title}?`,
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    }).then((result) => {
      if (result.isConfirmed) {
        remove(selected.event.id)
          .then(() => {
            get((data) => {
              setCurrentEvents(data); // Pastikan data adalah array
            });
          })
          .catch((error) => {
            Swal.fire({
              icon: "error",
              title: "Gagal!",
              text: "Error saat hapus data: " + error.message,
              showConfirmButton: false,
              timer: 1500,
            });
          });
      }
    });
  };

  return (
    <div className="grid grid-cols-2fr">
      <div className="bg-gray-50 p-[15px] rounded w-[200px] max-h-[85vh] overflow-auto">
        <h2 className="text-xl mt-0 mb-2 font-bold">
          {currentEvents.length} Jadwal
        </h2>
        <div className="flex flex-col gap-2">
          {currentEvents.length === 0 ? (
            <div className="ring-2 ring-blue-200 rounded p-2">
              <p>
                Belum ada jadwal yang tersedia, tambahkan jadwal baru di menu
                Calendar.
              </p>
            </div>
          ) : (
            currentEvents.map((event) => (
              <div
                key={event.id}
                className="overflow-auto grow border-b-2 border-l-2 border-sky-500 shadow p-2 rounded"
              >
                <div className="text-sky-500 font-bold">{event.title}</div>
                <div className="">
                  {formatDate(event.start, {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                  })}
                </div>
              </div>
            ))
          )}
        </div>
      </div>
      <div className="mx-3 max-h-[85vh] overflow-auto">
        <FullCalendar
          plugins={[
            dayGridPlugin,
            timeGridPlugin,
            interactionPlugin,
            listPlugin,
          ]}
          headerToolbar={{
            left: "prev,next today",
            center: "title",
            right: "dayGridMonth,timeGridWeek,timeGridDay,listMonth",
          }}
          initialView={view}
          editable={true}
          selectable={true}
          selectMirror={true}
          dayMaxEvents={true}
          events={currentEvents}
          select={handleDateClick}
          eventClick={handleEventClick}
        />
      </div>
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
                  value={formData.color}
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
    </div>
  );
};
