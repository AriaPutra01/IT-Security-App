import axios from "axios";

// export const getSags = axios.create({
//   baseURL: "http://localhost:1933", // Ganti dengan URL endpoint Golang Anda
//   timeout: 10000,
// });

export const getSag = axios.create({
  baseURL: "http://localhost:1933",
});


export const fetchDataSag = async () => {
  try {
    const response = await getSag.get('/sag');
    return response.data;
  } catch (error) {
    console.error(error);
    throw error;
  }
};