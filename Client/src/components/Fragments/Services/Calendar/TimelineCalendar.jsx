import { useEffect, useReducer, useState } from "react";
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

  useEffect(() => {
    const schedulerData = new SchedulerData(
      new dayjs().format(DATE_FORMAT),
      ViewType.Month,
      false,
      false,
      {
        schedulerWidth: "94%",
        responsiveByParent: true,
        customCellWidth: 30,
        nonAgendaDayCellHeaderFormat: "M/D|HH:mm",
        views: [
          {
            viewName: "Day",
            viewType: ViewType.Day,
            showAgenda: false,
            isEventPerspective: false,
          },
          {
            viewName: "Week",
            viewType: ViewType.Week,
            showAgenda: false,
            isEventPerspective: false,
          },
          {
            viewName: "Month",
            viewType: ViewType.Month,
            showAgenda: false,
            isEventPerspective: false,
          },
        ],
        schedulerMaxHeight: 400,
        dayMaxEvents: 2,
        weekMaxEvents: 4,
        monthMaxEvents: 4,
        quarterMaxEvents: 4,
        yearMaxEvents: 4,
      }
    );
    schedulerData.localeDayjs.locale("id");
    getEvents((eventData) => {
      getResources((resourceData) => {
        schedulerData.setResources(resourceData || []);
        schedulerData.setEvents(eventData || []);
        dispatch({ type: "INITIALIZE", payload: schedulerData });
      });
    });
  }, []);

  const prevClick = (schedulerData) => {
    schedulerData.prev();
    getEvents()
      .then((data) => {
        schedulerData.setEvents(data || []);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching timelines: ${error.message}`);
      });
  };

  const nextClick = (schedulerData) => {
    schedulerData.next();
    getEvents()
      .then((data) => {
        schedulerData.setEvents(data || []);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching timelines: ${error.message}`);
      });
  };

  const onViewChange = (schedulerData, view) => {
    schedulerData.setViewType(
      view.viewType,
      view.showAgenda,
      view.isEventPerspective
    );
    getEvents()
      .then((data) => {
        schedulerData.setEvents(data || []);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching timelines: ${error.message}`);
      });
  };

  const onSelectDate = (schedulerData, date) => {
    schedulerData.setDate(date);
    getEvents()
      .then((data) => {
        schedulerData.setEvents(data || []);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching timelines: ${error.message}`);
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
    let newFreshId = 0;
    schedulerData.events.forEach((item) => {
      if (item.id >= newFreshId) newFreshId = item.id + 1;
    });
    setFormModalOpen(true);
    setFormData({
      newFreshId,
      title: "",
      start,
      end,
      slotId,
      bgColor: "#2596be",
    });
  };

  const handleFormSubmit = (e) => {
    e.preventDefault();
    const newEvent = {
      id: formData.newFreshId,
      title: formData.title,
      start: formData.start,
      end: formData.end,
      resourceId: formData.slotId,
      bgColor: formData.bgColor,
    };
    try {
      insertEvent(newEvent);
      state.viewModel.addEvent(newEvent);
      dispatch({ type: "UPDATE_SCHEDULER", payload: state.viewModel });
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
        insertResource(name);
        schedulerData.addResource(name);
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
        const eventsToDelete = schedulerData.filter(
          (event) => event.resourceId === slot.slotId
        );
        for (const event of eventsToDelete) {
          await removeEvent(event.id);
        }

        // Hapus resource dan event dari client
        const updatedResources = schedulerData.filter(
          (resource) => resource.id !== slot.slotId
        );
        const updatedEvents = schedulerData.filter(
          (event) => event.resourceId !== slot.slotId
        );
        schedulerData.setResources(updatedResources);
        schedulerData.setEvents(updatedEvents);

        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
        Swal.fire(
          "Terhapus!",
          "Resource dan event terkait telah dihapus.",
          "success"
        );
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
    <div className="w-full flex justify-start">
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
          />
          {popover}
        </div>
      )}
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
                  className="block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 appearance-none  focus:outline-none focus:ring-0 focus:border-blue-600 peer"
                  placeholder="masukan event"
                  value={formData.title}
                  onChange={handleFormChange}
                  required
                />
              </div>
              <div className="flex flex-col gap-2 justify-start col-span-1">
                <Label htmlFor="bgColor" value="Color" />
                <ColorPick
                  name="bgColor"
                  value={formData.bgColor}
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
}

export default wrapperFun(Timeline);
