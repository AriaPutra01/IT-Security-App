import axios from "axios";

const API_URL = "http://localhost:8080/beritaAcara";

export function getBeritaAcaras(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.beritaAcaras);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addBeritaAcara(data) {
  const { username, ...rest } = data;
  return axios
    .post(`${API_URL}`, { ...rest }) // Tambahkan info
    .then((response) => {
      return response.data.beritaAcara;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updateBeritaAcara(id, data) {
  const { username, ...rest } = data;
  return axios
    .put(`${API_URL}/${id}`, { ...rest }) // Kirim username sebagai bagian dari request
    .then((response) => {
      return response.data.beritaAcara;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function getBeritaAcaraShow(id, callback) {
  return axios
    .get(`${API_URL}/${id}`)
    .then((response) => {
      callback(response.data.beritaAcara);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function deleteBeritaAcara(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}