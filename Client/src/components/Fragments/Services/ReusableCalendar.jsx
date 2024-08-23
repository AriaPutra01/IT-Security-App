import FullCalendar from "@fullcalendar/react";
import { formatDate } from "@fullcalendar/core";
import dayGridPlugin from "@fullcalendar/daygrid";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
import listPlugin from "@fullcalendar/list";
import "../../../calendar.css";

export const ReusableCalendar = (props) => {
  const { currentEvents, handleDateClick, handleEventClick } = props;
  return (
    <div className="calendar-container">
      <div className="calendar-content dark:text-white">
        <div className="events-list dark:bg-gray-800">
          <h2>Events</h2>
          <ul>
            {currentEvents.map((event) => (
              <li key={event.id} className="event-item">
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
  );
};
