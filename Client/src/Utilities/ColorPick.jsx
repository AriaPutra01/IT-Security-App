import React from "react";

export const ColorPick = (props) => {
  const { name, onChange, className, value } = props;
  return (
    <input
      value={value}
      className={className}
      type="color"
      id={name}
      name={name}
      onChange={onChange}
    />
  );
};
