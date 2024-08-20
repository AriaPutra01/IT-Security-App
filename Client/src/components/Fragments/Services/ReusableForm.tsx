import React, { useState } from "react";
import { Button, Label, TextInput } from "flowbite-react";

interface FormField {
  name: string;
  label: string;
  type: string;
  required: boolean;
}
interface FormConfig {
  fields: FormField[];
  onSubmit: (data: any) => void;
  action: string;
  services: string;
}
interface ReusableFormProps {
  config: FormConfig;
  formData: any;
  setFormData: (data: any) => void;
}

export const ReusableForm = ({
  config,
  formData,
  setFormData,
}: ReusableFormProps) => {
  const [errors, setErrors] = useState({});

  const handleInputChange = ({
    target: { name, value },
  }: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
  };

  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    // Validasi untuk setiap field
    config.fields.forEach((field) => {
      const value = formData[field.name]?.trim();
      if (field.required && !value) {
        newErrors[field.name] = `${field.label} is required.`;
      } else if (field.required && value.length === 0) {
        newErrors[field.name] = `${field.label} cannot be just whitespace.`;
      }
    });
    // Atur state error jika ada
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (validateForm()) {
      config.onSubmit(formData);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <h3 className="flex gap-1 justify-center text-xl font-medium text-gray-900 dark:text-white">
        {config.action === "add" ? "Tambah Data" : `Ubah Data`}
        <div className="uppercase">{config.services}</div>
      </h3>
      <div className="grid grid-cols-4 gap-6">
        {config.fields.map((field, index) => (
          <div className="col-span-2" key={index}>
            <div className="mb-2 block">
              <Label htmlFor={field.name} value={field.label} />
            </div>
            {field.type === "date" ? (
              <TextInput
                id={field.name}
                name={field.name}
                type={field.type}
                value={
                  formData[field.name]
                    ? new Date(formData[field.name]).toISOString().split("T")[0]
                    : ""
                }
                onChange={handleInputChange}
                required={field.required}
              />
            ) : (
              <TextInput
                id={field.name}
                name={field.name}
                type={field.type}
                value={formData[field.name]}
                onChange={handleInputChange}
                required={field.required}
              />
            )}
            {errors[field.name] && (
              <p className="text-red-600 text-sm">{errors[field.name]}</p>
            )}
          </div>
        ))}
      </div>
      <Button
        className="w-full"
        color={config.action === "add" ? "info" : "warning"}
        type="submit"
      >
        {config.action === "add" ? "Tambah" : "Ubah"}
      </Button>
    </form>
  );
};
