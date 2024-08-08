"use client";

import { Button, Label, TextInput } from "flowbite-react";
import { Link } from "react-router-dom";

export function LoginForm() {
  return (
    <form className="flex w-full flex-col gap-4">
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="email" value="Alamat Email" />
        </div>
        <TextInput
          id="email"
          type="email"
          placeholder="contoh@mail.com"
          required
        />
      </div>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="password" value="Password" />
        </div>
        <TextInput
          id="password"
          type="password"
          placeholder="••••••••"
          required
        />
      </div>
      <Button color="warning" type="submit">
        Kirim
      </Button>
    </form>
  );
}

export function RegisterForm() {
  return (
    <form className="flex w-full flex-col gap-4">
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="username" value="Username" />
        </div>
        <TextInput id="username" type="text" required />
      </div>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="email" value="Alamat Email" />
        </div>
        <TextInput
          id="email"
          type="email"
          placeholder="contoh@mail.com"
          required
        />
      </div>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="passsword" value="Password" />
        </div>
        <TextInput
          id="passsword"
          type="password"
          placeholder="••••••••"
          required
        />
      </div>
      <Button color="warning" type="submit">
        Kirim
      </Button>
    </form>
  );
}
