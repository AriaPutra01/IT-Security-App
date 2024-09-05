import axios from "axios";

const API_URL = "http://localhost:8080/memos";

export function getMemos(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.posts);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addMemo(data) {
  const { usernane, ...rest} = data;
  return axios
    .post(`${API_URL}`, { ...rest })
    .then((response) => {
      return response.data.posts;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updateMemo(id, data) {
  const { username, ...rest} = data
  return axios
    .put(`${API_URL}/${id}`, { ...rest })
    .then((response) => {
      return response.data.posts;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deleteMemo(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data.posts;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}
