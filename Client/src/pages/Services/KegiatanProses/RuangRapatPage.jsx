import { useState, useEffect } from "react";
import FullCalendar from "@fullcalendar/react";
import { formatDate } from "@fullcalendar/core";
import dayGridPlugin from "@fullcalendar/daygrid";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
import listPlugin from "@fullcalendar/list";
import axios from 'axios';
import { v4 as uuidv4 } from 'uuid'; // Import UUID
import "../../../calender.css";

export function RuangRapatPage() {
  const [currentEvents, setCurrentEvents] = useState([]);

  // Fetch events from API

  function getEvents(callback) {
    return axios
      .get("http://localhost:8080/ruang-rapat")
      .then((response) => {
        callback(response.data);
      })
      .catch((error) => {
        throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
      });
  }

  // Fetch events
  useEffect(() => {
    getEvents((data) => {
      setCurrentEvents(data);
    });
  }, []);

  // Save event to backend
  const saveEvent = async (event) => {
    try {
      await axios.post('http://localhost:8080/ruang-rapat', event)
        .then((response) => {
          return response.data
        })
      setCurrentEvents([...currentEvents, event])
    } catch (error) {
      console.error("Error saving event:", error);
    }
  };

  // Delete event from backend
  const deleteEvent = async (eventId) => {
    try {
      await axios.delete(`http://localhost:8080/ruang-rapat?id=${eventId}`);
      getEvents(); // Fetch events again to update the UI
    } catch (error) {
      console.error("Error deleting event:", error);
    }
  };

  // Handle date click to add new event
  const handleDateClick = async (selected) => {
    const title = prompt("Please enter a new title for your event");
    const calendarApi = selected.view.calendar;
    calendarApi.unselect();

    if (title) {
      const newEvent = {
        id: uuidv4(),
        title,
        start: selected.startStr,
        end: selected.endStr,
        allDay: selected.allDay,
      };

      setCurrentEvents((prevEvents) => {
        const updatedEvents = [...prevEvents, newEvent];
        saveEvent(newEvent); // Save event to backend
        return updatedEvents;
      });
      calendarApi.addEvent(newEvent); // Add event to calendar UI
    }
  };

  // Handle event click to delete event
  const handleEventClick = async (selected) => {
    if (
      window.confirm(
        `Are you sure you want to delete the event '${selected.event.title}'`
      )
    ) {
      selected.event.remove();
      setCurrentEvents((prevEvents) => {
        const updatedEvents = prevEvents.filter(event => event.id !== selected.event.id);
        deleteEvent(selected.event.id); // Delete event from backend
        return updatedEvents;
      });
    }
  };

  return (
    <div className="calendar-container">
      <header>
        <h1>Ruang Rapat</h1>
      </header>

      <div className="calendar-content">
        <div className="events-list">
          <h2>Events</h2>
          <ul>
            {currentEvents.map((event, index) => (
              <li key={index} className="event-item">
                <div className="event-title">{event.Title}</div>
                <div className="event-date">
                  {formatDate(event.Start, {
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
            plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin, listPlugin]}
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
            events={currentEvents} // Ensure events are set correctly
          />
        </div>
      </div>
    </div>
  );
}