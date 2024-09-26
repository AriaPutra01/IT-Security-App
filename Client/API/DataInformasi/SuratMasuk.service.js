import axios from "axios";

const API_URL = "http://localhost:8080/SuratMasuk";

export function getSuratMasuks(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.SuratMasuk);    
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addSuratMasuk(data) {
  const { username, ...rest } = data;
  return axios
    .post(`${API_URL}`, { ...rest })
    .then((response) => {
      return response.data.SuratMasuk;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updateSuratMasuk(id, data) {
  const { username, ...rest } = data;
  return axios
    .put(`${API_URL}/${id}`, { ...rest })
    .then((response) => {
      return response.data.SuratMasuk;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deleteSuratMasuk(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data.SuratMasuk;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}