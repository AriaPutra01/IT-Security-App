import React from "react";
import App from "../../../components/Layouts/App";
import Timeline from "../../../components/Fragments/Services/Calendar/TimelineCalendar";
import {
  getEventsProject,
  addEventProject,
  deleteEventProject,
  getResources,
  addResource,
  deleteResource,
} from "../../../../API/KegiatanProses/TimelineProject.service";

export const TimelineProjectPage = () => {
  return (
    <App services="Timeline Project">
      <Timeline
        getEvents={getEventsProject}
        insertEvent={addEventProject}
        removeEvent={deleteEventProject}
        getResources={getResources}
        insertResource={addResource}
        removeResource={deleteResource}
      />
    </App>
  );
};
