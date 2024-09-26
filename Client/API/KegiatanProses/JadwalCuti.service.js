import axios from "axios";

const API_URL = "http://localhost:8080/jadwal-cuti";

export function getCutis(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.cuti);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addCuti(data) {
  return axios
    .post(`${API_URL}`, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function deleteCuti(id) {return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}
