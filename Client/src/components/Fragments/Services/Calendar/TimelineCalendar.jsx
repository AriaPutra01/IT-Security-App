import { useEffect, useReducer, useState, useRef } from "react";
import { Button, Label, Modal, TextInput } from "flowbite-react";
import {
  Scheduler,
  SchedulerData,
  ViewType,
  DATE_FORMAT,
  AddMorePopover,
  wrapperFun,
} from "react-big-schedule";
import "../../../../calendar.css";
import dayjs from "dayjs";
import moment from "moment";
import Swal from "sweetalert2";
import { ColorPick } from "../../../../Utilities/ColorPick";

const initialState = {
  showScheduler: false,
  viewModel: {},
};

function reducer(state, action) {
  switch (action.type) {
    case "INITIALIZE":
      return { showScheduler: true, viewModel: action.payload };
    case "UPDATE_SCHEDULER":
      return { ...state, viewModel: action.payload };
    default:
      return state;
  }
}

function Timeline({
  getEvents,
  insertEvent,
  removeEvent,
  getResources,
  insertResource,
  removeResource,
}) {
  const [formModalOpen, setFormModalOpen] = useState(false);
  const [formData, setFormData] = useState({});
  const [state, dispatch] = useReducer(reducer, initialState);
  const [popoverState, setPopoverState] = useState({
    headerItem: undefined,
    left: 0,
    top: 0,
    height: 0,
  });

  const parentRef = useRef(null);

  useEffect(() => {
    const schedulerData = new SchedulerData(
      new dayjs().format(DATE_FORMAT),
      ViewType.Month,
      false,
      false,
      {
        resourceName: "Project",
        responsiveByParent: true,
        customCellWidth: 30,
        views: [
          {
            viewName: "Bulan",
            viewType: ViewType.Month,
            showAgenda: false,
            isEventPerspective: false,
          },
        ],
        schedulerMaxHeight: 440,
        dayMaxEvents: 2,
        weekMaxEvents: 4,
        monthMaxEvents: 4,
        quarterMaxEvents: 4,
        yearMaxEvents: 4,
      }
    );
    // schedulerData.localeDayjs.locale("en");
    //
    moment.locale("id", {
      months:
        "Januari_Februari_Maret_April_Mei_Juni_Juli_Agustus_September_Oktober_November_Desember".split(
          "_"
        ),
      monthsShort: "Jan_Feb_Mar_Apr_Mei_Jun_Jul_Agu_Sep_Okt_Nov_Des".split("_"),
      monthsParseExact: true,
      weekdays: "Minggu_Senin_Selasa_Rabu_Kamis_Jumat_Sabtu".split("_"),
      weekdaysShort: "Min_Sen_Sel_Rab_Kam_Jum_Sab".split("_"),
      weekdaysMin: "Mg_Sn_Sl_Rb_Km_Jm_Sb".split("_"),
      weekdaysParseExact: true,
      longDateFormat: {
        LT: "HH:mm",
        LTS: "HH:mm:ss",
        L: "DD/MM/YYYY",
        LL: "D MMMM YYYY",
        LLL: "D MMMM YYYY HH:mm",
        LLLL: "dddd, D MMMM YYYY HH:mm",
      },
      calendar: {
        sameDay: "[Hari ini pukul] LT",
        nextDay: "[Besok pukul] LT",
        nextWeek: "dddd [pukul] LT",
        lastDay: "[Kemarin pukul] LT",
        lastWeek: "dddd [lalu pukul] LT",
        sameElse: "L",
      },
      relativeTime: {
        future: "dalam %s",
        past: "%s yang lalu",
        s: "beberapa detik",
        m: "semenit",
        mm: "%d menit",
        h: "sejam",
        hh: "%d jam",
        d: "sehari",
        dd: "%d hari",
        M: "sebulan",
        MM: "%d bulan",
        y: "setahun",
        yy: "%d tahun",
      },
      dayOfMonthOrdinalParse: /\d{1,2}/,
      ordinal: function (number) {
        return number;
      },
      meridiemParse: /pagi|siang|sore|malam/,
      isPM: function (input) {
        return /^(siang|sore|malam)$/.test(input);
      },
      meridiem: function (hours, minutes, isLower) {
        if (hours < 11) {
          return "pagi";
        } else if (hours < 15) {
          return "siang";
        } else if (hours < 19) {
          return "sore";
        } else {
          return "malam";
        }
      },
      week: {
        dow: 1, // Senin adalah hari pertama dalam seminggu.
        doy: 4, // Digunakan untuk menentukan minggu pertama dalam setahun.
      },
    });
    //
    getEvents((eventData) => {
      getResources((resourceData) => {
        schedulerData.setResources(resourceData);
        schedulerData.setEvents(eventData);
        dispatch({ type: "INITIALIZE", payload: schedulerData });
      });
    });
  }, []);

  const prevClick = (schedulerData) => {
    schedulerData.prev();
    getEvents((data) => {
      schedulerData.setEvents(data);
      dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
    });
  };

  const nextClick = (schedulerData) => {
    schedulerData.next();
    getEvents((data) => {
      schedulerData.setEvents(data);
      dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
    });
  };

  const onViewChange = (schedulerData, view) => {
    schedulerData.setViewType(
      view.viewType,
      view.showAgenda,
      view.isEventPerspective
    );
    getEvents((data) => {
      schedulerData.setEvents(data);
      dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
    });
  };

  const onSelectDate = (schedulerData, date) => {
    schedulerData.setDate(date);
    getEvents((data) => {
      schedulerData.setEvents(data);
      dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
    });
  };

  const onCloseFormModal = () => {
    setFormModalOpen(false);
    setFormData({});
  };

  const handleFormChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  };

  const newEvent = (
    schedulerData,
    slotId,
    slotName,
    start,
    end,
    type,
    item
  ) => {
    setFormModalOpen(true);
    setFormData({
      schedulerData,
      title: "",
      start,
      end,
      slotId,
      slotName,
      bgColor: "#2596be",
      type,
      item,
    });
  };

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    const newEvent = {
      title: formData.title,
      start: formData.start,
      end: formData.end,
      resourceId: formData.slotId,
      bgColor: formData.bgColor,
    };
    try {
      const response = await insertEvent(newEvent);
      state.viewModel.addEvent(response);
      dispatch({ type: "UPDATE_SCHEDULER", payload: state.viewModel });
      onCloseFormModal();
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Gagal!",
        text: "Error saat menyimpan data: " + error,
        showConfirmButton: false,
        timer: 1500,
      });
    }
  };

  const updateEventStart = (schedulerData, event, newStart) => {
    if (
      confirm(
        `Do you want to adjust the start of the event? {eventId: ${event.id}, eventTitle: ${event.title}, newStart: ${newStart}}`
      )
    ) {
      schedulerData.updateEventStart(event, newStart);
    }
    dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
  };

  const updateEventEnd = (schedulerData, event, newEnd) => {
    if (
      confirm(
        `Do you want to adjust the end of the event? {eventId: ${event.id}, eventTitle: ${event.title}, newEnd: ${newEnd}}`
      )
    ) {
      schedulerData.updateEventEnd(event, newEnd);
    }
    dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
  };

  const moveEvent = (schedulerData, event, slotId, slotName, start, end) => {
    if (
      confirm(
        `Do you want to move the event? {eventId: ${event.id}, eventTitle: ${event.title}, newSlotId: ${slotId}, newSlotName: ${slotName}, newStart: ${start}, newEnd: ${end}`
      )
    ) {
      schedulerData.moveEvent(event, slotId, slotName, start, end);
      dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
    }
  };

  const deleteEvent = (schedulerData, event) => {
    Swal.fire({
      title: "Apakah Anda yakin?",
      text: `Anda akan menghapus event "${event.title}" ?`,
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, saya yakin",
      cancelButtonText: "Batal",
    }).then((result) => {
      if (result.isConfirmed) {
        removeEvent(event.id)
          .then(() => {
            schedulerData.removeEvent(event);
            dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
          })
          .catch((error) => {
            alert(`Error deleting event: ${error.message}`);
          });
      }
    });
  };

  const addResource = async (schedulerData) => {
    const { value: name } = await Swal.fire({
      title: "Masukan Resource!",
      input: "text",
      inputAttributes: {
        autocapitalize: "off",
      },
      showCancelButton: true,
      confirmButtonText: "Simpan",
      showLoaderOnConfirm: true,
      preConfirm: (e) => {
        return {
          name: e,
        };
      },
    });
    if (name) {
      try {
        const newResource = await insertResource(name);
        schedulerData.addResource(newResource);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      } catch (err) {
        Swal.fire({
          icon: "error",
          title: "Gagal!",
          text: "Error saat menyimpan Resource: " + err.message,
          showConfirmButton: false,
          timer: 1500,
        });
      }
    }
  };

  const slotClickedFunc = async (schedulerData, slot) => {
    const confirmResult = await Swal.fire({
      title: "Apakah Anda yakin?",
      text: "Anda tidak akan dapat mengembalikan ini!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonColor: "#3085d6",
      cancelButtonColor: "#d33",
      confirmButtonText: "Ya, hapus saja!",
    });

    if (confirmResult.isConfirmed) {
      try {
        // Hapus resource dari server
        await removeResource(slot.slotId);

        // Hapus event yang memiliki resourceId yang sama dari server
        const eventsToDelete = schedulerData.events.filter(
          (event) => event.resourceId === slot.slotId
        );

        for (const event of eventsToDelete) {
          await removeEvent(event.id);
        }

        // Hapus resource dan event dari client
        const updatedResources = schedulerData.resources.filter(
          (resource) => resource.id !== slot.slotId
        );
        const updatedEvents = schedulerData.events.filter(
          (event) => event.resourceId !== slot.slotId
        );
        schedulerData.setResources(updatedResources);
        schedulerData.setEvents(updatedEvents);

        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      } catch (error) {
        Swal.fire(
          "Gagal!",
          `Gagal menghapus resource atau event: ${error.message}`,
          "error"
        );
      }
    }
  };

  const onSetAddMoreState = (newState) => {
    if (newState === undefined) {
      setPopoverState({
        headerItem: undefined,
        left: 0,
        top: 0,
        height: 0,
      });
    } else {
      setPopoverState({
        ...newState,
      });
    }
  };

  const toggleExpandFunc = (schedulerData, slotId) => {
    schedulerData.toggleExpandStatus(slotId);
    dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
  };

  let popover = <div />;
  if (popoverState.headerItem !== undefined) {
    popover = (
      <AddMorePopover
        headerItem={popoverState.headerItem}
        eventItemClick={deleteEvent}
        schedulerData={state.viewModel}
        closeAction={onSetAddMoreState}
        left={popoverState.left}
        top={popoverState.top}
        height={popoverState.height}
        moveEvent={moveEvent}
      />
    );
  }

  const leftCustomHeader = (
    <div>
      <span style={{ fontWeight: "bold" }}>
        <Button className="ml-2" onClick={() => addResource(state.viewModel)}>
          Tambah Resource
        </Button>
      </span>
    </div>
  );

  return (
    <div className="w-full flex justify-start" ref={parentRef}>
      {" "}
      {/* Tambahkan ref ke elemen induk */}
      {state.showScheduler && (
        <div>
          <Scheduler
            schedulerData={state.viewModel}
            prevClick={prevClick}
            nextClick={nextClick}
            onSelectDate={onSelectDate}
            onViewChange={onViewChange}
            eventItemClick={deleteEvent}
            updateEventStart={updateEventStart}
            updateEventEnd={updateEventEnd}
            moveEvent={moveEvent}
            newEvent={newEvent}
            onSetAddMoreState={onSetAddMoreState}
            toggleExpandFunc={toggleExpandFunc}
            slotClickedFunc={slotClickedFunc}
            leftCustomHeader={leftCustomHeader}
            parentRef={parentRef} // Teruskan ref ke Scheduler
          />
          {popover}
        </div>
      )}
      {/* ModalForm */}
      <Modal show={formModalOpen} size="xl" onClose={onCloseFormModal} popup>
        <Modal.Header />
        <Modal.Body>
          <form onSubmit={handleFormSubmit}>
            <div className="flex flex-col gap-4">
              <div className="flex flex-col">
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
              <div className="flex flex-col gap-2 justify-start">
                <Label htmlFor="color" value="Color" />
                <ColorPick
                  colors={[
                    {
                      id: "blue",
                      label: "Blue",
                      value: "#4285f4",
                      checked: true,
                    },
                    { id: "red", label: "Red", value: "#db4437" },
                    { id: "yellow", label: "Yellow", value: "#fbbc05" },
                    { id: "green", label: "Green", value: "#0f9d58" },
                    { id: "teal", label: "Teal", value: "#00bfa5" },
                    { id: "purple", label: "Purple", value: "#9c27b0" },
                    { id: "pink", label: "Pink", value: "#e91e63" },
                  ]}
                  name="bgColor"
                  value={formData.bgColor}
                  onChange={handleFormChange}
                  className="mb-2 p-[2px]"
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
}

export default wrapperFun(Timeline);
