import axios from "axios";

const API_URL = "http://localhost:8080/timelineDesktop";
const RESOURCE_API_URL = "http://localhost:8080/resourceDesktop";

export function getEventsDesktop(callback) {
  return axios
    .get(API_URL)
    .then((response) => {
      callback(response.data.events);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addEventDesktop(data) {
  return axios
    .post(API_URL, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function deleteEventDesktop(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}

export function getResources(callback) {
  return axios
    .get(RESOURCE_API_URL)
    .then((response) => {
      callback(response.data.resources);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil resource. Alasan: ${error.message}`);
    });
}

export function addResource(data) {
  return axios
    .post(RESOURCE_API_URL, data)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan resource. Alasan: ${error.message}`);
    });
}

export function deleteResource(id) {
  return axios
    .delete(`${RESOURCE_API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus resource. Alasan: ${error.message}`);
    });
}
