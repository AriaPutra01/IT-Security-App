import axios from "axios";

const API_URL = "http://localhost:8080/user";

export function getUsers(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.users);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function updateUser(id, data) {
  return axios
    .put(`${API_URL}/${id}`, data)
    .then((response) => {
      return response.data.users;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deleteUser(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data.users;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}
