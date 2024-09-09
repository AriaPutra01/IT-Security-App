import { formatDate } from "@fullcalendar/core";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";
import interactionPlugin from "@fullcalendar/interaction";
import listPlugin from "@fullcalendar/list";

export const TimelineCalendar = ({
  children,
  currentEvents,
  handleDateClick,
  handleEventClick,
}) => {
  return (
    <div className="m-5 grid grid-cols-2fr">
      <div className="bg-gray-50 p-[15px] rounded w-[200px] max-h-[80vh] overflow-auto">
        <div className="flex gap-1 justify-between ">{children}</div>
        <h2 className="text-lg my-2 font-bold">
          {currentEvents.length} Jadwal
        </h2>
        <div className="flex flex-col gap-2">
          {currentEvents.map((event) => (
            <div
              key={event.id}
              className={`border-b-2 border-s-2 flex gap-3 items-center overflow-auto grow p-2 rounded`}
            >
              <div className={`size-4 rounded shadow bg-[${event.color}]`} />
              <div>
                <div className="font-bold">{event.title}</div>
                <div>
                  {formatDate(event.start, {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                  })}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
      <div className="mx-4 max-h-[80vh] overflow-auto border-b-2">
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
    </div>
  );
};
