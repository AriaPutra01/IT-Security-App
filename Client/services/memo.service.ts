import axios from "axios";

export const getIso = (callback) => {
  axios
    .get("http://localhost:8080/iso")
    .then((res) => {
      callback(res.data.posts);
      // console.log(res.data.posts);
    })
    .catch((err) => {
      console.log(err);
    });
};

export async function deleteIso(id) {
  try {
    const response = await axios.delete(`http://localhost:8080/iso/${id}`);
    return response.data;
  } catch (error) {
    throw new Error(
      `Gagal hapus ISO dengan id = ${id}. Alasan: ${error.message}`
    );
  }
}
