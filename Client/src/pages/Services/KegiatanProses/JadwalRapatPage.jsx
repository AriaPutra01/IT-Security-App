import App from "../../../components/Layouts/App";
import { Calendar } from "../../../components/Fragments/Services/Calendar/Calendar";
import {
  getRapats,
  addRapat,
  deleteRapat,
} from "../../../../API/KegiatanProses/JadwalRapat.service";

export function JadwalRapatPage() {
  return (
    <App services="Jadwal Rapat">
      <Calendar
        view="timeGridWeek"
        get={getRapats}
        add={addRapat}
        remove={deleteRapat}
      />
    </App>
  );
}
