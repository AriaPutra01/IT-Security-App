import App from "../../../components/Layouts/App";
import { Calendar } from "../../../components/Fragments/Services/Calendar/Calendar";
import {
  getCutis,
  addCuti,
  deleteCuti,
} from "../../../../API/KegiatanProses/JadwalCuti.service";

export function JadwalCutiPage() {
  return (
    <App services="Jadwal Cuti">
      <Calendar
        view="dayGridMonth"
        get={getCutis}
        add={addCuti}
        remove={deleteCuti}
      />
    </App>
  );
}
