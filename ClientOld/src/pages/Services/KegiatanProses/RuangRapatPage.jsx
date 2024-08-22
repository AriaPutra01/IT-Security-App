import { useState, useEffect } from "react";
import FullCalendar from "@fullcalendar/react";
import { formatDate } from "@fullcalendar/core";
import dayGridPlugin from "@fullcalendar/daygrid";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
import listPlugin from "@fullcalendar/list";
import { v4 as uuidv4 } from "uuid"; // Import UUID
import "../../../calender.css";
import Swal from "sweetalert2";
import {
  getEvents,
  addEvent,
  deleteEvent,
} from "../../../../API/KegiatanProses/rapat.service";
import App from "../../../components/Layouts/App";

export function RuangRapatPage() {
  const [currentEvents, setCurrentEvents] = useState([]);
  console.log("🚀 ~ RuangRapatPage ~ currentEvents:", currentEvents);

  // Fetch events
  useEffect(() => {
    getEvents((data) => {
      setCurrentEvents(data.reverse());
    });
  }, []);

  // Handle date click to add new event
  const handleDateClick = async (selected) => {
    const title = prompt("Please enter a new title for your event");
    const calendarApi = selected.view.calendar;
    if (title) {
      const newEvent = {
        id: uuidv4(),
        title,
        start: selected.startStr,
        end: selected.endStr,
        allDay: selected.allDay,
      };
      try {
        await addEvent(newEvent);
        setCurrentEvents([...currentEvents, newEvent]);
      } catch {
        (error) => {
          Swal.fire({
            icon: error.message,
            title: "Gagal!",
            text: "Error saat menyimpan data",
            showConfirmButton: false,
            timer: 1500,
          });
        };
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
          await deleteEvent(selected.event.id);
          setCurrentEvents((prevEvents) =>
            prevEvents.filter((event) => event.id !== selected.event.id)
          );
          Swal.fire({
            icon: "info",
            title: "Berhasil!",
            text: "Data berhasil dihapus",
            showConfirmButton: false,
            timer: 1500,
          });
        } catch (error) { // Perbaikan di sini
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
    <App services="Ruang Rapat">
      <div className="calendar-container">
        <div className="calendar-content">
          <div className="events-list">
            <h2>Events</h2>
            <ul>
              {currentEvents.map((event, index) => (
                <li key={index} className="event-item">
                  <div className="event-title">{event.title}</div>
                  <div className="event-date">
                    {formatDate(event.start, {
                      year: "numeric",
                      month: "short",
                      day: "numeric",
                    })}
                  </div>
                </li>
              ))}
            </ul>
          </div>

          <div className="calendar">
            <FullCalendar
              height="75vh"
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
              initialView="dayGridMonth"
              editable={true}
              selectable={true}
              selectMirror={true}
              dayMaxEvents={true}
              select={handleDateClick}
              eventClick={handleEventClick}
              events={currentEvents}
            />
          </div>
        </div>
      </div>
    </App>
  );
}