import axios from "axios";

export const GetNotifRapat = (callback) => {
  return axios
    .get(`http://localhost:8080/notifications`)
    .then((response) => {
      callback(response.data.notifications);
    })
    .catch((error) => {
      throw new Error(`Gagal mengambil data. Alasan: ${error.message}`);
    });
};
