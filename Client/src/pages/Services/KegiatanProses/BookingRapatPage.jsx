import { useState, useEffect } from "react";
import { v4 as uuidv4 } from "uuid"; // Import UUID
import Swal from "sweetalert2";
import App from "../../../components/Layouts/App";
import { BookingRapatCalendar } from "../../../components/Fragments/Services/Calendar/BookingRapatCalendar";
import {
  getBookingRapat,
  addBookingRapat,
  deleteBookingRapat,
} from "../../../../API/KegiatanProses/BookingRapat.service";

export function BookingRapatPage() {
  const [currentEvents, setCurrentEvents] = useState([]);
  const [eventColor, setEventColor] = useState(""); // Menambahkan state untuk warna
  const handleColorChange = (color) => {
    setEventColor(color);
  };
  // Fetch events
  useEffect(() => {
    getBookingRapat((data) => {
      setCurrentEvents(data.reverse());
    });
  }, []);

  // Handle date click to add new event
  const handleDateClick = async (selected) => {
    const { value: title } = await Swal.fire({
      title: "Masukan Event!",
      input: "text",
      inputAttributes: {
        autocapitalize: "off",
      },
      showCancelButton: true,
      confirmButtonText: "Simpan",
      showLoaderOnConfirm: true,
      preConfirm: (e) => {
        return {
          id: uuidv4(),
          title: e,
          start: selected.startStr,
          end: selected.endStr,
          allDay: selected.allDay,
          color: eventColor, // Pastikan ini menggunakan state eventColor yang terbaru
        };
      },
    });
    if (title) {
      try {
        const newEvent = await addBookingRapat(title);
        setCurrentEvents((prevEvents) => [...prevEvents, newEvent]); // Pastikan event baru ditambahkan
        window.location.reload(); // Ini mungkin perlu dihapus untuk menghindari reload
      } catch (error) {
        Swal.fire({
          icon: "error",
          title: "Gagal!",
          text: "Error saat menyimpan data: " + error.message,
          showConfirmButton: false,
          timer: 1500,
        });
      }
    } else {
      calendarApi.unselect();
    }
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
          await deleteBookingRapat(selected.event.id);
          setCurrentEvents((prevEvents) =>
            prevEvents.filter((event) => event.id !== selected.event.id)
          );
          window.location.reload();
        } catch (error) {
          // Perbaikan di sini
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
    <App services="Booking Rapat">
      <BookingRapatCalendar
        currentEvents={currentEvents}
        handleDateClick={handleDateClick}
        handleEventClick={handleEventClick}
        onColorChange={handleColorChange} // Pastikan ini sesuai
        eventColor={eventColor}
      />
    </App>
  );
}
