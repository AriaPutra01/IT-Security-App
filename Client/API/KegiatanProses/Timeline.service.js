import axios from "axios";

const API_URL = "http://localhost:8080/timelines";
const RESOURCE_API_URL = "http://localhost:8080/resources";

export function getTimelines() {
  return axios
    .get(API_URL)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addTimeline(data) {
  console.log("Data yang dikirim ke server:", data); // Tambahkan log ini
  return axios
    .post(API_URL, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      console.error("Error response:", error.response); // Tambahkan log ini
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function deleteTimeline(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}

export function getResourcesTimeline() {
  return axios
    .get(RESOURCE_API_URL)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil resource. Alasan: ${error.message}`);
    });
}

export function addResourceTimeline(data) {
  return axios
    .post(RESOURCE_API_URL, data)
    .then((response) => {
      return response.data; // Pastikan response data memiliki struktur yang benar
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan resource. Alasan: ${error.message}`);
    });
}

export function deleteResourceTimeline(id) {
  return axios
    .delete(`${RESOURCE_API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus resource. Alasan: ${error.message}`);
    });
}
