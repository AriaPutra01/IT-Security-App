import axios from "axios";

const API_URL = "http://localhost:8080/booking-rapat";

export function getBookingRapat(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.booking);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addBookingRapat(data) {
  return axios
    .post(`${API_URL}`, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function deleteBookingRapat(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}
