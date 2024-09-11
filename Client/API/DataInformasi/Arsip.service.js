import axios from "axios";

const API_URL = "http://localhost:8080/Arsip"; // Ganti dengan URL API Anda

export const getArsip = (callback) => {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data.arsip);
    })
    .catch((err) => {
      alert("Error fetching arsip:", err);
    });
};

export const addArsip = (data) => {
  return axios
    .post(`${API_URL}`, data)
    .then((response) => {
      return response.data.arsip;
    })
    .catch((error) => {
      alert("Error adding arsip:", error);
    });
};

export const updateArsip = (id, data) => {
  return axios
    .put(`${API_URL}/${id}`, data)
    .then((response) => {
      return response.data.arsip;
    })
    .catch((error) => {
      alert("Error updating arsip:", error);
    });
};

export const deleteArsip = (id) => {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      alert("Error deleting arsip:", error);
    });
};
