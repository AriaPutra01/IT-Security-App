import React from "react";
import App from "../../../components/Layouts/App";
import Timeline from "../../../components/Fragments/Services/Calendar/TimelineCalendar";
import {
  getEventsDesktop,
  addEventDesktop,
  deleteEventDesktop,
  getResources,
  addResource,
  deleteResource,
} from "../../../../API/KegiatanProses/TimelineDesktop.service";

export const TimelineDesktopPage = () => {
  return (
    <App services="Timeline Wallpaper Desktop">
      <Timeline
        getEvents={getEventsDesktop}
        insertEvent={addEventDesktop}
        removeEvent={deleteEventDesktop}
        getResources={getResources}
        insertResource={addResource}
        removeResource={deleteResource}
      />
    </App>
  );
};
