import axios from "axios";

const API_URL = "http://localhost:8080/sagiso";

export function getSagisos(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.sagiso.map(sagiso => ({
        ...sagiso,
        info: sagiso.info // Menyertakan kolom Info
      })));
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addSagiso(data) {
  const { username, ...rest } = data;
  return axios
    .post(`${API_URL}`, { ...rest, info: username })
    .then((response) => {
      return response.data.sagiso;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updateSagiso(id, data) {
  const { username, ...rest } = data;
  return axios
    .put(`${API_URL}/${id}`, { ...rest })
    .then((response) => {
      return response.data.sagiso;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deleteSagiso(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}