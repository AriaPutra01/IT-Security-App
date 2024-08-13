import axios from "axios";

export const getSags = (callback) => {
  axios
    .get("http://localhost:8080/sag")
    .then((res) => {
      callback(res.data.posts);
    })
    .catch((err) => {
      console.log(err);
    });
};

export async function deleteSag(id) {
  try {
    const response = await axios.delete(`http://localhost:8080/sag/${id}`);
    return response.data;
  } catch (error) {
    throw new Error(
      `Gagal hapus SAG dengan id = ${id}. Alasan: ${error.message}`
    );
  }
}
