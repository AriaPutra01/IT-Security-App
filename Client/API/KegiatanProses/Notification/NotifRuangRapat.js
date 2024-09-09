import axios from "axios";

const API_URL = "http://localhost:8080/notifications";

export const GetNotifRapat = (callback) => {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.notifications); // ensure this matches the structure sent from the backend
    })
    .catch((error) => {
      throw new Error(`Gagal menampilkan notif. Alasan: ${error.message}`);
    });
};

export function deleteNotifRapat(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data.notifications;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus notif. Alasan: ${error.message}`);
    });
}
