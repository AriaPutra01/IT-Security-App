import React from "react";

export const ColorPick = ({ name, onChange, value, colors }) => {
  return (
    <div className="flex justify-around gap-2">
      {colors.map((color) => (
        <div key={color.id}>
          <input
            type="radio"
            id={color.id}
            name={name}
            value={color.value || value}
            defaultChecked={color.checked}
            onChange={onChange}
            style={{
              backgroundColor: color.value,
              width: "2rem",
              height: "2rem",
            }}
          />
        </div>
      ))}
    </div>
  );
};
