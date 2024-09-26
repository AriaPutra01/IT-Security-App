import axios from "axios";

const API_URL = "http://localhost:8080/Perdin";

export function getPerdins(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.perdin);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addPerdin(data) {
  const { username, ...rest } = data;
  return axios
    .post(`${API_URL}`, { ...rest })
    .then((response) => {
      return response.data.perdin;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updatePerdin(id, data) {
  const { username,...rest } = data;
  return axios
    .put(`${API_URL}/${id}`, { ...rest })
    .then((response) => {
      return response.data.perdin;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deletePerdin(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data.perdin;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}