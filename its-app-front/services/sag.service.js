import axios from "axios";

export const getSags = (callback) => {
  axios
    .get("http://localhost:8080/sag")
    .then((res) => {
      // callback(res.data);
      callback(res.data.posts);
    })
    .catch((err) => {
      console.log(err);
    });
};
