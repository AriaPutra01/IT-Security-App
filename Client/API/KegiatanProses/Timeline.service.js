import axios from "axios";

const API_URL = "http://localhost:8080/timeline";

export function getTimelines(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addTimeline(data) {
  return axios
    .post(`${API_URL}`, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function deleteTimeline(id) {
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