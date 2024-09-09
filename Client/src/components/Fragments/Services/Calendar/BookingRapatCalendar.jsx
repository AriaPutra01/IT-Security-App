import FullCalendar from "@fullcalendar/react";
import { formatDate } from "@fullcalendar/core";
import dayGridPlugin from "@fullcalendar/daygrid";
import interactionPlugin from "@fullcalendar/interaction";
import listPlugin from "@fullcalendar/list";

export const BookingRapatCalendar = ({
  currentEvents,
  handleDateClick,
  handleEventClick,
  onColorChange,
}) => {
  return (
    <div className="m-5 grid grid-cols-2fr">
      <div className="bg-gray-50 p-[15px] rounded w-[200px] max-h-[80vh] overflow-auto">
        <h2 className="text-xl mt-0 mb-2 font-bold">
          {currentEvents.length} Jadwal
        </h2>
        <div className="flex flex-col gap-2">
          {currentEvents.map((event) => (
            <div
              key={event.ID}
              className="overflow-auto grow border-b-2 border-l-2 border-sky-500 shadow p-2 rounded"
            >
              <div className="text-sky-500 font-bold">{event.title}</div>
              <div>
                {formatDate(event.start, {
                  year: "numeric",
                  month: "short",
                  day: "numeric",
                })}
              </div>
            </div>
          ))}
        </div>
      </div>
      <div className="mx-4">
        <FullCalendar
          plugins={[dayGridPlugin, interactionPlugin, listPlugin]}
          headerToolbar={{
            left: "prev,next today",
            center: "title",
            right: "dayGridMonth,listMonth",
          }}
          initialView="dayGridMonth"
          editable={true}
          selectable={true}
          selectMirror={true}
          dayMaxEvents={true}
          events={currentEvents}
          select={handleDateClick}
          eventClick={handleEventClick}
        />
      </div>
      <div>
        <button onClick={() => onColorChange("red")}>Merah</button>
        <button onClick={() => onColorChange("green")}>Hijau</button>
        <button onClick={() => onColorChange("yellow")}>Kuning</button>
      </div>
    </div>
  );
};
