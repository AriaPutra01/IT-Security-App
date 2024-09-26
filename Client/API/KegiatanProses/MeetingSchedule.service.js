import axios from "axios";

const API_URL = "http://localhost:8080/meetingSchedule";

export function getMeetingList(callback) {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.meetingschedule);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
}

export function addMeetingList(data) {
  const { username, ...rest } = data;
  return axios
    .post(`${API_URL}`, { ...rest }) // Tambahkan info
    .then((response) => {
      return response.data.meetingschedule;
    })
    .catch((error) => {
      throw new Error(`Gagal menambahkan data. Alasan: ${error.message}`);
    });
}

export function updateMeetingList(id, data) {
  const { username, ...rest } = data;
  return axios
    .put(`${API_URL}/${id}`, { ...rest })
    .then((response) => {
      return response.data.meetingschedule;
    })
    .catch((error) => {
      throw new Error(`Gagal mengubah data. Alasan: ${error.message}`);
    });
}

export function deleteMeetingList(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus data. Alasan: ${error.message}`);
    });
}