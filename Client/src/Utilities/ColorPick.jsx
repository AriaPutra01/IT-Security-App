import React from "react";

export const ColorPick = (props) => {
  const { checked, color, handleColor } = props;
  return (
    <input
      type="radio"
      defaultChecked={checked}
      name="color"
      className={`size-6 bg-[${color}] rounded-full cursor-pointer checked:bg-[${color}] checked:ring-[${color}] active:ring-[${color}] hover:scale-105`}
      onClick={() => handleColor(color)}
    />
  );
};
