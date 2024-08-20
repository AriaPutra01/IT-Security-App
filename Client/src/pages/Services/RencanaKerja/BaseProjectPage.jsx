import React from "react";
import App from "../../../components/Layouts/App";

export const BaseProjectPage = () => {
  return (
    <App services="Base Project">
      <div className="flex justify-center">
        <img
          className="rounded-xl w-7/12 border-4 border-black"
          src="../../../../public/images/image.png"
          alt=""
        />
      </div>
    </App>
  );
};
