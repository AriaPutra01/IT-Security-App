import axios from "axios";

const API_URL = "http://localhost:8080/sag";

export function getSags(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.posts);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addSag(data) {
  return axios
    .post(`${API_URL}`, data)
    .then((response) => {
      return response.data.posts;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updateSag(id, data) {
  return axios
    .put(`${API_URL}/${id}`, data)
    .then((response) => {
      return response.data.posts;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deleteSag(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data.posts;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}