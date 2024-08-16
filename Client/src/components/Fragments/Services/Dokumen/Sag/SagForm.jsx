  import { Button, Label, TextInput } from "flowbite-react";
  import React from "react";

  export const SagForm = (props) => {
    const {
      onSubmit,
      action,
      services,
      tanggal,
      noMemo,
      perihal,
      pic,
      setTanggal,
      setNoMemo,
      setPerihal,
      setPic,
    } = props;

    return (
      <form onSubmit={onSubmit} className="space-y-6">
        <h3 className="flex gap-1 justify-center text-xl font-medium text-gray-900 dark:text-white">
          {action === "add" ? "Tambah Data" : `Ubah Data`}
          <div className="uppercase">{services}</div>
        </h3>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="Tanggal" value="Tanggal" />
          </div>
          <TextInput
            type="date"
            name="Tanggal"
            id="Tanggal"
            value={tanggal}
            onChange={(e) => setTanggal(e.target.value)}
            required
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="NoMemo" value="Nomor Memo" />
          </div>
          <TextInput
            id="NoMemo"
            name="NoMemo"
            type="text"
            value={noMemo}
            onChange={(e) => setNoMemo(e.target.value)}
            required
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="Perihal" value="Perihal" />
          </div>
          <TextInput
            id="Perihal"
            name="Perihal"
            type="text"
            value={perihal}
            onChange={(e) => setPerihal(e.target.value)}
            required
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="Pic" value="Pic" />
          </div>
          <TextInput
            id="Pic"
            name="Pic"
            type="text"
            value={pic}
            onChange={(e) => setPic(e.target.value)}
            required
          />
        </div>
        <Button
          className="w-full"
          color={action === "add" ? "info" : "warning"}
          type="submit"
        >
          {action === "add" ? "Tambah" : "Ubah"}
        </Button>
      </form>
    );
  };
