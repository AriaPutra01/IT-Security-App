import { useEffect, useReducer, useState } from "react";
import { Button } from "flowbite-react";
import { v4 as uuidv4 } from "uuid"; // Import UUID
import {
  Scheduler,
  SchedulerData,
  ViewType,
  DATE_FORMAT,
  AddMorePopover,
  wrapperFun,
} from "react-big-schedule";
import "react-big-schedule/dist/css/style.css";
import {
  getTimelines,
  addTimeline,
  deleteTimeline,
  getResourcesTimeline,
  addResourceTimeline,
  deleteResourceTimeline,
} from "../../../../API/KegiatanProses/Timeline.service";
import dayjs from "dayjs";

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

function TimelinePage() {
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
      ViewType.Month
    );
    schedulerData.localeDayjs.locale("en");
    Promise.all([getTimelines(), getResourcesTimeline()])
      .then(([timelineData, resourceData]) => {
        schedulerData.setResources(resourceData.resources || []);
        schedulerData.setEvents(timelineData.events || []);
        dispatch({ type: "INITIALIZE", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching data: ${error.message}`);
      });
  }, []);

  const prevClick = (schedulerData) => {
    schedulerData.prev();
    getTimelines()
      .then((data) => {
        schedulerData.setEvents(data.events || []);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching timelines: ${error.message}`);
      });
  };

  const nextClick = (schedulerData) => {
    schedulerData.next();
    getTimelines()
      .then((data) => {
        schedulerData.setEvents(data.events || []);
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
    getTimelines()
      .then((data) => {
        schedulerData.setEvents(data.events || []);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching timelines: ${error.message}`);
      });
  };

  const onSelectDate = (schedulerData, date) => {
    schedulerData.setDate(date);
    getTimelines()
      .then((data) => {
        schedulerData.setEvents(data.events || []);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      })
      .catch((error) => {
        alert(`Error fetching timelines: ${error.message}`);
      });
  };

  const eventClicked = (schedulerData, event) => {
    alert(
      `You just clicked an event: {id: ${event.id}, title: ${event.title}}`
    );
  };

  const newEvent = async (
    schedulerData,
    slotId,
    slotName,
    start,
    end,
    type,
    item
  ) => {
    if (
      confirm(
        `Do you want to create a new event? {slotId: ${slotId}, slotName: ${slotName}, start: ${start}, end: ${end}, type: ${type}, item: ${item} }`
      )
    ) {
      const newEvent = {
        title: "Test oyyyy",
        start: start,
        end: end,
        resourceId: slotId,
        bgColor: "purple",
      };
      const response = await addTimeline(newEvent);
      if (response && response.id) {
        schedulerData.addEvent([...schedulerData.events, response]);
        dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
      } else {
        alert("Gagal menambahkan event");
      }
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

  const addResource = (schedulerData) => {
    const resourceName = prompt("Enter resource name:");
    if (resourceName) {
      const newResource = {
        name: resourceName,
      };
      addResourceTimeline(newResource)
        .then((data) => {
          if (data && data.id) {
            schedulerData.setResources([
              ...schedulerData.resources,
              { id: data.id, name: data.name },
            ]);
            dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
          } else {
            throw new Error("Invalid resource data");
          }
        })
        .catch((error) => {
          alert(`Error adding resource: ${error.message}`);
        });
    }
  };

  const deleteEvent = (schedulerData, event) => {
    if (
      confirm(
        `Do you want to delete the event? {eventId: ${event.id}, eventTitle: ${event.title}}`
      )
    ) {
      deleteTimeline(event.id)
        .then(() => {
          schedulerData.removeEvent(event);
          dispatch({ type: "UPDATE_SCHEDULER", payload: schedulerData });
        })
        .catch((error) => {
          alert(`Error deleting event: ${error.message}`);
        });
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
        eventItemClick={eventClicked}
        viewEventClick={deleteEvent}
        viewEventText="Delete"
        schedulerData={state.viewModel}
        closeAction={onSetAddMoreState}
        left={popoverState.left}
        top={popoverState.top}
        height={popoverState.height}
        moveEvent={moveEvent}
      />
    );
  }

  return (
    <>
      {state.showScheduler && (
        <div>
          <Button onClick={() => addResource(state.viewModel)}>
            Add Resource
          </Button>
          <Scheduler
            schedulerData={state.viewModel}
            prevClick={prevClick}
            nextClick={nextClick}
            onSelectDate={onSelectDate}
            onViewChange={onViewChange}
            eventItemClick={eventClicked}
            viewEventClick={deleteEvent}
            viewEventText="Delete"
            updateEventStart={updateEventStart}
            updateEventEnd={updateEventEnd}
            moveEvent={moveEvent}
            newEvent={newEvent}
            onSetAddMoreState={onSetAddMoreState}
            toggleExpandFunc={toggleExpandFunc}
          />
          {popover}
        </div>
      )}
    </>
  );
}

export default wrapperFun(TimelinePage);
