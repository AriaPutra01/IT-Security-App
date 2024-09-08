import axios from "axios";

const API_URL = "http://localhost:8080/booking-rapat";

export function getBookingRapat(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addBookingRapat(data) {
  return axios
    .post(`${API_URL}`, {
      ...data,
      color: data.color // Pastikan warna termasuk dalam data yang dikirim
    })
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function deleteBookingRapat(id) {
  if (!id) {
    throw new Error("ID harus disertakan untuk menghapus data.");
  }
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}