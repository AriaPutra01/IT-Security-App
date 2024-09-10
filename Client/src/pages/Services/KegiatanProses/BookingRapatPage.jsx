import App from "../../../components/Layouts/App";
import { Calendar } from "../../../components/Fragments/Services/Calendar/Calendar";
import {
  getBookingRapat,
  addBookingRapat,
  deleteBookingRapat,
} from "../../../../API/KegiatanProses/BookingRapat.service";

export function BookingRapatPage() {
  return (
    <App services="Booking Ruang Rapat">
      <Calendar
        view="dayGridMonth"
        get={getBookingRapat}
        add={addBookingRapat}
        remove={deleteBookingRapat}
      />
    </App>
  );
}
