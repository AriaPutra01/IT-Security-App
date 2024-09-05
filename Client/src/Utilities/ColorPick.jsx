import React from "react";

export const ColorPick = ({ checked, color, onColorChange }) => {
  return (
    <input
      type="radio"
      defaultChecked={checked}
      name="color"
      className={`size-6 bg-[${color}] rounded-full cursor-pointer checked:bg-[${color}] checked:ring-[${color}] active:ring-[${color}] hover:scale-105`}
      onClick={() => onColorChange(color)} // Panggil fungsi onColorChange
    />
  );
};
