import axios from "axios";

const API_URL = "http://localhost:8080/ruang-rapat";

export function getEvents(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addEvent(data) {
  return axios
    .post(`${API_URL}`, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updateEvent(id, data) {
  return axios
    .put(`${API_URL}/${id}`, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deleteEvent(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}
